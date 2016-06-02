package mongo

import (
	"github.com/samalba/dockerclient"
)

const (
	CONTAINER_STATE_WAIT_FOR_SERVICE = iota
	CONTAINER_STATE_STARTED
	CONTAINER_STATE_ALL_UP
	CONTAINER_STATE_ERROR
	CONTAINER_STATE_DELETE
)

const (
	SERVICE_STATE_PRE_DEPLOYMENT = iota
	SERVICE_STATE_IN_DEPLOYMENT
	SERVICE_STATE_ALL_UP_HEALTHY
	SERVICE_STATE_UNHEALTHY
)

const (
	GROUP_TYPE_ONLINE = iota
	GROUP_TYPE_PRE_TEST
	GROUP_TYPE_TEST
	GROUP_TYPE_VM
	GROUP_TYPE_PRE_ONLINE
)

const (
	GROUP_STATE_DEPLOYING = iota
	GROUP_STATE_FAIL
	GROUP_STATE_NORMAL
)

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

type Host struct {
	State    int
	HostIp   string
	Vlan     string
	NameNode bool
}
