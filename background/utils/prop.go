package utils

import (
	"github.com/tjz101/goprop"
	// "log"
	"os"
)

// global var
var Path PathProps
var Param ParamProps

type PathProps struct {

	// marathon
	MarathonAppsUrl   string
	MarathonGroupsUrl string

	// docker registry
	DockerRegistryUrl       string
	DockerRegistrySearchUrl string

	// mongo
	MongoUrl string

	// redis
	RedisUrl string
}

type ParamProps struct {
	MongoDB string
}

func init() {
	if os.Getenv("GOPATH") == "" {
		return
	}
	Path = PathProps{}
	prop := goprop.NewProp()
	prop.Read(os.Getenv("GOPATH") + PROP_FILE)
	Path.MarathonAppsUrl = prop.Get("marathon.apps.url")
	Path.MarathonGroupsUrl = prop.Get("marathon.groups.url")
	Path.DockerRegistryUrl = prop.Get("docker.registry.url")
	Path.DockerRegistrySearchUrl = Path.DockerRegistryUrl + prop.Get("docker.registry.search.url")
	Path.MongoUrl = prop.Get("mongo.url")
	Path.RedisUrl = prop.Get("redis.url")

	Param.MongoDB = prop.Get("mongo.db")
}
