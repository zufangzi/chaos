package utils

import (
	"errors"
	"os/exec"
	"strconv"
	"strings"
)

func GetShell(cmdStr string) string {
	cmd := exec.Command("/bin/sh", "-c", cmdStr)
	out, err := cmd.Output()
	if err != nil {
		panic(errors.New("error occur while process shell cmd"))
	}
	cleanOut := string(out)
	realOut := cleanOut[0:strings.Index(cleanOut, "\n")]
	return realOut
}

func GetHostName() string {
	return GetShell("hostname")
}

func GetHostIp() string {
	return GetShell("hostname -ip")
}

func PortExists(servicePort string) bool {
	cnt, _ := strconv.Atoi(GetShell("netstat -nlp | grep \":" + servicePort + "\"|wc -l"))
	return cnt > 0
}

func FileExists(signPath, signFile string) bool {
	cnt, _ := strconv.Atoi(GetShell("find " + signPath + " -name " + signFile + " | wc -l"))
	return cnt > 0
}
