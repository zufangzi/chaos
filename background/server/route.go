package main

import (
	"encoding/json"
	"github.com/ant0ine/go-json-rest/rest"
	"io/ioutil"
	"log"
	"net/http"
	"opensource/chaos/background/server/dto"
	"opensource/chaos/background/server/handler"
	webUtils "opensource/chaos/background/server/utils"
)

func Route() (rest.App, error) {

	return rest.MakeRouter(
		// 新增某个服务
		rest.Post("/apps", restGuarder(handler.CreateAppsHandler)),
		// 获取所有服务的基本信息
		rest.Get("/apps", restGuarder(handler.GetInfoAppsHandler)),
		// TODO 获取某一个服务的详细信息
		rest.Get("/apps/*appId", restGuarder(handler.GetSingleAppsHandler)),
		// 删除某一个服务的所有实例
		rest.Delete("/apps/*appId", restGuarder(handler.DeleteAppsHandler)),
		// 新增或者更新一批服务
		rest.Post("/apps/updater", restGuarder(handler.CreateOrUpdateAppsHandler)),
		// 回滚一批服务
		rest.Post("/apps/rollback", restGuarder(handler.RollbackAppsHandler)),
		// 新增一批组信息
		rest.Post("/groups", restGuarder(handler.DeployGroupsHandler)),
	)
}

func restGuarder(method RestFunc) rest.HandlerFunc {
	return func(w rest.ResponseWriter, r *rest.Request) {
		// begin := time.Now().UnixNano()
		defer func() {
			if e, ok := recover().(error); ok {
				rest.Error(w, e.Error(), http.StatusInternalServerError)
				log.Println("catchable system error occur: ", e)
			}
			// log.Printf("the request: %s cost: %d ms\n", r.URL.RequestURI(), ((time.Now().UnixNano() - begin) / 1000000))
		}()

		pathParams := r.PathParams
		var request dto.CommonRequest
		content, err := ioutil.ReadAll(r.Body)
		defer r.Body.Close()
		webUtils.CheckError(err)
		if len(content) == 0 {
			w.WriteJson(method(pathParams, nil))
			return
		}
		err = json.Unmarshal(content, &request)
		webUtils.CheckError(err)
		switch request.SyncType {
		case "sync":
			log.Println("now use sync mode")
			w.WriteJson(method(pathParams, content))
		case "async":
			log.Println("now use async mode")
			go method(pathParams, content)
			w.WriteJson(map[string]string{"status": "ok"})
		default:
			log.Println("now use default mode(sync)", request.SyncType)
			w.WriteJson(method(pathParams, content))

		}

	}
}

type RestFunc func(pathParams map[string]string, data []byte) interface{}
