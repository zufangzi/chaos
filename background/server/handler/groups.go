package handler

import (
	"opensource/chaos/background/server/dto"
	"opensource/chaos/background/server/dto/marathon"
	webUtils "opensource/chaos/background/server/utils"
	"opensource/chaos/background/utils"
	"opensource/chaos/background/utils/fasthttp"
)

func DeployGroupsHandler(pathParams map[string]string, data []byte) interface{} {
	var request dto.DeployGroupsRequest
	webUtils.ParseOuterRequest(data, &request)

	var groupsRequest marathon.MarathonGroupsRequest
	groupsRequest.Id = request.Id
	perGroups := make([]marathon.MarathonGroupsInfo, len(request.Groups))
	for i, v := range request.Groups {
		var group marathon.MarathonGroupsInfo
		group.Id = v.Id
		perApps := make([]marathon.MarathonAppsRequest, len(v.Apps))
		for j, app := range v.Apps {
			perApps[j] = *webUtils.BuildAppsRequest(app)
		}
		group.Apps = perApps
		perGroups[i] = group
	}
	groupsRequest.Groups = perGroups
	var resData map[string]interface{}
	resCode := fasthttp.JsonReqAndResHandler(utils.Path.MarathonGroupsUrl, groupsRequest, &resData, "POST")
	return webUtils.ProcessResponse(resCode, resData)
}
