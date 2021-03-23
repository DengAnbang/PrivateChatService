package socket

import (
	"encoding/json"
	"fmt"
	"gitee.com/DengAnbang/PrivateChatService/src/bean"

	"gitee.com/DengAnbang/goutils/loge"
	"github.com/gorilla/websocket"
)

type WsConn struct {
	Id   string
	Conn *websocket.Conn
}

func (conn *WsConn) Response(err error, messageType string) {
	if resultData, ok := err.(*bean.ResultData); ok {
		resultData.Type = messageType
		_ = conn.SendMessageToConn(resultData)
	} else if err, ok := err.(error); ok {
		data := bean.NewErrorMessage("服务器内部错误")
		data.DebugMessage = fmt.Sprintf("%v", err)
		_ = conn.SendMessageToConn(data)
	}
}
func (conn *WsConn) SetId(id string) {
	conn.Id = id
}
func (conn *WsConn) GetId() string {
	return conn.Id
}
func (conn *WsConn) SendMessageToConn(msg interface{}) (err error) {
	bytes, err := json.Marshal(msg)
	if err != nil {
		return err
	}
	loge.W("ws:发送的内容:", string(bytes))
	err = conn.Conn.WriteMessage(1, bytes)
	return err
}
