package handler

import (
	// "fmt"
	"net/http"
	"opensource/chaos/domain/dao/mongo"
	model "opensource/chaos/domain/model/mongo"
	webUtils "opensource/chaos/modules/oneclickdep/utils"
	"strconv"
)

func CreateContainerInfo(pathParams map[string]string, data []byte) interface{} {
	var request model.Container
	webUtils.ParseOuterRequest(data, &request)
	mongo.InsertContainer(request)
	return webUtils.ProcessResponseFully(http.StatusOK, nil, true)
}

func UpdateStateContainerInfo(pathParams map[string]string, data []byte) interface{} {
	request := make(map[string]int)
	cId := pathParams["cId"]
	cState := pathParams["cState"]
	request["state"], _ = strconv.Atoi(cState)
	mongo.UpdateContainerByIdUsingMap(cId, request)
	return webUtils.ProcessResponseFully(http.StatusOK, nil, true)
}

func MaskContainerInfo(pathParams map[string]string, data []byte) interface{} {
	request := make(map[string]int)
	cId := pathParams["cId"]
	request["state"] = model.CONTAINER_STATE_DELETE
	mongo.UpdateContainerByIdUsingMap(cId, request)
	return webUtils.ProcessResponseFully(http.StatusOK, nil, true)
}
