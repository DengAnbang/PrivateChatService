package socketConst

//socket相关常量
const (
	TypeConnect            = "1"     //连接成功
	TypeOtherLogin         = "10001" //其他人登陆
	TypeLogin              = "10002" //登录
	TYPE_MSG_STATUS_SEND   = "20000" //状态,表示对方收到
	TYPE_MSG               = "20001" //消息
	TypePushAskForLeave    = 10002   //请假消息的推送
	TypePushCommissioned   = 10003   //委托消息的推送
	TypePushCreateWorksite = 10004   //申请创建工点的推送

	TypePushChat = 11001 //聊天消息的推送
)
