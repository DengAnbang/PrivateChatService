package main

import (
	"gitee.com/DengAnbang/PrivateChatService/src/api"
	"gitee.com/DengAnbang/PrivateChatService/src/config"
	"gitee.com/DengAnbang/PrivateChatService/src/socket"
	"gitee.com/DengAnbang/goutils/loge"
	"gitee.com/DengAnbang/goutils/sysUtils"
	"net/http"
)

//http://47.108.172.20:9090/upload
//http://192.168.31.213:9090/log
//http://192.168.31.213:9090/app/user/register?account=test?pwd=123
//nohup ./PrivateChatService &
func main() {
	err := sysUtils.Install("服务", "服务", "此服务程序为后端服务功能",
		func() {
			loge.W("准备开启服务..", config.ConfigBean.ApiPort)
			mux := http.NewServeMux()
			//signal.Notify(stop_chan, os.Interrupt)
			go socket.TcpRun(config.ConfigBean.SocketPort)
			//mux.HandleFunc("/websocket", socket.WebSocketRun)
			api.Run(config.ConfigBean.ApiPort, mux)
		})
	if err != nil {
		loge.W(err)
		panic(err)
	}
}
