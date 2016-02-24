package handler

import (
	"opensource/chaos/server/dto"
	"opensource/chaos/server/dto/marathon"
	"opensource/chaos/server/utils"
	"opensource/chaos/server/utils/fasthttp"
)

func DeployGroupsHandler(pathParams map[string]string, data []byte) interface{} {
	var request dto.DeployGroupsRequest
	utils.ParseOuterRequest(data, &request)

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
	var resData map[string]interface{}
	resCode := fasthttp.JsonReqAndResHandler(utils.Path.MarathonGroupsUrl, groupsRequest, &resData, "POST")
	return utils.ProcessResponse(resCode, resData)
}
