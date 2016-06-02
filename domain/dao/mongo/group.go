package mongo

import (
	// "fmt"
	"opensource/chaos/domain/model/mongo"
	// "opensource/chaos/background/utils"
)

const COLLECTION_GROUP = "cgroup"

func GetGroupInfoById(id string) mongo.Group {
	result := mongo.Group{}
	queryById(COLLECTION_GROUP, &result, id)
	return result
}

func InsertGroup(group mongo.Group) {
	insert(COLLECTION_GROUP, group)
}

func UpdateGroupById(id string, group mongo.Group) {
	// TODO
}

func DeleteGroupById(id string) {
	delById(COLLECTION_GROUP, id)
}
