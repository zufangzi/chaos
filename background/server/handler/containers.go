package handler

import (
	// "fmt"
	"net/http"
	"opensource/chaos/background/server/domain/mongo"
	"opensource/chaos/background/server/dto/model"
	webUtils "opensource/chaos/background/server/utils"
	"strconv"
)

func CreateContainerInfo(pathParams map[string]string, data []byte) interface{} {
	var request model.Container
	webUtils.ParseOuterRequest(data, &request)
	mongo.InsertContainer(request)
	return webUtils.ProcessResponseFully(http.StatusOK, nil, true)
}

func UpdateStateContainerInfo(pathParams map[string]string, data []byte) interface{} {
	var request map[string]int
	cId := pathParams["cId"]
	cState := pathParams["state"]
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
