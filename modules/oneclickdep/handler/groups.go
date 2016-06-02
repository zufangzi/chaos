package handler

import (
	"opensource/chaos/common"
	"opensource/chaos/common/fasthttp"
	"opensource/chaos/domain/model/marathon"
	"opensource/chaos/modules/oneclickdep/entity"
	webUtils "opensource/chaos/modules/oneclickdep/utils"
)

func DeployGroupsHandler(pathParams map[string]string, data []byte) interface{} {
	var request entity.DeployGroupsRequest
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
	resCode := fasthttp.JsonReqAndResHandler(common.Path.MarathonGroupsUrl, groupsRequest, &resData, "POST")
	return webUtils.ProcessResponse(resCode, resData)
}
