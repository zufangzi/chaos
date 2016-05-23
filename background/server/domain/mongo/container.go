package mongo

import (
	// "fmt"
	"opensource/chaos/background/server/dto/model"
	// "opensource/chaos/background/utils"
)

const COLLECTION_CONTAINER = "container"

func GetContainerInfoById(id string) model.Container {
	result := model.Container{}
	queryById(COLLECTION_CONTAINER, &result, id)
	return result
}

func InsertContainer(container model.Container) {
	insert(COLLECTION_CONTAINER, container)
}

// 获取虚拟IP时候，不但要返回，而且还要写入DB
func UpdateContainerByIdUsingMap(id string, data map[string]int) {
	updateById(COLLECTION_CONTAINER, &data, id)
}

// 获取虚拟IP时候，不但要返回，而且还要写入DB
func UpdateContainerById(id string, container model.Container) {
	updateById(COLLECTION_CONTAINER, &container, id)
}

func DeleteContainerById(id string) {
	delById(COLLECTION_CONTAINER, id)
}
