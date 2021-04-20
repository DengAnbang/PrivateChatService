package push

import (
	"gitee.com/DengAnbang/PrivateChatService/src/bean"
	"gitee.com/DengAnbang/PrivateChatService/src/socket/socketConst"
	"gitee.com/DengAnbang/goutils/loge"
	"sync"
)

var SocketManage Manage

func init() {
	SocketManage = Manage{Conns: make(map[string]ResponseAble)}
}

type Manage struct {
	heartbeat   bool                    //心跳
	Lock        sync.RWMutex            //互斥(保证线程安全)
	SocketName  string                  //客户端Socket名称
	MaxLifeTime int64                   //垃圾回收时间
	Conns       map[string]ResponseAble //保存packet的指针[SocketID] = packet
}

func (s Manage) SendMessageToKey(key string, msg interface{}) {
	s.Lock.Lock()
	defer s.Lock.Unlock()
	conn, ok := s.Conns[key]
	if ok {
		_ = conn.SendMessageToConn(msg)
	} else {
		loge.W("发送消息失败,%v不存在", key)
	}
}

type ResponseAble interface {
	GetId() string
	SetId(id string)
	Response(err error, messageType string)
	SendMessageToConn(msg interface{}) (err error)
	Close() (err error)
}

func Register(key string, wc ResponseAble) (last ResponseAble) {
	//SocketManage.Lock.Lock()
	//defer SocketManage.Lock.Unlock()
	if conn, ok := SocketManage.Conns[key]; ok {
		last = conn
		delete(SocketManage.Conns, key)
	}
	SocketManage.Conns[key] = wc
	NotificationOnline(wc.GetId(), true)
	wc.SetId(key)
	return
}
func CheckOnline(key string) (online bool) {
	if _, ok := SocketManage.Conns[key]; ok {
		return true
	}
	return false
}
func UnRegister(wc ResponseAble) (last ResponseAble) {
	//SocketManage.Lock.Lock()
	//defer SocketManage.Lock.Unlock()
	delete(SocketManage.Conns, wc.GetId())
	NotificationOnline(wc.GetId(), false)
	last = wc
	return
}

type UserSelectFriendCallback func(user_id, friend_type string) ([]bean.UserBean, error)

var UserSelectFriendCall UserSelectFriendCallback

func NotificationOnline(user_id string, online bool) {
	beans, _ := UserSelectFriendCall(user_id, "1")
	for _, value := range beans {
		data := bean.SocketData{
			TargetId: value.UserId,
			SenderId: user_id,
			Type:     socketConst.TYPE_FRIEND_CHANGE,
			Msg:      "",
			DebugMsg: "",
			Data:     online,
		}
		PushSocketByTargetId(&data, value.UserId)
	}

}
func Push(userId string, pushType string, msg string) {
	succeedMessage := bean.NewSucceedMessage(msg)
	succeedMessage.Type = pushType
	SocketManage.SendMessageToKey(userId, succeedMessage)
}
func PushSocket(socketData *bean.SocketData) {
	SocketManage.SendMessageToKey(socketData.TargetId, socketData)
}
func PushSocketByTargetId(socketData *bean.SocketData, targetId string) {
	SocketManage.SendMessageToKey(targetId, socketData)
}
