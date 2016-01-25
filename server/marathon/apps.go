package marathon

// ================= /v2/apps begin =================
// 请求"POST/v2/apps"
type MarathonAppsRequest struct {
	Id        string                 `json:"id,omitempty"`
	Instances int                    `json:"instances,omitempty"`
	Cpus      float64                `json:"cpus,omitempty"`
	Mem       float64                `json:"mem,omitempty"`
	Disk      float64                `json:"disk,omitempty"`
	Version   string                 `json:"version,omitempty"`
	Container map[string]interface{} `json:"container,omitempty"`
}

// only be used in request
func NewMarathonAppsRequest() *MarathonAppsRequest {
	var request MarathonAppsRequest
	request.Cpus = 0.1
	request.Mem = 2000
	request.Container = make(map[string]interface{})
	request.Container["type"] = "DOCKER"
	return &request
}

// 该响应所对应的请求是"GET /v2/apps"
type MarathonAppsGlobalInfoResponse struct {
	Apps []MarathonPerAppInfo
}

type MarathonPerAppInfo struct {
	Id             string                 `json:"id,omitempty"`
	Instances      int                    `json:"instances,omitempty"`
	Cpus           float64                `json:"cpus,omitempty"`
	Mem            float64                `json:"mem,omitempty"`
	Disk           float64                `json:"disk,omitempty"`
	Version        string                 `json:"version,omitempty"`
	TasksRunning   int                    `json:"tasksRunning,omitempty"`
	TasksHealthy   int                    `json:"tasksHealthy,omitempty"`
	TasksUnhealthy int                    `json:"tasksUnhealthy,omitempty"`
	TasksStaged    int                    `json:"tasksStaged,omitempty"`
	VersionInfo    VersionInfos           `json:"versionInfo,omitempty"`
	Container      map[string]interface{} `json:"container,omitempty"`
}

type VersionInfos struct {
	LastScalingAt      string `json:"lastScalingAt,omitempty"`
	LastConfigChangeAt string `json:"lastConfigChangeAt,omitempty"`
}

type MarathonDockerContainer struct {
	Image        string               `json:"image,omitempty"`
	Network      string               `json:"network,omitempty"`
	Type         string               `json:"type,omitempty"`
	PortMappings []MarathonDockerPort `json:"portMappings,omitempty"`
	Volumes      []interface{}        `json:"volumes"`
}

// only be used in request
func NewMarathonDockerContainer() *MarathonDockerContainer {
	var container MarathonDockerContainer
	container.Network = "BRIDGE"
	container.Type = "DOCKER"
	container.Volumes = make([]interface{}, 0)
	return &container
}

type MarathonDockerPort struct {
	ContainerPort int    `json:"containerPort,omitempty"`
	HostPort      int    `json:"hostPort"`
	ServicePort   int    `json:"servicePort,omitempty"`
	Protocol      string `json:"protocol,omitempty"`
}

// only be used in request
func NewMarathonDockerPort() *MarathonDockerPort {
	var port MarathonDockerPort
	// port.Protocol = "tcp"
	port.HostPort = 0
	return &port
}

// ================= /v2/apps end =================
