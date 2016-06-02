package marathon

// ================= /v2/apps begin =================
// 请求"POST/v2/apps"
type MarathonAppsRequest struct {
	Id        string            `json:"id,omitempty"`
	Instances int               `json:"instances,omitempty"`
	Cpus      float64           `json:"cpus,omitempty"`
	Mem       float64           `json:"mem,omitempty"`
	Disk      float64           `json:"disk,omitempty"`
	Version   string            `json:"version,omitempty"`
	Container MarathonContainer `json:"container,omitempty"`
	Labels    map[string]string `json:"labels,omitempty"`
}

// 在cloud-server访问marathon时候使用
func NewMarathonAppsRequest() *MarathonAppsRequest {
	var request MarathonAppsRequest
	// 由于采用了ovs+none的方式，所以不需要做端口映射了
	request.Cpus = 0.1
	request.Mem = 2000
	container := NewMarathonContainer()
	request.Container = *container
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

type MarathonContainer struct {
	Type    string                  `json:"type,omitempty"`
	Docker  MarathonDockerContainer `json:"docker,omitempty"`
	Volumes []MarathonPerVolume     `json:"volumes"`
}

type MarathonDockerContainer struct {
	Image        string               `json:"image,omitempty"`
	Network      string               `json:"network,omitempty"`
	PortMappings []MarathonDockerPort `json:"portMappings,omitempty"`
}

type MarathonPerVolume struct {
	ContainerPath string `json:"containerPath,omitempty"`
	HostPath      string `json:"hostPath,omitempty"`
	Mode          string `json:"mode,omitempty"`
}

// 在cloud-server访问marathon时候使用
func NewMarathonContainer() *MarathonContainer {
	var container MarathonContainer
	container.Type = "DOCKER"

	var docker MarathonDockerContainer
	docker.Network = "NONE"
	container.Docker = docker

	container.Volumes = make([]MarathonPerVolume, 1)
	var coreVolume MarathonPerVolume
	coreVolume.ContainerPath = "/home/work/data"
	coreVolume.HostPath = "/data/container"
	coreVolume.Mode = "RW"
	container.Volumes[0] = coreVolume
	return &container
}

type MarathonDockerPort struct {
	ContainerPort int    `json:"containerPort,omitempty"`
	HostPort      int    `json:"hostPort"·`
	ServicePort   int    `json:"servicePort,omitempty"`
	Protocol      string `json:"protocol,omitempty"`
}

// 在cloud-server访问marathon时候使用
func NewMarathonDockerPort() *MarathonDockerPort {
	var port MarathonDockerPort
	// port.Protocol = "tcp"
	port.HostPort = 0
	return &port
}

// ================= /v2/apps end =================
