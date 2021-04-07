package socket

import (
	"encoding/json"
	"fmt"
	"gitee.com/DengAnbang/PrivateChatService/src/bean"
	"gitee.com/DengAnbang/PrivateChatService/src/socket/push"
	"gitee.com/DengAnbang/PrivateChatService/src/socket/socketConst"
	"gitee.com/DengAnbang/goutils/loge"
	"github.com/gorilla/websocket"
	"log"
	"net/http"
)

func WebSocketRun(w http.ResponseWriter, r *http.Request) {
	upgrader := websocket.Upgrader{
		// 解决跨域问题
		CheckOrigin: func(r *http.Request) bool {
			return true
		},
	}
	c, err := upgrader.Upgrade(w, r, nil)
	if err != nil {
		log.Print("upgrade:", err)
		return
	}
	err = webSocketHandler(c)
	if err != nil {
		log.Print("upgrade:", err)
		return
	}
}
func webSocketHandler(ws *websocket.Conn) error {
	defer ws.Close()
	conn := WsConn{Conn: ws}
	message := bean.NewSucceedMessage("连接成功!")
	message.Type = socketConst.TypeConnect
	err := conn.SendMessageToConn(message)
	if err != nil {
		return err
	}
	for {
		if _, message, err := ws.ReadMessage(); err != nil {
			loge.W("receive failed:", err.Error())
			if _, ok := err.(*websocket.CloseError); ok {
				push.UnRegister(&conn)
				return err
			}
			break
		} else {
			var sm bean.SocketData
			err := json.Unmarshal(message, &sm)
			s := string(message)
			loge.W("读取到的消息:" + s)
			if err == nil {
				go Dispense(&sm, &conn)
			} else {
				conn.Response(bean.NewErrorMessage(fmt.Sprintf("编码错误,%s", err)), "0")
			}
		}

	}
	return nil
}
