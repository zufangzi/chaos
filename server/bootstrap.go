package main

import (
	// "container/list"
	// "encoding/json"
	"fmt"
	"github.com/ant0ine/go-json-rest/rest"
	"opensource/chaos/server/dto"
	"opensource/chaos/server/fasthttp"
	"opensource/chaos/server/marathon"
	// "github.com/gorilla/mux"
	"github.com/tjz101/goprop"
	// "html"
	"html/template"
	"io/ioutil"
	"log"
	// "net"
	// "bytes"
	// "encoding/gob"
	"net/http"
	"os"
	"path"
	"strconv"
	"strings"
	"time"
)

// Constant
const (
	STATIC_DIR = "/src/opensource/chaos/foreground"
	PROP_FILE  = "/src/opensource/chaos/resources/core.properties"
	TMPL_NAME  = "html"
)

// global vars
var goCore GoCore

// templates := make(map[string]*template.Template) will occur "non-declaration statement outside function body" error
var templates = make(map[string]*template.Template)

func main() {
	fmt.Println("now begin to run server, please wait...")
	api := rest.NewApi()
	api.Use(rest.DefaultDevStack...)
	route, _ := rest.MakeRouter(
		rest.Post("/simpledeploy/apps", deploySimpleHandler),
		rest.Get("/info", appInfoHandler),
		rest.Get("/info/#appId", appDetailInfoHandler),
	)
	api.SetApp(route)

	fileServer := http.FileServer(http.Dir(os.Getenv("GOPATH") + STATIC_DIR))
	http.HandleFunc("/", entranceGuarder(indexHandler))
	http.Handle("/static/", http.StripPrefix("/static", fileServer))
	http.Handle("/api/", http.StripPrefix("/api", api.MakeHandler()))
	http.ListenAndServe(":8080", nil)
}

func deploySimpleHandler(w rest.ResponseWriter, r *rest.Request) {

	var request dto.DeploySimpleRequest
	err := r.DecodeJsonPayload(&request)
	CheckError(err)
	log.Println("the request data is: ", request)

	deployInfo := marathon.NewMarathonAppsRequest()
	deployInfo.Id = request.Id
	deployInfo.Cpus, _ = strconv.ParseFloat(request.Cpus, 64)
	deployInfo.Mem, _ = strconv.ParseFloat(request.Mem, 64)
	deployInfo.Instances, _ = strconv.Atoi(request.Instances)

	container := marathon.NewMarathonDockerContainer()
	container.Image = "10.32.27.82:5000/" + request.Image

	ports := make([]marathon.MarathonDockerPort, len(request.ExportPorts))
	for i, v := range request.ExportPorts {
		port := marathon.NewMarathonDockerPort()
		port.ContainerPort, _ = strconv.Atoi(v.ContainerPort)
		ports[i] = *port
	}

	container.PortMappings = ports
	container.Volumes = make([]interface{}, 0)

	deployInfo.Container = make(map[string]interface{})
	deployInfo.Container["docker"] = container

	// request
	var response map[string]interface{}
	code := fasthttp.JsonReqAndResHandler(goCore.MarathonAppsUrl, deployInfo, &response, "POST")
	fmt.Println("now code is: " + strconv.Itoa(code))
	if code != http.StatusCreated {
		w.WriteJson(response)
	} else {
		w.WriteJson(map[string]string{"status": "ok"})
	}
}

func appInfoHandler(w rest.ResponseWriter, r *rest.Request) {
	var marathonApps marathon.MarathonAppsGlobalInfoResponse
	fasthttp.JsonReqAndResHandler(goCore.MarathonAppsUrl, nil, &marathonApps, "GET")
	appsCnt := len(marathonApps.Apps)

	// should not code like this: appsGlobalInfos := [appsCnt]entity.AppsGlobalInfo{}
	appsGlobalInfos := make([]dto.AppsGlobalInfoResponse, appsCnt)

	for i, v := range marathonApps.Apps {
		var perApp dto.AppsGlobalInfoResponse
		if strings.LastIndex(v.Id, "/") == -1 {
			perApp.Id = v.Id
		} else {
			perApp.Id = v.Id[strings.LastIndex(v.Id, "/")+1:]
		}
		perApp.Cpus = strconv.FormatFloat(v.Cpus, 'f', 1, 64)
		perApp.CurrentInstances = strconv.Itoa(v.TasksRunning)
		fmt.Println(v)
		if strings.LastIndex(v.Id, "/") <= 0 { // exclude like /zk or zk
			perApp.Group = "No Groups"
		} else {
			perApp.Group = v.Id[0:strings.LastIndex(v.Id, "/")]
		}
		perApp.Instances = strconv.Itoa(v.Instances)
		perApp.Mem = strconv.FormatFloat(v.Mem, 'f', 1, 64)
		if v.TasksHealthy == 0 && v.TasksUnhealthy == 0 { // when no and healthy check
			perApp.Healthy = "100"
		} else {
			perApp.Healthy = strconv.FormatFloat(float64(v.TasksHealthy)/float64(v.TasksHealthy+v.TasksUnhealthy), 'f', 1, 64)
		}
		perApp.FormatStatus(v.TasksStaged)
		appsGlobalInfos[i] = perApp
	}
	w.WriteJson(appsGlobalInfos)
}

func appDetailInfoHandler(w rest.ResponseWriter, r *rest.Request) {
	fmt.Println(r.PathParam("appId"))
}

func entranceGuarder(method http.HandlerFunc) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		begin := time.Now().UnixNano()
		defer func() {
			// func().(xx) means method return type cast
			// if in args type cast case. you can use string(xx) or xx.(string)
			if e, ok := recover().(error); ok {
				http.Error(w, e.Error(), http.StatusInternalServerError)
				log.Println("catchable system error occur. " + e.Error())
			}
			log.Printf("the request: %s cost: %d ms\n", r.URL.RequestURI(), ((time.Now().UnixNano() - begin) / 1000000))
		}()
		method(w, r)
	}
}

func indexHandler(w http.ResponseWriter, r *http.Request) {
	renderPage(w, "index", nil)
}

func renderPage(w http.ResponseWriter, tmpl string, values map[string]interface{}) {
	err := templates[tmpl+"."+TMPL_NAME].Execute(w, values)
	CheckError(err)
}

func CheckError(err error) {
	if err != nil {
		log.Fatalln("fatal occur.", err.Error())
		panic(err)
	}
}

func init() {
	goCore = GoCore{"", ""}
	prop := goprop.NewProp()
	prop.Read(os.Getenv("GOPATH") + PROP_FILE)
	goCore.MarathonAppsUrl = prop.Get("marathon.apps.url")
	goCore.MarathonGroupsUrl = prop.Get("marathon.groups.url")

	htmlPath := os.Getenv("GOPATH") + STATIC_DIR + string(os.PathSeparator) + TMPL_NAME + string(os.PathSeparator)
	fileInfoArray, err := ioutil.ReadDir(htmlPath)
	CheckError(err)

	var fileName, filePath string
	for _, fileInfo := range fileInfoArray {
		fileName = fileInfo.Name()
		if ext := path.Ext(fileName); ext != ("." + TMPL_NAME) {
			continue
		}
		filePath = htmlPath + fileName
		log.Println("loading template: " + filePath)
		t := template.Must(template.New(fileName).Delims("[[", "]]").ParseFiles(filePath))
		templates[fileName] = t
	}
}

type GoCore struct {
	MarathonAppsUrl   string
	MarathonGroupsUrl string
}
