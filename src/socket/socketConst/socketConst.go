package socketConst

//socket相关常量
const (
	TypeConnect         = "1"     //连接成功
	TypeHeartbeat       = "2"     //心跳
	TypeOtherLogin      = "10001" //其他人登陆
	TypeLogin           = "10002" //登录
	TYPE_LOGIN_OUT      = "10003" //退出登录
	TYPE_MSG_UPDATE     = "20000" //更新消息状态
	TYPE_MSG_SEND       = "20001" //发送消息
	TYPE_MSG_RECEIVE    = "20002" //接收消息
	TYPE_MSG_GROUP_SEND = "20003" //发送群消息
	TYPE_FRIEND_ADD     = "21001" //有好友申请的推送
)
