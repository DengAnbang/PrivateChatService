package code

const (
	OK        = "0"  //成功
	NormalErr = "-1" //普通错误
)

var (
	RootName         = "/res/"
	RootPath         = CurrentPath + RootName
	ViewRootPath     = RootPath + "view/"
	LogRootPath      = RootPath + "log/"
	FileRootPath     = RootPath + "file/"
	DatabaseRootPath = RootPath + "database/"
	ConfigRootPath   = RootPath + "config/"
	FileAppPath      = FileRootPath + "app/"
	FileAppPathName  = FileAppPath + "/app.apk"
)
