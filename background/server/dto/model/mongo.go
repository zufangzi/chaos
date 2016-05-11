package model

import (
	"github.com/samalba/dockerclient"
)

const (
	CONTAINER_STATE_WAIT_FOR_SERVICE = itoa
	CONTAINER_STATE_ALL_UP
	CONTAINER_STATE_ERROR
)

const (
	SERVICE_STATE_HEALTHY = itoa
	SERVICE_STATE_UNHEALTHY
)

const (
	GROUP_TYPE_ONLINE = itoa
	GROUP_TYPE_PRE_TEST
	GROUP_TYPE_TEST
	GROUP_TYPE_VM
	GROUP_TYPE_PRE_ONLINE
)

const (
	GROUP_STATE_DEPLOYING = itoa
	GROUP_STATE_FAIL
	GROUP_STATE_NORMAL
)

const itoa = 0

// 不能用组合方式来做，mongo那边做insert和query时候会
// 错乱，将mongo字段带入作为根
type Mongo struct {
	State   int
	BizId   string
	BizName string
}

type Container struct {
	State       int
	BizId       string
	BizName     string
	ContainerIp string
	ServiceId   string
	Inspect     dockerclient.ContainerInfo
}

type Service struct {
	State        int
	BizId        string
	BizName      string
	Dependencies []string
	GroupId      string
	Instances    int
	Cpu          float64
	Mem          float64
	Disk         float64
	MarathonId   string
}

type Group struct {
	State      int
	BizId      string
	BizName    string
	ServiceCnt int
	Type       int
}
