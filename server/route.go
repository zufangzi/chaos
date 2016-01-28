package main

import (
	"github.com/ant0ine/go-json-rest/rest"
	"opensource/chaos/server/handler"
	"opensource/chaos/server/utils"
)

func Route() (rest.App, error) {
	return rest.MakeRouter(
		// 新增某个服务
		rest.Post("/apps", utils.RestGuarder(handler.CreateAppsHandler)),
		// 获取所有服务的基本信息
		rest.Get("/apps", utils.RestGuarder(handler.GetInfoAppsHandler)),
		// TODO 获取某一个服务的详细信息
		rest.Get("/apps/#appId", utils.RestGuarder(handler.GetSingleAppsHandler)),
		// TODO 删除某一个服务的所有实例
		rest.Delete("/apps/#appId", utils.RestGuarder(handler.DeleteAppsHandler)),
		// 新增或者更新一批服务
		rest.Post("/apps/updater", utils.RestGuarder(handler.CreateOrUpdateAppsHandler)),
		// 回滚一批服务
		rest.Post("/apps/rollback", utils.RestGuarder(handler.RollbackAppsHandler)),
		// 新增一批组信息
		rest.Post("/groups", utils.RestGuarder(handler.DeployGroupsHandler)),
	)
}
