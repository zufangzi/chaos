package main

import (
	"opensource/chaos/server/dto"
	"opensource/chaos/server/dto/marathon"
	"opensource/chaos/server/utils"
	"opensource/chaos/server/utils/fasthttp"
)

func createOrUpdateAppsService(request dto.DeployAppsSimpleRequest) (interface{}, int) {
	deployInfo := utils.BuildAppsRequest(request)
	finalRequest := make([]marathon.MarathonAppsRequest, 1)
	finalRequest[0] = *deployInfo
	var response map[string]interface{}
	code := fasthttp.JsonReqAndResHandler(utils.Path.MarathonAppsUrl, finalRequest, &response, "PUT")
	return response, code
}
