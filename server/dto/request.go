package dto

type DeploySimpleRequest struct {
	Id          string
	ExportPorts []ExportPort
	Cpus        string
	Mem         string
	Instances   string
	Image       string
}

type ExportPort struct {
	ContainerPort string
	Desc          string
}
