package service

import (
	"opensource/chaos/background/server/dto"
	"opensource/chaos/background/server/dto/marathon"
	webUtils "opensource/chaos/background/server/utils"
	"opensource/chaos/background/utils"
	"opensource/chaos/background/utils/fasthttp"
)

func CreateOrUpdateAppsService(request dto.DeployAppsBatchRequest) (interface{}, int) {
	finalRequest := make([]marathon.MarathonAppsRequest, len(request.Batch))
	for i, v := range request.Batch {
		deployInfo := webUtils.BuildAppsRequest(v)
		finalRequest[i] = *deployInfo
	}
	var response map[string]interface{}
	code := fasthttp.JsonReqAndResHandler(utils.Path.MarathonAppsUrl, finalRequest, &response, "PUT")
	return response, code
}
