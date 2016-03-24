package main

import (
	"flag"
	"fmt"
	"log"
	"opensource/chaos/background/utils"
	"os"
	"time"
)

var signPath, signFile string
var monitorInterval = flag.Int("monitor-interval", 2, "Interval(in millisecond) between monitor attemps")
var SCRIPT_WGET_HOME = "http://10.32.27.11/cc/preonline/container/"

func initEnv() {
	flag.Parse()
	signPath = "/home/work/data/" + utils.GetHostName() + "/"
	signFile = "kickoff_sign_file"
}

func main() {
	initEnv()
	log.Printf("[GUARDKEEPER]Now begin to start guardkeeper...")

	if *monitorInterval <= 0 {
		log.Println("[GUARDKEEPER]the param: monitorInterval is less than 0!")
		os.Exit(2)
	}
	for {
		if !utils.FileExists(signPath, signFile) {
			log.Println("[GUARDKEEPER]sign file not found yet, path: " + signPath + ", filename: " + signFile + " please wait...")
			time.Sleep(time.Duration(*monitorInterval) * time.Second)
			continue
		}
		log.Println("[GUARDKEEPER]found the sign file!")
		break
	}

	// hostIp := utils.GetHostIp()
	hostName := utils.GetHostName()
	realIp := utils.GetShell("ifconfig eth0 | grep \"inet addr\" | awk '{print $2}' | awk -F: '{print $2}'")
	// 选用--net=bridge方式的话，则进行ip替换。
	// commandCore := fmt.Sprintf("s/%s\t%s/%s\t%s/g", hostIp, hostName, realIp, hostName)
	// utils.GetShell("sed \"" + commandCore + "\" /etc/hosts > out.tmp && cat out.tmp > /etc/hosts && rm -f out.tmp")

	// 选用--net=none方式的话，那么/etc/hosts里面不会有ip\tport信息。此时直接往里面插就行。
	utils.GetShell(fmt.Sprintf("echo \"%s\t%s\" >> /etc/hosts", realIp, hostName))

	log.Println("[GUARDKEEPER]now begin to execute supervisord!")
	log.Fatalln("[GUARDKEEPER]Error...", utils.GetShell("/usr/bin/supervisord"))
}
