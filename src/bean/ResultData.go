package bean

import (
	"encoding/json"
	"gitee.com/DengAnbang/PrivateChatService/src/code"
	"gitee.com/DengAnbang/goutils/loge"
	"net/http"
)

type ResultData struct {
	Code         int         `json:"code"`
	Type         int         `json:"type"`
	Message      string      `json:"message"`
	DebugMessage string      `json:"debug_message"`
	Data         interface{} `json:"data"`
}
type RequestData struct {
	Code     int                    `json:"code"`
	Msg      string                 `json:"msg"`
	DebugMsg string                 `json:"debug_msg"`
	Data     map[string]interface{} `json:"data"`
}

func (r *ResultData) Error() string {
	return r.Message
}
func NewSucceedMessage(data interface{}) *ResultData {
	return &ResultData{Code: code.OK, Message: "", Data: data}
}
func NewErrorMessage(message string) *ResultData {
	return &ResultData{Code: code.NormalErr, Message: message}
}
func (r *ResultData) SetDeBugMessage(message string) *ResultData {
	r.DebugMessage = message
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
