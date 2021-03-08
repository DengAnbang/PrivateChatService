package code

const (
	OK        = 0  //成功
	NormalErr = -1 //普通错误
)
const (
	AdministratorAccount = "admin"
	AdministratorPwd     = "111111"
	AdministratorName    = "管理员"
)

var (
	RootName         = "/res/"
	RootPath         = CurrentPath + RootName
	ViewRootPath     = RootPath + "view/"
	LogRootPath      = RootPath + "log/"
	FileRootPath     = RootPath + "file/"
	DatabaseRootPath = RootPath + "database/"
	ConfigRootPath   = RootPath + "config/"
)

//socket相关常量
const (
	TypeConnect    = 1     //连接成功
	TypeOtherLogin = 10001 //其他人登陆

	TypePushAskForLeave    = 10002 //请假消息的推送
	TypePushCommissioned   = 10003 //委托消息的推送
	TypePushCreateWorksite = 10004 //申请创建工点的推送

	TypePushChat = 11001 //聊天消息的推送
)
