package main

import (
	"errors"
	"flag"
	// "fmt"
	consulApi "github.com/hashicorp/consul/api"
	"github.com/hashicorp/go-cleanhttp"
	"log"
	"os"
	"os/exec"
	"strconv"
	"strings"
	"time"
)

var consulClient *consulApi.Client
var serviceId, serviceName, servicePort, serviceIp string
var serviceTags []string
var isNormal bool
var fatalExit chan int
var monitorInterval = flag.Int("monitor-interval", 2, "Interval(in millisecond) between monitor attemps")

func assert(param interface{}) {
	if param == nil {
		os.Exit(2)
	}
	switch param.(type) {
	case string:
		if param == "" {
			os.Exit(2)
		}
	}
}

func initEnv() {
	// 让命令行的命令生效
	flag.Parse()
	// 初始化一些初始变量
	serviceIp = getShell("hostname -i")
	serviceName = os.Getenv("SERVICE_NAME")
	assert(serviceName)
	servicePort = os.Getenv("SERVICE_PORT")
	assert(servicePort)
	serviceId = strconv.FormatInt(time.Now().Unix(), 10) + "_" + serviceName
	isNormal = false

	// 给一个channel，如果没有严重错误就永久等待
	fatalExit = make(chan int)
}

func defaultConfig() *consulApi.Config {
	// 拿ip最后一位设置为1，即为宿主机ip。默认宿主机上必须有consul
	hostIp := serviceIp[:strings.LastIndex(serviceIp, ".")] + ".1:5000"
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
		if !monitor() { // 说明启动还未成功
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
	log.Println("[REGISTRATOR]Now finished registration!")

	// 启动自动监控功能，当端口不在的时候，能够取消注册
	ticker := time.NewTicker(time.Duration(*monitorInterval) * time.Second)
	go func() {
		for {
			select {
			case <-ticker.C:
				log.Println("[REGISTRATOR]Now in a loop monitor ticket time!")
				if monitor() {
					if !isNormal {
						log.Println("[REGISTRATOR]Now begin to register!")
						isNormal = true
						register()
					}
				} else {
					if isNormal {
						log.Println("[REGISTRATOR]Now begin to deregister!")
						isNormal = false
						deregister()
					}
				}
			}
		}
	}()

	log.Fatalln("[REGISTRATOR]Error...", <-fatalExit)
}

func getShell(cmdStr string) string {
	cmd := exec.Command("/bin/sh", "-c", cmdStr)
	out, err := cmd.Output()
	if err != nil {
		panic(errors.New("error occur while process shell cmd"))
	}
	cleanOut := string(out)
	realOut := cleanOut[0:strings.Index(cleanOut, "\n")]
	return realOut
}

func monitor() bool {
	cnt, _ := strconv.Atoi(getShell("netstat -nlp | grep \":" + servicePort + "\"|wc -l"))
	log.Println("[REGISTRATOR]Now found counts: ", cnt)
	return cnt > 0
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
