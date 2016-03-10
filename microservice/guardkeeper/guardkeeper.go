package main

import (
	"flag"
	// "fmt"
	"log"
	"opensource/chaos/microservice/utils"
	"os"
	"time"
)

var signPath, signFile string
var monitorInterval = flag.Int("monitor-interval", 2, "Interval(in millisecond) between monitor attemps")

func initEnv() {
	flag.Parse()
	signPath = "/home/work"
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
			log.Println("[GUARDKEEPER]sign file not found yet, please wait...")
			time.Sleep(time.Duration(*monitorInterval) * time.Millisecond)
			continue
		}
		log.Println("[GUARDKEEPER]found the sign file!")
		break
	}

	log.Println("[GUARDKEEPER]now begin to execute supervisord!")
	log.Fatalln("[GUARDKEEPER]Error...", utils.GetShell("/usr/bin/supervisord"))
}
