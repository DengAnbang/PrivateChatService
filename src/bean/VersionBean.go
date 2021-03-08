package bean

type VersionBean struct {
	FileBean
	VersionCode    string `json:"version_code"`
	VersionName    string `json:"version_name"`
	VersionMsg     string `json:"version_msg"`
	VersionChannel string `json:"version_channel"`
}

func NewVersionBean(data map[string]string) VersionBean {
	return VersionBean{
		FileBean:       NewFileBean(data),
		VersionCode:    data["version_code"],
		VersionName:    data["version_name"],
		VersionMsg:     data["version_msg"],
		VersionChannel: data["version_channel"],
	}
}
