package oneclickdep

import (
	"fmt"
)

func init(){
	fixDataPerLoops()
	AutoFixData()
	AutoTopologyDeploy()
	AutoUpdateServiceAndGroupState()
}

// 启动之后就开始轮询,事务问题由AutoFixData来做
func AutoTopologyDeploy() {
	ticker := time.NewTicker(time.Duration(*monitorInterval) * time.Second)
	go func() {
		for {
			select {
			case <-ticker.C:
				log.Println("[SENTINEL]Now in a loop ticket time!")
				// 先找出所有可以被部署的服务
				services := lookupForReadyServices()
				// 针对以上所有service去部署，并修改状态表示在部署中
				deployServices(services)
			}
		}
	}()
}

func AutoUpdateServiceAndGroupState(){
	// TODO
}

func AutoFixData(){
	// TODO
}

func fixDataPerLoops(){

}

// 遍历Group的所有Service，找到Dependency已经都ok的所有Service，进行部署
func lookupForReadyServices() []string {
	// TODO
	return nil
}

func deployServices(services []string) {
	// TODO
}
