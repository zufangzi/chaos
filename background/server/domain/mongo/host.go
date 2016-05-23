package mongo

import (
	"opensource/chaos/background/server/dto/model"
)

const COLLECTION_HOST = "host"

func GetHostInfoById(id string) model.Host {
	result := model.Host{}
	queryById(COLLECTION_HOST, &result, id)
	return result
}
