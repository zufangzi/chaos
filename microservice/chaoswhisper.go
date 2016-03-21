package main

import (
	dockerapi "github.com/fsouza/go-dockerclient"
	consulApi "github.com/hashicorp/consul/api"
	"github.com/hashicorp/go-cleanhttp"
	"log"
	"opensource/chaos/microservice/utils"
	// "opensource/chaos/server/utils/fasthttp"
	"os"
	"strconv"
	"strings"
)

var SCRIPT_HOME = "/usr/local/scripts/"
var START_SCRIPT = SCRIPT_HOME + "docker-cnet-start.sh"
var START_LOG_CLT_SCRIPT = SCRIPT_HOME + "logstash-start.sh"
var STOP_SCRIPT = SCRIPT_HOME + "docker-cnet-stop.sh"
var DELETE_SCRIPT = SCRIPT_HOME + "docker-cnet-delete.sh"

var consulClient *consulApi.Client

func main() {
	if os.Getenv("DOCKER_HOST") == "" {
		os.Setenv("DOCKER_HOST", "unix:///var/run/docker.sock")
	}
	docker, err := dockerapi.NewClientFromEnv()
	assert(err)
	events := make(chan *dockerapi.APIEvents)
	assert(docker.AddEventListener(events))
	log.Println("Listening for Docker events ...")

	consulClient, _ = consulApi.NewClient(defaultConfig())

	for msg := range events {
		log.Printf("get docker event: %s now... \n", msg)

		// 不能用inspect了。都改成自己本地或者远程存储。
		if msg.Status != "destroy" {
			cnt, _ := strconv.Atoi(utils.GetShell("docker inspect " + msg.ID + " | grep Name | grep mesos | wc -l"))
			if cnt == 0 {
				continue
			}
		}

		switch msg.Status {
		case "start":
			go startProcessor(msg)
		case "die":
			go dieProcessor(msg)
		case "destroy":
			go destroyProcessor(msg)
		}
	}

	log.Fatalln("Docker listener closed!")
	// TODO 需要一个守护进程。定时扫有没有忘了删除的服务
	// 此处的逻辑为： 扫所有consul注册的服务，获取虚拟IP，进一步扫本机所有的容器，
	// 去进行一一比对。如果本机容器Name带有"mesos"字样，ip在上面存活的，且PID=0的，
	// 则认为服务已经停止，此时删除consul上的节点数据。
	// ------ 以上为删除已被停止的容器服务注册信息 -------
	// 每个宿主机上各自定时扫描汇报给cloud-server自己机器上的所有存活节点。
	// 由cloud-server抓取consul-server进行比对。取出不存在列表中的consul服务信息，
	// 再等待下一次比对。如果下一次比对仍然没有，则判定为不存在。此时进行删除。
	// ------ 以上为删除已不存在的容器服务注册信息 -------
}

func startProcessor(msg *dockerapi.APIEvents) {
	log.Println("now found a START event!")
	log.Println("Step1: begin to process the container network")
	utils.GetShell("sh " + START_SCRIPT + " " + getHostname(msg.ID, false))
	log.Println("Step2: begin to process log collection")
	// TODO 处理logstash、kafka topic相关脚本
}

func dieProcessor(msg *dockerapi.APIEvents) {
	log.Println("now found a DIE event! ")
	log.Println("Step1: begin to clear link")
	log.Println("sh " + STOP_SCRIPT + " " + getHostname(msg.ID, true))
	utils.GetShell("sh " + STOP_SCRIPT + " " + getHostname(msg.ID, true))
	log.Println("Step2: begin to process log collection")
	// TODO
	// globalDeregister(msg)
}

func destroyProcessor(msg *dockerapi.APIEvents) {
	log.Println("now found a DESTROY event! ")
	log.Println("Step1: begin to clear up all info abt this container")
	utils.GetShell("sh " + DELETE_SCRIPT + " " + getHostname(msg.ID, true))
	log.Println("Step2: begin to process log collection")
	// TODO
}

func globalDeregister(msg *dockerapi.APIEvents) {
	// var request dto.DeployAppsRequest
	// utils.ParseOuterRequest(data, &request)
	// deployInfo := utils.BuildAppsRequest(request)
	// var resData map[string]interface{}
	// resCode := fasthttp.JsonReqAndResHandler(utils.Path.MarathonAppsUrl, deployInfo, &resData, "POST")
	// return utils.ProcessResponse(resCode, resData)
}

// Deprecated because of die maybe will followed with destroy so that 'inspect' command will be not allowed
// 逻辑改为，只要包含容器ID。就可以删掉。
func globalDeregisterOld(msg *dockerapi.APIEvents) {
	deregisterEnv := utils.GetShell("docker inspect -f '{{.Config.Env}}' " + msg.ID)
	output := strings.Split(deregisterEnv, " ")
	var deregisterServiceName string
	for _, data := range output {
		if !strings.Contains(data, "SERVICE_NAME") {
			continue
		}
		deregisterServiceName = strings.Split(data, "=")[1]
		break
	}
	serviceId := getHostname(msg.ID, true) + "_" + deregisterServiceName
	consulClient.Agent().ServiceDeregister(serviceId)
	log.Printf("consul info clear up. serviceId: %s", serviceId)
}

func defaultConfig() *consulApi.Config {
	hostIp := "127.0.0.1:8500"
	config := &consulApi.Config{
		Address:    hostIp,
		Scheme:     "http",
		HttpClient: cleanhttp.DefaultClient(),
	}
	return config
}

// 需要被替代。不能用inspect
func getHostname(id string, byFullId bool) string {
	if !byFullId {
		return utils.GetShell("docker inspect -f {{.Config.Hostname}} " + id)
	} else {
		return id[:12] // 这边考虑到了destroy的话就没有hostname了，这时候就通过id进行分割获取。
	}

}

func assert(err error) {
	if err != nil {
		log.Fatal(err)
	}
}
