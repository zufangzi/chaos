package test

import (
	"fmt"
	"opensource/chaos/background/server/domain/mongo"
	"opensource/chaos/background/server/dto/model"
	"testing"
)

func Test_mongo_service_simple_insert_and_query(t *testing.T) {
	service := model.Service{}
	service.BizId = "9228817231"
	service.BizName = "helloworld"
	service.Dependencies = []string{"a", "b", "c"}
	mongo.InsertService(service)
	fmt.Println(mongo.GetServiceInfoById(service.BizId).Dependencies)
	if len(mongo.GetServiceInfoById(service.BizId).Dependencies) != 3 {
		t.Fail()
	}
	mongo.DeleteServiceById(service.BizId)
}

func Test_mongo_group_simple_insert_and_query(t *testing.T) {
	group := model.Group{}
	group.BizId = "189182cef22sdf"
	group.BizName = "hello-group"
	group.ServiceCnt = 4
	mongo.InsertGroup(group)
	fmt.Println(mongo.GetGroupInfoById(group.BizId).ServiceCnt)
	if mongo.GetGroupInfoById(group.BizId).ServiceCnt != 4 {
		t.Fail()
	}
	mongo.DeleteGroupById(group.BizId)
}
