package main

import (
	"fmt"
	"gitee.com/DengAnbang/PrivateChatService/src/api"
	"gitee.com/DengAnbang/PrivateChatService/src/config"
	"gitee.com/DengAnbang/PrivateChatService/src/socket"
	"gitee.com/DengAnbang/goutils/loge"
	"gitee.com/DengAnbang/goutils/sysUtils"
	"gitee.com/DengAnbang/goutils/timeUtils"
	"net/http"
	"time"
)

//http://47.108.172.20:9090/upload
//http://192.168.31.213:9090/log
//http://192.168.31.213:9090/app/user/register?account=test?pwd=123
//nohup ./PrivateChatService &
func main() {
	b1 := time.Now().UTC()
	fmt.Println("b1:", b1)

	loge.W(b1.Format("2006-01-02 15:04:05"))
	loge.W(timeUtils.GetTimeFormat(b1.Unix(), timeUtils.DATE_TIME_FMT))
	loge.W(time.Now().Format("2006-01-02 15:04:05"))

	err := sysUtils.Install("服务", "服务", "此服务程序为后端服务功能",
		func() {
			loge.W("准备开启服务..", config.ConfigBean.ApiPort)
			mux := http.NewServeMux()
			//signal.Notify(stop_chan, os.Interrupt)
			//go socket.TcpRun(config.ConfigBean.SocketPort)
			mux.HandleFunc("/websocket", socket.WebSocketRun)
			api.Run(config.ConfigBean.ApiPort, mux)
		})
	if err != nil {
		loge.W(err)
		panic(err)
	}
}
