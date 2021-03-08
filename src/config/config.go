package config

import (
	"bytes"
	"encoding/json"
	"gitee.com/DengAnbang/PrivateChatService/src/code"

	"gitee.com/DengAnbang/goutils/fileUtil"
	"gitee.com/DengAnbang/goutils/loge"
	"io/ioutil"
	"os"
)

type Bean struct {
	ApiPort          string `json:"api_port"`
	SocketPort       string `json:"socket_port"`
	DatabaseUserName string `json:"database_user_name"`
	DatabasePassWord string `json:"database_pass_word"`
	DatabaseAddress  string `json:"database_address"`
	DatabasePort     string `json:"database_port"`
	DatabaseName     string `json:"database_name"`
	DebugLog         bool   `json:"debug_log"`
}

var ConfigBean = Bean{
	ApiPort:          "9090",
	SocketPort:       "9091",
	DatabaseUserName: "root",
	DatabasePassWord: "root123456",
	DatabaseAddress:  "47.108.172.20",
	DatabasePort:     "13306",
	DatabaseName:     "PrivateChat",
	DebugLog:         false}

func init() {
	_ = os.MkdirAll(code.ViewRootPath, os.ModePerm)
	_ = os.MkdirAll(code.FileRootPath, os.ModePerm)
	_ = os.MkdirAll(code.DatabaseRootPath, os.ModePerm)
	_ = os.MkdirAll(code.ConfigRootPath, os.ModePerm)
	loge.SetLogPath(code.LogRootPath)
	fileName := code.ConfigRootPath + "config.cfg"
	if !fileUtil.PathExists(fileName) {
		configBeanBytes, err := json.Marshal(ConfigBean)
		if err != nil {
			loge.W(err)
			panic(err)
		}
		var str bytes.Buffer
		err = json.Indent(&str, configBeanBytes, "", "  ")
		if err != nil {
			loge.W(err)
			panic(err)
		}
		_ = ioutil.WriteFile(fileName, str.Bytes(), os.ModePerm)
	} else {
		fileBytes, err := ioutil.ReadFile(fileName)
		if err != nil {
			loge.W(err)
			panic(err)
		}
		err = json.Unmarshal(fileBytes, &ConfigBean)
		if err != nil {
			loge.W(err)
			panic(err)
		}
		var str bytes.Buffer
		b, _ := json.Marshal(ConfigBean)
		err = json.Indent(&str, b, "", "  ")
		_ = ioutil.WriteFile(fileName, str.Bytes(), os.ModePerm)
	}
	loge.IsDebug = ConfigBean.DebugLog
}
