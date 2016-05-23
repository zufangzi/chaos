package main

import (
	"flag"
	"fmt"
	dockerapi "github.com/fsouza/go-dockerclient"
	consulApi "github.com/hashicorp/consul/api"
	"github.com/hashicorp/go-cleanhttp"
	"log"
	"net/http"
	"opensource/chaos/background/server/dto/model"
	"opensource/chaos/background/utils"
	"opensource/chaos/background/utils/fasthttp"
	"os"
	"strconv"
	"strings"
)

var ARGS_SCRIPT_WGET_HOME = flag.String("script-home", "http://10.32.27.11/cc/preonline/host/", "Script home with http:// prefix")
var ARGS_CLOUD_SERVER_URL = flag.String("cloud-server", "10.32.27.76:8080", "Cloud Server Address Without 'http://' prefix")

var PATH_CONTAINER = "http://" + *ARGS_CLOUD_SERVER_URL + "/api/containers"
var PATH_IPHOLDER = "http://" + *ARGS_CLOUD_SERVER_URL + "/api/ipholder/%s"

var SCRIPT_HOME = "/usr/local/scripts/"
var START_SCRIPT = "docker-cnet-start.sh"
var START_LOG_CLT_SCRIPT = "logstash-start.sh"
var STOP_SCRIPT = "docker-cnet-stop.sh"
var DELETE_SCRIPT = "docker-cnet-delete.sh"

var consulClient *consulApi.Client
var dockerClient *dockerapi.Client

func main() {

	flag.Parse()

	if os.Getenv("DOCKER_HOST") == "" {
		os.Setenv("DOCKER_HOST", "unix:///var/run/docker.sock")
	}

	utils.GetShell("wget " + *ARGS_SCRIPT_WGET_HOME + START_SCRIPT + " && mv " + START_SCRIPT + " " + SCRIPT_HOME)
	utils.GetShell("wget " + *ARGS_SCRIPT_WGET_HOME + STOP_SCRIPT + " && mv " + STOP_SCRIPT + " " + SCRIPT_HOME)
	utils.GetShell("wget " + *ARGS_SCRIPT_WGET_HOME + DELETE_SCRIPT + " && mv " + DELETE_SCRIPT + " " + SCRIPT_HOME)

	dockerClient, _ = dockerapi.NewClientFromEnv()
	events := make(chan *dockerapi.APIEvents)
	assert(dockerClient.AddEventListener(events))
	log.Println("Listening for Docker events ...")

	consulClient, _ = consulApi.NewClient(defaultConfig())

	for msg := range events {
		log.Printf("get docker event: %s now... \n", msg)
		switch msg.Status {
		case "start":
			go startProcessor(msg, getHostname(msg.ID, true))
		case "die":
			go dieProcessor(msg, getHostname(msg.ID, true))
		case "destroy":
			go destroyProcessor(msg, getHostname(msg.ID, true))
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

func startProcessor(msg *dockerapi.APIEvents, shortId string) {
	log.Println("now found a START event!")
	cnt, _ := strconv.Atoi(utils.GetShell("docker inspect " + shortId + " | grep SERVICE_NAME | wc -l"))
	if cnt == 0 {
		return
	}

	log.Println("Step1: begin to process container info restore")
	storeContainerInfo(shortId)

	log.Println("Step2: begin to process the container network")
	ip := getIp(shortId)
	fmt.Println("IP is: " + ip)
	utils.GetShell("sh " + SCRIPT_HOME + START_SCRIPT + " " + shortId + " " + ip)

	log.Println("Step3: begin to process log collection")
	// TODO 处理logstash、kafka topic相关脚本
}

func dieProcessor(msg *dockerapi.APIEvents, shortId string) {
	log.Println("now found a DIE event! ")
	log.Println("Step1: deregister")
	needDeregister := globalDeregister(msg, shortId)
	if !needDeregister {
		log.Println("No need to deregister. No Service Found")
		return
	}
	log.Println("Step2: begin to clear link")
	utils.GetShell("sh " + SCRIPT_HOME + STOP_SCRIPT + " " + shortId)
	log.Println("Step3: begin to process log collection")
}

func destroyProcessor(msg *dockerapi.APIEvents, shortId string) {
	log.Println("now found a DESTROY event! ")
	log.Println("Step1: begin to clear up all info abt this container")
	// 即使没有start和stop，去调用delete脚本也不会有问题。
	utils.GetShell("sh " + SCRIPT_HOME + DELETE_SCRIPT + " " + shortId)
	delIp(shortId)
	log.Println("Step2: begin to mask container info")
	maskContainerInfo(shortId)
	log.Println("Step3: begin to process log collection")
	// TODO
}

func globalDeregister(msg *dockerapi.APIEvents, shortId string) bool {
	data, _ := consulClient.Agent().Services()
	for key, _ := range data {
		if !strings.Contains(key, "_") {
			continue
		}
		if strings.Split(key, "_")[0] == shortId {
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

func getHostname(id string, simpleCut bool) string {
	if !simpleCut {
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

func storeContainerInfo(cId string) {
	var resData map[string]interface{}
	var request model.Container
	request.BizId = cId
	docker, _ := dockerClient.InspectContainer(cId)
	request.BizName = docker.Name
	request.State = model.CONTAINER_STATE_STARTED
	fasthttp.JsonReqAndResHandler(PATH_CONTAINER, request, &resData, "POST")
	if int(resData["status"].(float64)) != http.StatusOK {
		log.Println("[CHAOSWHISPER] store container info fail... CID: " + cId)
	}
}

func maskContainerInfo(cId string) {
	var resData map[string]interface{}
	fasthttp.JsonReqAndResHandler(PATH_CONTAINER+"/soft/"+cId, nil, &resData, "DELETE")
	if int(resData["status"].(float64)) != http.StatusOK {
		log.Println("[CHAOSWHISPER] mask container info fail... CID: " + cId)
	}
}

func getIp(cId string) string {
	var resData map[string]interface{}
	uri := fmt.Sprintf(PATH_IPHOLDER, cId)
	fasthttp.JsonReqAndResHandler(uri, nil, &resData, "GET")
	if int(resData["status"].(float64)) == http.StatusOK {
		return resData["data"].(string)
	} else {
		// ALARM. TODO
		return "127.0.0.1"
	}
}

func delIp(cId string) {
	var resData map[string]interface{}
	uri := fmt.Sprintf(PATH_IPHOLDER, cId)
	fasthttp.JsonReqAndResHandler(uri, nil, &resData, "DELETE")
	if int(resData["status"].(float64)) != http.StatusOK {
		log.Println("[CHAOSWHISPER] delete ip fail... CID: " + cId)
	}
}
