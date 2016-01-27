// TODO 对于response需要进行struct的统一。code/message/data
package dto

type AppsGlobalInfoResponse struct {
	Id               string
	Instances        string
	Cpus             string
	Mem              string
	CurrentInstances string
	Healthy          string
	Group            string
	Status           string
}

var statusInfoMap map[int]string

func init() {
	statusInfoMap = map[int]string{
		0: "正常",
	}
}

func (info *AppsGlobalInfoResponse) FormatStatus(stage int) {
	info.Status = statusInfoMap[stage]
}
