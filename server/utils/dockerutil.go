package utils

import (
	"fmt"
	"github.com/samalba/dockerclient"
	"opensource/chaos/server/dto/docker"
	"opensource/chaos/server/utils/fasthttp"
	"sort"
	"strconv"
	"strings"
)

var (
	DockerClient = new(FastDocker)
)

type FastDocker struct {
	dockerclient.DockerClient
}

// 回滚时候使用
// 此处希望tag即每次上线版本号是按照"毫秒时间戳_CommitId_Job号"组成
// 则获取前一次的镜像即变成去找Job号倒数第二大的即可。
// 但是现在默认采用的是纯粹按照时间戳来打tag。这样的话只需单纯进行排序取倒数第二个即可。
func (f *FastDocker) GetPreviousImageAndTag(repository string, version string, tagStyle string) (string, string, string) {
	return f.GetImageAndTagByFreshness(repository, version, tagStyle, 1, true)
}

func (f *FastDocker) GetImageAndTagByFreshness(repository string, version string, tagStyle string, previous int, needHttpPrefix bool) (string, string, string) {
	if version != "" {
		return sweepHttpPrefix(Path.DockerRegistryUrl+"/"+repository+":"+version, needHttpPrefix)
	}

	uri := fmt.Sprintf(Path.DockerRegistrySearchUrl, repository)
	var res docker.DockerRegistryTagsResponse

	fasthttp.JsonReqAndResHandler(uri, nil, &res, "GET")

	if tagStyle == "" {
		tagStyle = IMAGES_TAG_STYLE_SIMPLE
	}

	if len(res.Tags) == 1 {
		return sweepHttpPrefix(Path.DockerRegistryUrl+"/"+repository+":"+res.Tags[0], needHttpPrefix)
	}

	tagArray := make([]int, len(res.Tags))
	for i, v := range res.Tags {
		tagArray[i], _ = strconv.Atoi(v)
	}

	sort.Ints(tagArray)
	fullImage := Path.DockerRegistryUrl + "/" + repository + ":" + strconv.Itoa(tagArray[len(tagArray)-1-previous])
	return sweepHttpPrefix(fullImage, needHttpPrefix)
}

func sweepHttpPrefix(fullImage string, needHttpPrefix bool) (string, string, string) {
	cutFullImage := fullImage[strings.Index(fullImage, "http://")+7:]
	if !needHttpPrefix {
		fullImage = cutFullImage
	}

	tag := cutFullImage[strings.LastIndex(cutFullImage, ":")+1:]
	image := cutFullImage[:strings.LastIndex(cutFullImage, ":")]
	return fullImage, image, tag
}
