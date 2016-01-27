package utils

import (
	"fmt"
	"github.com/ant0ine/go-json-rest/rest"
	"log"
	"net/http"
	"opensource/chaos/server/dto"
	"opensource/chaos/server/dto/marathon"
	"strconv"
	"time"
)

var fastDocker *FastDocker

func init() {
	fastDocker = new(FastDocker)
}

func ParseOuterRequest(r *rest.Request, request interface{}) {
	err := r.DecodeJsonPayload(&request)
	CheckError(err)
	log.Println("the request data is: ", request)
}

func BuildAppsRequest(request dto.DeployAppsSimpleRequest) (appsReq *marathon.MarathonAppsRequest) {
	deployInfo := marathon.NewMarathonAppsRequest()
	deployInfo.Id = request.Id
	deployInfo.Cpus, _ = strconv.ParseFloat(request.Cpus, 64)
	deployInfo.Mem, _ = strconv.ParseFloat(request.Mem, 64)
	deployInfo.Instances, _ = strconv.Atoi(request.Instances)
	container := marathon.NewMarathonDockerContainer()
	if request.Version != "" {
		container.Image = "10.32.27.82:5000/" + request.Image + ":" + request.Version
	} else {
		// 拿到最新的时间戳的tag
		container.Image = fastDocker.GetImageByFreshness(request.Image, "", "", 0, false)

	}
	fmt.Println("d")

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

func ProcessResponse(w rest.ResponseWriter, code int, response interface{}) {
	ProcessResponseFully(w, code, response, true)
}

func ProcessResponseFully(w rest.ResponseWriter, code int, response interface{}, shouldHideSuccessInfo bool) {

	if !shouldHideSuccessInfo {
		w.WriteJson(response)
		return
	}

	if code != http.StatusCreated && code != http.StatusOK && code != http.StatusAccepted {
		w.WriteJson(response)
	} else {
		w.WriteJson(map[string]string{"status": "ok"})
	}
}

func CheckError(err error) {
	if err != nil {
		log.Fatalln("fatal occur.", err.Error())
		panic(err)
	}
}

func RestGuarder(method rest.HandlerFunc) rest.HandlerFunc {
	return func(w rest.ResponseWriter, r *rest.Request) {
		begin := time.Now().UnixNano()
		defer func() {
			// func().(xx) means method return type cast
			// if in args type cast case. you can use string(xx) or xx.(string)
			if e, ok := recover().(error); ok {
				rest.Error(w, e.Error(), http.StatusInternalServerError)
				log.Println("catchable system error occur. " + e.Error())
			}
			log.Printf("the request: %s cost: %d ms\n", r.URL.RequestURI(), ((time.Now().UnixNano() - begin) / 1000000))
		}()
		method(w, r)
	}
}
