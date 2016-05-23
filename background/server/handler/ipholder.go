package handler

import (
	// "fmt"
	"log"
	"net/http"
	"opensource/chaos/background/server/domain/mongo"
	"opensource/chaos/background/server/domain/redis"
	"opensource/chaos/background/server/dto/model"
	webUtils "opensource/chaos/background/server/utils"
)

var IP_QUEUE = "ip_queue"

func DeleteIpForContainer(pathParams map[string]string, data []byte) interface{} {
	cId := pathParams["cId"]
	container := mongo.GetContainerInfoById(cId)
	if container.ContainerIp == "" {
		return webUtils.ProcessResponse(http.StatusNotFound, "")
	}
	// 将ip重新push到队列里面去
	ok := redis.Lpush(IP_QUEUE, container.ContainerIp)
	if !ok {
		log.Fatalln("[CHAOS]push to queue error!")
	}
	return webUtils.ProcessResponse(http.StatusOK, "")
}

func CreateIpForContainer(pathParams map[string]string, data []byte) interface{} {
	cId := pathParams["cId"]
	// 查一下库里面是不是已经有了。有了的话直接返回。
	container := mongo.GetContainerInfoById(cId)
	if container.ContainerIp != "" {
		return webUtils.ProcessResponseFully(http.StatusOK, container.ContainerIp, false)
	}
	// 如果库里没有，那么就从redis的队列里面拿一个，然后存到mongo里，最后返回
	ip := redis.Rpop(IP_QUEUE)
	if ip != "" {
		return saveAndReturnIp(container, ip)
	}
	// 如果库里没有，而且redis队列也空了，那么就触发一次队列更新。如果队列没有，那么返回异常code。由请求方自行处理。
	if !triggerRedisQueueUpdated() {
		return webUtils.ProcessResponse(http.StatusNotFound, "")
	}
	// 如果有，那么就开始处理。可能存在高并发的场景，刷完之后又被其他请求先抢了。这里就不加锁了，由客户端来判断，如果内容
	ip = redis.Rpop(IP_QUEUE)
	if ip == "" {
		return webUtils.ProcessResponse(http.StatusNotFound, "")
	}
	// 保存ip并返回
	return saveAndReturnIp(container, ip)
}

func saveAndReturnIp(container model.Container, ip string) interface{} {
	container.ContainerIp = ip
	mongo.UpdateContainerById(container.BizId, container)
	return webUtils.ProcessResponseFully(http.StatusOK, ip, false)
}

func triggerRedisQueueUpdated() bool {
	// TODO
	return true
}
