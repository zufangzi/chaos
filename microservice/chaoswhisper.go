package main

import (
	dockerapi "github.com/fsouza/go-dockerclient"
	consulApi "github.com/hashicorp/consul/api"
	"github.com/hashicorp/go-cleanhttp"
	"log"
	"opensource/chaos/microservice/utils"
	"os"
	"strconv"
	"strings"
)

var SCRIPT_HOME = "/usr/local/script/"
var START_SCRIPT = SCRIPT_HOME + "cnet-start.sh"
var START_LOG_CLT_SCRIPT = SCRIPT_HOME + "logstash-start.sh"
var STOP_SCRIPT = SCRIPT_HOME + "cnet-stop.sh"
var DELETE_SCRIPT = SCRIPT_HOME + "cnet-delete.sh"

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
		cnt, _ := strconv.Atoi(utils.GetShell("docker inspect " + msg.ID + " | grep Name | grep mesos | wc -l"))
		if cnt == 0 {
			continue
		}
		switch msg.Status {
		case "start":
			go startProcessor(msg)
		case "stop":
			go stopProcessor(msg)
		case "destroy":
			go destroyProcessor(msg)
		}
	}

	log.Fatalln("Docker listener closed!")
	// TODO 需要一个守护进程。定时扫有没有忘了删除的服务

}

func startProcessor(msg *dockerapi.APIEvents) {
	log.Println("now found a CREAT event!")
	log.Println("Step1: begin to process the container network")
	utils.GetShell("sh " + START_SCRIPT + " " + msg.ID)
	log.Println("Step2: begin to process log collection")
	// TODO
}

func stopProcessor(msg *dockerapi.APIEvents) {
	log.Println("now found a STOP event! ")
	log.Println("Step1: begin to clear link")
	utils.GetShell("sh " + STOP_SCRIPT + " " + msg.ID)
	log.Println("Step2: begin to process log collection")
	// TODO
	globalDeregister(msg)
}

func destroyProcessor(msg *dockerapi.APIEvents) {
	log.Println("now found a DESTROY event! ")
	log.Println("Step1: begin to clear up all info abt this container")
	utils.GetShell("sh " + DELETE_SCRIPT + " " + msg.ID)
	log.Println("Step2: begin to process log collection")
	// TODO
	globalDeregister(msg)
}

func globalDeregister(msg *dockerapi.APIEvents) {
	deregisterHostName := utils.GetShell("docker inspect -f '{{.Config.Hostname}}' " + msg.ID)
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
	serviceId := deregisterHostName + "_" + deregisterServiceName
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

func assert(err error) {
	if err != nil {
		log.Fatal(err)
	}
}
