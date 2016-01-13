package entity

type MarathonAppsInfo struct {
	Apps []MarathonPerAppInfo
}

type MarathonPerAppInfo struct {
	Id          string
	Instances   int
	Cpus        float32
	Mem         int
	Disk        int
	Version     string
	VersionInfo VersionInfos
	Container   MarathonContainer
}

type MarathonContainer struct {
	Type   string
	Docker MarathonDocker
}

type MarathonDocker struct {
	Image        string
	Network      string
	PortMappings []MarathonDockerPort
}

type MarathonDockerPort struct {
	ContainerPort int
	HostPort      int
	ServicePort   int
	Protocol      string
}

type VersionInfos struct {
	LastScalingAt      string
	LastConfigChangeAt string
}
