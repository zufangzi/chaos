package utils

import (
	"encoding/json"
	"log"
	"net/http"
	"opensource/chaos/server/dto"
	"opensource/chaos/server/dto/marathon"
	"strconv"
)

var fastDocker *FastDocker

func init() {
	fastDocker = new(FastDocker)
}

func ParseOuterRequest(body []byte, request interface{}) {
	err := json.Unmarshal(body, &request)
	CheckError(err)
	log.Println("the request data is: ", request)
}

func BuildAppsRequest(request dto.DeployAppsRequest) (appsReq *marathon.MarathonAppsRequest) {
	deployInfo := marathon.NewMarathonAppsRequest()
	deployInfo.Id = request.Id
	deployInfo.Cpus, _ = strconv.ParseFloat(request.Cpus, 64)
	deployInfo.Mem, _ = strconv.ParseFloat(request.Mem, 64)
	deployInfo.Instances, _ = strconv.Atoi(request.Instances)
	container := marathon.NewMarathonDockerContainer()
	if request.Version != "" {
		container.Image = request.Image + ":" + request.Version
	} else {
		// 拿到最新的时间戳的tag
		container.Image, _, _ = fastDocker.GetImageAndTagByFreshness(request.Image, "", "", 0, false)

	}

	// 如果设定了端口，那么就处理
	var ports []marathon.MarathonDockerPort
	flag := false
	if request.ExportPorts != nil && len(request.ExportPorts) > 0 {
		ports = make([]marathon.MarathonDockerPort, len(request.ExportPorts), len(request.ExportPorts)+10)
		for i, v := range request.ExportPorts {
			port := marathon.NewMarathonDockerPort()
			port.ContainerPort, _ = strconv.Atoi(v.ContainerPort)
			ports[i] = *port
			if v.ContainerPort == "22" {
				flag = true
			}
		}
	} else {
		ports = make([]marathon.MarathonDockerPort, 1)
		ports[0] = *AddDefaultPorts()
		flag = true
	}
	if !flag {
		ports = append(ports, *AddDefaultPorts())
	}
	container.PortMappings = ports
	container.Volumes = make([]interface{}, 0)
	deployInfo.Container = make(map[string]interface{})
	deployInfo.Container["docker"] = container
	deployInfo.Labels = make(map[string]string)
	deployInfo.Labels["dd-version"] = request.Version
	return deployInfo
}

func AddDefaultPorts() *marathon.MarathonDockerPort {
	port := marathon.NewMarathonDockerPort()
	port.ContainerPort = 22
	port.Protocol = "tcp"
	return port
}

func ProcessResponse(code int, response interface{}) interface{} {
	return ProcessResponseFully(code, response, true)
}

func ProcessResponseFully(code int, response interface{}, shouldHideSuccessInfo bool) interface{} {

	if !shouldHideSuccessInfo {
		return response
	}

	if code != http.StatusCreated && code != http.StatusOK && code != http.StatusAccepted {
		return response
	} else {
		return map[string]string{"status": "ok"}
	}
}

func CheckError(err error) {
	if err != nil {
		log.Fatalln("fatal occur.", err.Error())
		panic(err)
	}
}
