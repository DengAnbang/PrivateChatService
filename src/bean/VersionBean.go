package bean

type VersionBean struct {
	VersionCode int32  `json:"version_code"`
	VersionName string `json:"version_name"`
	VersionMsg  string `json:"version_msg"`
	Packages    string `json:"packages"`
}
