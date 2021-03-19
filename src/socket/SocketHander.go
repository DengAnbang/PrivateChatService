package socket

import (
	"fmt"
	"gitee.com/DengAnbang/PrivateChatService/src/bean"
	"gitee.com/DengAnbang/PrivateChatService/src/socket/push"
	"gitee.com/DengAnbang/PrivateChatService/src/socket/socketConst"
)

func Dispense(message *bean.SocketData, conn push.ResponseAble) {
	defer func() {
		i := recover()
		if err, ok := i.(error); ok {
			conn.Response(err, "0")
		}
	}()

	//err := json.Unmarshal(message.Data, &parameters)
	//panic(err)
	switch message.Type {
	//用户相关
	case socketConst.TypeLogin: //登录
		userId := message.SenderId
		lastCon := push.Register(userId, conn)
		if lastCon != nil && lastCon != conn {
			succeedMessage := bean.NewSucceedMessage("账号在其他地方登录了!")
			succeedMessage.Type = socketConst.TypeOtherLogin
			_ = lastCon.SendMessageToConn(succeedMessage)
		}
		succeedMessage := bean.NewSucceedMessage("登录成功!")
		succeedMessage.Type = socketConst.TypeLogin
		_ = conn.SendMessageToConn(succeedMessage)

	case socketConst.TYPE_MSG: //消息
		targetId := message.TargetId
		message := message.Data
		push.Push(targetId, socketConst.TYPE_MSG, message)

	case socketConst.TYPE_MSG_STATUS_SEND:
		targetId := message.TargetId
		senderId := message.SenderId
		push.Push(senderId, socketConst.TYPE_MSG_STATUS_SEND, targetId)

	default:
		conn.Response(bean.NewErrorMessage(fmt.Sprintf("未知的消息类型%v", message.Type)), "0")
	}
}
func GetString(parameters map[string]string, key string) string {
	return parameters[key]

}

//func GetString(parameters map[string]interface{}, key string) string {
//	i := parameters[key]
//	if i == nil {
//		return ""
//	}
//	return fmt.Sprintf("%s", i)
//}