package utils

import (
	"encoding/json"
	"log"
	"opensource/chaos/common"
	"opensource/chaos/domain/model/marathon"
	"opensource/chaos/modules/oneclickdep/entity"
	"strconv"
)

var fastDocker *common.FastDocker

func init() {
	fastDocker = new(common.FastDocker)
}

func ParseOuterRequest(body []byte, request interface{}) {
	err := json.Unmarshal(body, &request)
	CheckError(err)
	log.Println("the request data is: ", request)
}

func BuildAppsRequest(request entity.DeployAppsRequest) (appsReq *marathon.MarathonAppsRequest) {
	deployInfo := marathon.NewMarathonAppsRequest()
	deployInfo.Id = request.Id
	deployInfo.Cpus, _ = strconv.ParseFloat(request.Cpus, 64)
	deployInfo.Mem, _ = strconv.ParseFloat(request.Mem, 64)
	deployInfo.Instances, _ = strconv.Atoi(request.Instances)
	container := marathon.NewMarathonContainer()
	// 不能带上http前缀，image不能有http前缀
	container.Docker.Image, _, _ = fastDocker.GetImageAndTagByFreshness(request.Image, request.Version, "", 0, false)

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
	// 由于采用了ovs+none的方式，所以不需要做端口映射了
	// container.PortMappings = ports
	deployInfo.Container = *container
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
	return ProcessResponseFully(code, response, false)
}

func ProcessResponseFully(code int, response interface{}, shouldHideSuccessInfo bool) interface{} {

	if shouldHideSuccessInfo {
		return map[string]interface{}{
			"status": code,
		}
	}
	return map[string]interface{}{
		"status": code,
		"data":   response,
	}
}

func CheckError(err error) {
	if err != nil {
		log.Fatalln("fatal occur.", err.Error())
		panic(err)
	}
}
