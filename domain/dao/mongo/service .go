package mongo

import (
	"opensource/chaos/domain/model/mongo"
)

const COLLECTION_SERVICE = "service"

func GetServiceInfoById(id string) mongo.Service {
	result := mongo.Service{}
	queryById(COLLECTION_SERVICE, &result, id)
	return result
}

func InsertService(service mongo.Service) {
	insert(COLLECTION_SERVICE, service)
}

// 这边主要是先针对State。监控报警可能会调用该接口。先不管
func UpdateServiceById(id string, service mongo.Service) {
	// TODO
}

func DeleteServiceById(id string) {
	delById(COLLECTION_SERVICE, id)
}
