package main

import (
	dockerapi "github.com/fsouza/go-dockerclient"
	consulApi "github.com/hashicorp/consul/api"
	"github.com/hashicorp/go-cleanhttp"
	"log"
	"opensource/chaos/background/utils"
	// "opensource/chaos/server/utils/fasthttp"
	"os"
	"strconv"
	"strings"
)

var SCRIPT_WGET_HOME = "http://10.32.27.11/cc/preonline/host/"
var SCRIPT_HOME = "/usr/local/scripts/"
var START_SCRIPT = "docker-cnet-start.sh"
var START_LOG_CLT_SCRIPT = "logstash-start.sh"
var STOP_SCRIPT = "docker-cnet-stop.sh"
var DELETE_SCRIPT = "docker-cnet-delete.sh"

var consulClient *consulApi.Client
var dockerClinet *dockerapi.Client

func init() {
	utils.GetShell("wget " + SCRIPT_WGET_HOME + START_SCRIPT + " && mv " + START_SCRIPT + " " + SCRIPT_HOME)
	utils.GetShell("wget " + SCRIPT_WGET_HOME + STOP_SCRIPT + " && mv " + STOP_SCRIPT + " " + SCRIPT_HOME)
	utils.GetShell("wget " + SCRIPT_WGET_HOME + DELETE_SCRIPT + " && mv " + DELETE_SCRIPT + " " + SCRIPT_HOME)
}

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
	cnt, _ := strconv.Atoi(utils.GetShell("docker inspect " + msg.ID + " | grep SERVICE_NAME | wc -l"))
	if cnt == 0 {
		return
	}
	log.Println("Step1: begin to process the container network")
	utils.GetShell("sh " + SCRIPT_HOME + START_SCRIPT + " " + getHostname(msg.ID, false))
	log.Println("Step2: begin to process log collection")
	// TODO 处理logstash、kafka topic相关脚本
}

func dieProcessor(msg *dockerapi.APIEvents) {
	log.Println("now found a DIE event! ")
	log.Println("Step1: deregister")
	needDeregister := globalDeregister(msg)
	if !needDeregister {
		log.Println("No need to deregister. No Service Found")
		return
	}
	log.Println("Step2: begin to clear link")
	utils.GetShell("sh " + SCRIPT_HOME + STOP_SCRIPT + " " + getHostname(msg.ID, true))
	log.Println("Step3: begin to process log collection")
}

func destroyProcessor(msg *dockerapi.APIEvents) {
	log.Println("now found a DESTROY event! ")
	log.Println("Step1: begin to clear up all info abt this container")
	// 即使没有start和stop，去调用delete脚本也不会有问题。
	utils.GetShell("sh " + SCRIPT_HOME + DELETE_SCRIPT + " " + getHostname(msg.ID, true))
	log.Println("Step2: begin to process log collection")
	// TODO
}

func globalDeregister(msg *dockerapi.APIEvents) bool {
	data, _ := consulClient.Agent().Services()
	for key, _ := range data {
		if !strings.Contains(key, "_") {
			continue
		}
		serviceId := getHostname(msg.ID, true)
		if strings.Split(key, "_")[0] == serviceId {
			log.Println("found serviceId: " + key)
			consulClient.Agent().ServiceDeregister(key)
			return true
		}
	}
	return false
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
