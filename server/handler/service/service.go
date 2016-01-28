package service

import (
	"opensource/chaos/server/dto"
	"opensource/chaos/server/dto/marathon"
	"opensource/chaos/server/utils"
	"opensource/chaos/server/utils/fasthttp"
)

func CreateOrUpdateAppsService(request dto.DeployAppsBatchRequest) (interface{}, int) {
	finalRequest := make([]marathon.MarathonAppsRequest, len(request.Batch))
	for i, v := range request.Batch {
		deployInfo := utils.BuildAppsRequest(v)
		finalRequest[i] = *deployInfo
	}
	var response map[string]interface{}
	code := fasthttp.JsonReqAndResHandler(utils.Path.MarathonAppsUrl, finalRequest, &response, "PUT")
	return response, code
}
