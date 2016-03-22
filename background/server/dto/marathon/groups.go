package marathon

type MarathonGroupsRequest struct {
	Id     string               `json:"id,omitempty"`
	Groups []MarathonGroupsInfo `json:"groups,omitempty"`
}

type MarathonGroupsInfo struct {
	Id   string                `json:"id,omitempty"`
	Apps []MarathonAppsRequest `json:"apps,omitempty"`
}
