package dto

type DeployAppsSimpleRequest struct {
	Id          string
	ExportPorts []ExportPort
	Cpus        string
	Mem         string
	Instances   string
	Image       string
	Version     string
}

type ExportPort struct {
	ContainerPort string
	Desc          string
}

type DeployGroupsSimpleRequest struct {
	Id     string
	Groups []GroupsInfo
}

type GroupsInfo struct {
	Id   string
	Apps []DeployAppsSimpleRequest
}

// 使用的时候提供的接口提供批量接口
type RollbackAppsRequest struct {
	Id      string
	Version string
}
