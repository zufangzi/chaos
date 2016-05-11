package mongo

import (
	"opensource/chaos/background/server/dto/model"
)

const COLLECTION_SERVICE = "service"

func GetServiceInfoById(id string) model.Service {
	result := model.Service{}
	queryById(COLLECTION_SERVICE, &result, id)
	return result
}

func InsertService(service model.Service) {
	insert(COLLECTION_SERVICE, service)
}

// 这边主要是先针对State。监控报警可能会调用该接口。先不管
func UpdateServiceById(id string, service model.Service) {
	// TODO
}

func DeleteServiceById(id string) {
	delById(COLLECTION_SERVICE, id)
}
