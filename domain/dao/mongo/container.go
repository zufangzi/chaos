package mongo

import (
	// "fmt"
	"opensource/chaos/domain/model/mongo"
	// "opensource/chaos/background/utils"
)

const COLLECTION_CONTAINER = "container"

func GetContainerInfoById(id string) mongo.Container {
	result := mongo.Container{}
	queryById(COLLECTION_CONTAINER, &result, id)
	return result
}

func InsertContainer(container mongo.Container) {
	insert(COLLECTION_CONTAINER, container)
}

// 获取虚拟IP时候，不但要返回，而且还要写入DB
func UpdateContainerByIdUsingMap(id string, data map[string]int) {
	updateById(COLLECTION_CONTAINER, &data, id)
}

// 获取虚拟IP时候，不但要返回，而且还要写入DB
func UpdateContainerById(id string, container mongo.Container) {
	updateById(COLLECTION_CONTAINER, &container, id)
}

func DeleteContainerById(id string) {
	delById(COLLECTION_CONTAINER, id)
}
