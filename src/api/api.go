package api

import (
	"fmt"
	"gitee.com/DengAnbang/PrivateChatService/src/bean"
	"gitee.com/DengAnbang/PrivateChatService/src/code"
	"gitee.com/DengAnbang/goutils/fileUtil"
	"gitee.com/DengAnbang/goutils/httpUtils"
	"gitee.com/DengAnbang/goutils/loge"
	"gitee.com/DengAnbang/goutils/update"
	"net/http"
	"net/url"
	"os"
)

var (
	UploadWorkTime = "/UploadWorkTime"
	GetWorkTime    = "/GetWorkTime"
	//公共部分
	PublicDatabaseBackups = "/public/database/backups"
	PublicDatabaseRestore = "/public/database/restore"
	PublicFileUpload      = "/public/file/upload"
	PublicFileDelete      = "/public/file/delete"
	PublicFileUploadChat  = "/public/file/upload/chat"
	PublicUpdatesUpload   = "/public/app/updates/upload"
	PublicUpdatesCheck    = "/public/app/updates/check"
	//用户相关
	UserRegister = "/app/user/register"
	UserLogin    = "/app/user/login"
)

func Run(port string, mux *http.ServeMux) {
	//Apis["/"] = test

	//公共部分
	Apis[PublicFileUpload] = PublicFileUploadHttp
	Apis[PublicFileDelete] = PublicFileDeleteHttp
	Apis[PublicFileUploadChat] = PublicFileUploadChatHttp
	Apis[PublicUpdatesUpload] = PublicUpdatesUploadHttp
	Apis[PublicUpdatesCheck] = PublicUpdatesCheckHttp
	//Apis[PublicDatabaseBackups] = DatabaseBackupsHttp
	//Apis[PublicDatabaseRestore] = DatabaseRestoreHttp
	//用户相关
	Apis[UserRegister] = UserRegisterHttp
	Apis[UserLogin] = UserLoginHttp

	httpUtils.FileHandle(mux, code.RootName, code.RootPath)
	for k, v := range Apis {
		mux.HandleFunc(k, AppHandleFunc(v))
	}
	loge.SetHttp(mux)
	_ = update.UpgradeService(":"+port, mux)
}

var Apis = make(map[string]AppHandleFuncErr, 10)

type AppHandleFuncErr func(w http.ResponseWriter, r *http.Request) error
type FileHandler struct{}

func (f FileHandler) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	values, _ := url.PathUnescape(r.URL.String())
	fileName := code.RootPath + values
	if fileUtil.PathExists(fileName) {
		f, _ := os.Stat(fileName)
		if !f.IsDir() {
			http.ServeFile(w, r, fileName)
			return
		}
	}
	http.NotFound(w, r)
}

func AppHandleFunc(appHandle AppHandleFuncErr) func(http.ResponseWriter, *http.Request) {
	return func(w http.ResponseWriter, r *http.Request) {
		defer func() {
			i := recover()
			if resultData, ok := i.(*bean.ResultData); ok {
				resultData.WriterResponse(w)
			} else if err, ok := i.(error); ok {
				data := bean.NewErrorMessage("服务器内部错误")
				data.DebugMessage = fmt.Sprintf("%v", err)
				data.WriterResponse(w)
				loge.W(r.URL, err)
			}
		}()
		err := appHandle(w, r)
		loge.WDf("请求url%v,请求数据:%v", r.URL, r.Form)
		if resultData, ok := err.(*bean.ResultData); ok {
			resultData.WriterResponse(w)
		} else if err != nil {
			data := bean.NewErrorMessage("服务器内部错误")
			data.DebugMessage = fmt.Sprintf("%v", err)
			data.WriterResponse(w)
			loge.W(r.URL, err)
		}
	}
}
