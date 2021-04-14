package bean

import (
	"encoding/json"
	"gitee.com/DengAnbang/PrivateChatService/src/code"
	"gitee.com/DengAnbang/goutils/loge"
	"net/http"
)

type ResultData struct {
	Code         string      `json:"code"`
	Type         string      `json:"type"`
	Msg          string      `json:"msg"`
	DebugMessage string      `json:"debug_message"`
	Data         interface{} `json:"data"`
}
type SocketData struct {
	TargetId string `json:"targetId"`
	SenderId string `json:"senderId"`
	Type     string `json:"type"`
	Msg      string `json:"msg"`
	DebugMsg string `json:"debug_msg"`
	Data     string `json:"data"`
}

func (r *ResultData) Error() string {
	return r.Msg
}
func NewSucceedMessage(data interface{}) *ResultData {
	return &ResultData{Code: code.OK, Msg: "", Data: data}
}
func NewErrorMessage(message string) *ResultData {
	return &ResultData{Code: code.NormalErr, Msg: message}
}
func (r *ResultData) SetDeBugMessage(message string) *ResultData {
	r.DebugMessage = message
	return r
}
func (r *ResultData) SetCode(code string) *ResultData {
	r.Code = code
	return r
}
func (r *ResultData) SetMsg(msg string) *ResultData {
	r.Msg = msg
	return r
}
func (r *ResultData) WriterResponse(w http.ResponseWriter) {
	bytes, err := json.Marshal(r)
	if err != nil {
		NewErrorMessage("编码错误").WriterResponse(w)
		return
	}
	w.Write(bytes)
	loge.WDf("返回数据----->%v", string(bytes))
}
func (r *ResultData) GetJson() string {
	bytes, err := json.Marshal(r)
	if err != nil {
		return NewErrorMessage("编码错误").GetJson()
	}
	return string(bytes)
}
