package mongo

import (
	"opensource/chaos/domain/model/mongo"
)

const COLLECTION_HOST = "host"

func GetHostInfoById(id string) mongo.Host {
	result := mongo.Host{}
	queryById(COLLECTION_HOST, &result, id)
	return result
}
