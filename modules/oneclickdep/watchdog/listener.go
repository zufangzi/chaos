package oneclickdep

import (
	"fmt"
)

// 启动之后就开始轮询
func Sentinel() {
	// 启动自动监控功能，当端口不在的时候，能够取消注册
	ticker := time.NewTicker(time.Duration(*monitorInterval) * time.Second)
	go func() {
		for {
			select {
			case <-ticker.C:
				log.Println("[SENTINEL]Now in a loop ticket time!")
				// 先找出所有可以被部署的服务
				services := lookupForReadyServices()
				// 针对以上所有service去部署
				deployServices()
			}
		}
	}()
}

// 遍历Group的所有Service，找到Dependency已经都ok的所有Service，进行部署
func lookupForReadyServices() []string {
	return nil
}

func deployServices() {

}
