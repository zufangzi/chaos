package feo

// 部署一个或者多个apps的请求
type DeployAppsBatchRequest struct {
	Batch []DeployAppsRequest
	CommonRequest
}

type DeployAppsRequest struct {
	Id          string
	ExportPorts []ExportPort
	Cpus        string
	Mem         string
	Instances   string
	Image       string
	Version     string
	CommonRequest
}

type ExportPort struct {
	ContainerPort string
	Desc          string
}

// 部署一个或者多个groups的请求
type DeployGroupsBatchRequest struct {
	Batch []DeployGroupsRequest
	CommonRequest
}

type DeployGroupsRequest struct {
	Id     string
	Groups []GroupsInfo
	CommonRequest
}

type GroupsInfo struct {
	Id   string
	Apps []DeployAppsRequest
}

// 回滚一个或者多个apps的请求
type RollbackAppsBatchRequest struct {
	Batch []RollbackAppsRequest
	CommonRequest
}

type RollbackAppsRequest struct {
	Id      string
	Version string
	CommonRequest
}

type CommonRequest struct {
	SyncType string
}
