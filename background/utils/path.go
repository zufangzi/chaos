package utils

import (
	"github.com/tjz101/goprop"
	"log"
	"os"
)

// global var
var Path PathProps

type PathProps struct {
	MarathonAppsUrl         string
	MarathonGroupsUrl       string
	DockerRegistryUrl       string
	DockerRegistrySearchUrl string
}

func init() {
	log.Println("hi iam now in path...")
	if os.Getenv("GOPATH") == "" {
		return
	}
	log.Println("hi iam now in path2...")
	Path = PathProps{}
	prop := goprop.NewProp()
	prop.Read(os.Getenv("GOPATH") + PROP_FILE)
	Path.MarathonAppsUrl = prop.Get("marathon.apps.url")
	Path.MarathonGroupsUrl = prop.Get("marathon.groups.url")
	Path.DockerRegistryUrl = prop.Get("docker.registry.url")
	Path.DockerRegistrySearchUrl = Path.DockerRegistryUrl + prop.Get("docker.registry.search.url")
	log.Println("hi iam now in path3...", Path.DockerRegistryUrl)
}
