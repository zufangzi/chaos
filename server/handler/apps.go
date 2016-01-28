package handler

import (
	"net/http"
	"opensource/chaos/server/dto"
	"opensource/chaos/server/dto/marathon"
	"opensource/chaos/server/handler/service"
	"opensource/chaos/server/utils"
	"opensource/chaos/server/utils/fasthttp"
	"strconv"
	"strings"
)

// 逻辑为：放入仓库的时候，即每个模块携带时间戳，每次前端构建时候传入比如zookeeper
// 则后端则从私库里捞出zookeeper所有模块，并按时间倒叙取出最新的zk模块镜像进行部署
// 则，在回滚时候，捞出倒数第二新的模块进行重新部署。部署时候更新labels即可。
func RollbackAppsHandler(data []byte) interface{} {
	var request dto.RollbackAppsBatchRequest
	utils.ParseOuterRequest(data, &request)

	requestBatch := make([]dto.DeployAppsRequest, len(request.Batch))
	for i, v := range request.Batch {
		_, image, tag := utils.DockerClient.GetPreviousImageAndTag(v.Id, v.Version, "")
		var request dto.DeployAppsRequest
		request.Id = v.Id
		request.Image = image
		request.Version = tag
		requestBatch[i] = request
	}
	appsBatchRequest := dto.DeployAppsBatchRequest{}
	appsBatchRequest.Batch = requestBatch
	resData, code := service.CreateOrUpdateAppsService(appsBatchRequest)
	return utils.ProcessResponse(code, resData)
}

func CreateAppsHandler(data []byte) interface{} {
	var request dto.DeployAppsRequest
	utils.ParseOuterRequest(data, &request)
	deployInfo := utils.BuildAppsRequest(request)
	var resData map[string]interface{}
	resCode := fasthttp.JsonReqAndResHandler(utils.Path.MarathonAppsUrl, deployInfo, &resData, "POST")
	return utils.ProcessResponse(resCode, resData)
}

func CreateOrUpdateAppsHandler(data []byte) interface{} {
	var request dto.DeployAppsBatchRequest
	utils.ParseOuterRequest(data, &request)
	resData, resCode := service.CreateOrUpdateAppsService(request)
	return utils.ProcessResponse(resCode, resData)
}

func GetInfoAppsHandler(data []byte) interface{} {
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
	return utils.ProcessResponseFully(http.StatusOK, appsGlobalInfos, false)
}

func GetSingleAppsHandler(data []byte) interface{} {
	return nil
}

func DeleteAppsHandler(data []byte) interface{} {
	return nil
}
