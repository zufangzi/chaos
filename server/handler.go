package main

import (
	"fmt"
	"github.com/ant0ine/go-json-rest/rest"
	"net/http"
	"opensource/chaos/server/dto"
	"opensource/chaos/server/dto/marathon"
	"opensource/chaos/server/utils"
	"opensource/chaos/server/utils/fasthttp"
	"strconv"
	"strings"
)

var docker *utils.FastDocker

func init() {
	docker = new(utils.FastDocker)
}

func Route() (rest.App, error) {

	return rest.MakeRouter(
		rest.Post("/deploy/apps/rollback", utils.RestGuarder(rollbackAppsHandler)),
		rest.Post("/deploy/apps", utils.RestGuarder(createAppsHandler)),
		rest.Post("/deploy/apps/updater", utils.RestGuarder(createOrUpdateAppsHandler)),
		rest.Post("/deploy/groups", utils.RestGuarder(deployGroupsHandler)),
		rest.Get("/info", utils.RestGuarder(appInfoHandler)),
		rest.Get("/info/#appId", utils.RestGuarder(appDetailInfoHandler)),
	)
}

// 逻辑为：放入仓库的时候，即每个模块携带时间戳，每次前端构建时候传入比如zookeeper
// 则后端则从私库里捞出zookeeper所有模块，并按时间倒叙取出最新的zk模块镜像进行部署
// 则，在回滚时候，捞出倒数第二新的模块进行重新部署。部署时候更新labels即可。
func rollbackAppsHandler(w rest.ResponseWriter, r *rest.Request) {

	var request []dto.RollbackAppsRequest
	utils.ParseOuterRequest(r, &request)

	resData := make(map[string]int)
	for _, v := range request {
		image := docker.GetPreviousImage(v.Id, v.Version, "")
		var request dto.DeployAppsSimpleRequest
		request.Id = v.Id
		request.Image = image
		fmt.Println("helloL ", request)
		response, code := createOrUpdateAppsService(request)
		fmt.Println(response)
		resData[v.Id] = code
	}
	utils.ProcessResponseFully(w, http.StatusOK, resData, false)
}

func createAppsHandler(w rest.ResponseWriter, r *rest.Request) {
	var request dto.DeployAppsSimpleRequest
	utils.ParseOuterRequest(r, &request)
	deployInfo := utils.BuildAppsRequest(request)
	var response map[string]interface{}
	code := fasthttp.JsonReqAndResHandler(utils.Path.MarathonAppsUrl, deployInfo, &response, "POST")
	utils.ProcessResponse(w, code, response)
}

func createOrUpdateAppsHandler(w rest.ResponseWriter, r *rest.Request) {
	var request dto.DeployAppsSimpleRequest
	utils.ParseOuterRequest(r, &request)
	response, code := createOrUpdateAppsService(request)
	utils.ProcessResponse(w, code, response)
}

func deployGroupsHandler(w rest.ResponseWriter, r *rest.Request) {
	var request dto.DeployGroupsSimpleRequest
	utils.ParseOuterRequest(r, &request)
	var groupsRequest marathon.MarathonGroupsRequest
	groupsRequest.Id = request.Id
	perGroups := make([]marathon.MarathonGroupsInfo, len(request.Groups))
	for i, v := range request.Groups {
		var group marathon.MarathonGroupsInfo
		group.Id = v.Id
		perApps := make([]marathon.MarathonAppsRequest, len(v.Apps))
		for j, app := range v.Apps {
			perApps[j] = *utils.BuildAppsRequest(app)
		}
		group.Apps = perApps
		perGroups[i] = group
	}
	groupsRequest.Groups = perGroups
	var response map[string]interface{}
	code := fasthttp.JsonReqAndResHandler(utils.Path.MarathonGroupsUrl, groupsRequest, &response, "POST")
	utils.ProcessResponse(w, code, response)
}

func appInfoHandler(w rest.ResponseWriter, r *rest.Request) {
	var marathonApps marathon.MarathonAppsGlobalInfoResponse
	fasthttp.JsonReqAndResHandler(utils.Path.MarathonAppsUrl, nil, &marathonApps, "GET")
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
		if strings.LastIndex(v.Id, "/") <= 0 { // exclude like /zk or zk
			perApp.Group = "No Groups"
		} else {
			perApp.Group = v.Id[0:strings.LastIndex(v.Id, "/")]
		}
		perApp.Instances = strconv.Itoa(v.Instances)
		perApp.Mem = strconv.FormatFloat(v.Mem, 'f', 1, 64)
		perApp.Healthy = strconv.FormatFloat(100*float64(v.TasksRunning)/float64(v.Instances), 'f', 1, 64)
		perApp.FormatStatus(v.TasksStaged)
		appsGlobalInfos[i] = perApp
	}
	w.WriteJson(appsGlobalInfos)
}

func appDetailInfoHandler(w rest.ResponseWriter, r *rest.Request) {
	fmt.Println(r.PathParam("appId"))
}
