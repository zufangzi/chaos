package main

import (
	"flag"
	"fmt"
	consulApi "github.com/hashicorp/consul/api"
	"github.com/hashicorp/go-cleanhttp"
	"log"
	"net/http"
	"opensource/chaos/common"
	"opensource/chaos/common/fasthttp"
	"opensource/chaos/domain/model/mongo"
	"os"
	"strconv"
	"time"
)

var consulClient *consulApi.Client
var signFile string
var serviceId, serviceName, servicePort, serviceIp string
var serviceTags []string
var isNormal bool
var fatalExit chan int
var monitorInterval = flag.Int("monitor-interval", 2000, "Interval(in millisecond) between monitor attemps")
var ARGS_CLOUD_SERVER_URL = flag.String("cloud-server", "10.32.27.76:8080", "Cloud Server Address Without 'http://' prefix")

var PATH_CONTAINER_STATE = "http://" + *ARGS_CLOUD_SERVER_URL + "/api/containers/%s/%d"

func initEnv() {
	// 让命令行的命令生效
	flag.Parse()
	// 初始化一些初始变量
	serviceIp = common.GetHostIp()
	serviceName = os.Getenv("SERVICE_NAME")
	common.Assert(serviceName)
	servicePort = os.Getenv("SERVICE_PORT")
	common.Assert(servicePort)
	// 直接以容器ID + ServiceName来命名，确保唯一。如果后续发现容器ID不唯一，那么就再加上宿主机IP
	// serviceId = strconv.FormatInt(time.Now().Unix(), 10) + "_" + serviceName
	serviceId = common.GetHostName() + "_" + serviceName
	signFile = "/home/work/data/$(hostname)/kickoff_sign_file"
	isNormal = false

	// 给一个channel，如果没有严重错误就永久等待
	fatalExit = make(chan int)
}

func defaultConfig() *consulApi.Config {
	// 拿ip最后一位设置为1，即为宿主机ip。默认宿主机上必须有consul
	// Deprecated。采用OVS+none划分VLAN的方式重做二层网络。采用读取信号量文件来拿
	// hostIp := serviceIp[:strings.LastIndex(serviceIp, ".")] + ".1:5000"
	hostIp := common.GetShell("cat "+signFile+" | head -n1") + ":8500"
	// hostIp := "10.14.5.14:8500"
	config := &consulApi.Config{
		Address:    hostIp,
		Scheme:     "http",
		HttpClient: cleanhttp.DefaultClient(),
	}
	return config
}

func main() {

	initEnv()

	// 持续监听指定端口号，采用netstat -nlp | grep ":$PORT"来实现，隔5s来一次
	log.Printf("[REGISTRATOR]Now begin to listen service: %s with service port: %s.", serviceName, servicePort)

	if *monitorInterval <= 0 {
		log.Println("[REGISTRATOR]the param: monitorInterval is less than 0!")
		os.Exit(2)
	}
	// ticker := time.NewTicker(time.Duration(*monitorInterval) * time.Millisecond)
	log.Println("[REGISTRATOR]Now begin to listen the port...")
	for {
		if !common.PortExists(servicePort) { // 说明启动还未成功
			log.Println("[REGISTRATOR]Not found port for registration, please wait...")
			time.Sleep(time.Duration(*monitorInterval) * time.Millisecond)
			continue
		}
		log.Println("[REGISTRATOR]found the port for registration!")
		isNormal = true
		break
	}

	// 启动成功的话，就先初始化consul客户端
	log.Println("[REGISTRATOR]Now begin to process registration!")
	consulClient, _ = consulApi.NewClient(defaultConfig())

	// 然后调用consulapi进行注册
	register()
	updateState(common.GetHostName(), mongo.CONTAINER_STATE_ALL_UP)
	log.Println("[REGISTRATOR]Now finished registration!")

	// 启动自动监控功能，当端口不在的时候，能够取消注册
	ticker := time.NewTicker(time.Duration(*monitorInterval) * time.Second)
	go func() {
		for {
			select {
			case <-ticker.C:
				log.Println("[REGISTRATOR]Now in a loop monitor ticket time!")
				if common.PortExists(servicePort) {
					if !isNormal {
						log.Println("[REGISTRATOR]Now begin to register!")
						isNormal = true
						register()
						// 注册信号量
						updateState(common.GetHostName(), mongo.CONTAINER_STATE_ALL_UP)
					}
				} else {
					if isNormal {
						log.Println("[REGISTRATOR]Now begin to deregister!")
						isNormal = false
						deregister()
						// 取消信号量。如果是容器北山，就已经至少被mask掉了。
						updateState(common.GetHostName(), mongo.CONTAINER_STATE_ERROR)
					}
				}
			}
		}
	}()

	log.Fatalln("[REGISTRATOR]Error...", <-fatalExit)
}

func register() {
	registration := new(consulApi.AgentServiceRegistration)
	registration.ID = serviceId
	registration.Name = serviceName
	registration.Port, _ = strconv.Atoi(servicePort)
	registration.Tags = serviceTags
	registration.Address = serviceIp
	// TODO build Check latter
	log.Println("ID is: ", registration.ID)
	log.Println("Name is: ", registration.Name)
	log.Println("Port is: ", registration.Port)
	log.Println("Tags is: ", registration.Tags)
	log.Println("Address is: ", registration.Address)
	consulClient.Agent().ServiceRegister(registration)
}

func deregister() {
	consulClient.Agent().ServiceDeregister(serviceId)
}

func updateState(cId string, cState int) {
	log.Printf("[REGISTRATOR]now begin to update state, cid: %s, cstate: %s %n", cId, cState)
	var resData map[string]interface{}
	uri := fmt.Sprintf(PATH_CONTAINER_STATE, cId, cState)
	fasthttp.JsonReqAndResHandler(uri, make(map[string]interface{}), &resData, "POST")
	if int(resData["status"].(float64)) != http.StatusOK {
		log.Println("[CHAOSWHISPER] request cloud server fail... CID: " + cId)
	}
}
