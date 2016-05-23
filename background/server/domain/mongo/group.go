package mongo

import (
	// "fmt"
	"opensource/chaos/background/server/dto/model"
	// "opensource/chaos/background/utils"
)

const COLLECTION_GROUP = "cgroup"

func GetGroupInfoById(id string) model.Group {
	result := model.Group{}
	queryById(COLLECTION_GROUP, &result, id)
	return result
}

func InsertGroup(group model.Group) {
	insert(COLLECTION_GROUP, group)
}

func UpdateGroupById(id string, group model.Group) {
	// TODO
}

func DeleteGroupById(id string) {
	delById(COLLECTION_GROUP, id)
}
