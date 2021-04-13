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
	UserRegister               = "/app/user/register"
	UserLogin                  = "/app/user/login"
	UserUpdate                 = "/app/user/update"
	UserSecurityUpdate         = "/app/user/security/update"
	UserSelectSecurity         = "/app/user/select/security"
	UserSelectByAccount        = "/app/user/select/by/account"
	UserSelectById             = "/app/user/select/by/id"
	UserSelectByFuzzySearch    = "/app/user/select/by/fuzzy/search"
	UserSelectByFuzzySearchAll = "/app/user/select/by/fuzzy/search/all"
	UserFriendAdd              = "/app/user/friend/add"
	UserFriendDelete           = "/app/user/friend/delete"
	UserSelectFriend           = "/app/user/select/friend"

	GroupRegister      = "/app/group/register"
	GroupAddUser       = "/app/group/add/users"
	GroupRemoveUser    = "/app/group/remove/user"
	GroupRemoveUserAll = "/app/group/remove/user/all"
	GroupSelectList    = "/app/group/select/list"
	GroupSelectUser    = "/app/group/select/user"
	GroupSelectUserMsg = "/app/group/select/user/msg"

	RechargeAdd                     = "/app/recharge/add"
	RechargeSelectByType            = "/app/select/by/type"
	RechargeSelectByUserId          = "/app/select/by/user/id"
	RechargeSelectByExecutionUserId = "/app/select/by/execution/user/id"
	RechargeSelectByTime            = "/app/select/by/time"
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
	Apis[UserUpdate] = UserUpdateHttp
	Apis[UserSecurityUpdate] = UserSecurityUpdateHttp
	Apis[UserSelectSecurity] = UserSelectSecurityByAccountHttp
	Apis[UserSelectByAccount] = UserSelectByAccountHttp
	Apis[UserSelectById] = UserSelectByIdHttp
	Apis[UserSelectByFuzzySearch] = UserSelectByFuzzySearchHttp
	Apis[UserSelectByFuzzySearchAll] = UserSelectByFuzzySearchAllHttp
	Apis[UserFriendAdd] = UserFriendAddHttp
	Apis[UserFriendDelete] = UserFriendDeleteHttp

	Apis[UserSelectFriend] = UserSelectFriendHttp
	//群相关
	Apis[GroupRegister] = GroupRegisterHttp
	Apis[GroupAddUser] = GroupAddUserHttp
	Apis[GroupRemoveUser] = GroupRemoveUserHttp
	Apis[GroupRemoveUserAll] = GroupRemoveUserAllHttp
	Apis[GroupSelectList] = GroupSelectListHttp
	Apis[GroupSelectUser] = GroupSelectUserHttp
	Apis[GroupSelectUserMsg] = GroupSelectUserMsgHttp
	//充值相关
	Apis[RechargeAdd] = RechargeAddHttp
	Apis[RechargeSelectByType] = RechargeSelectByTypeHttp
	Apis[RechargeSelectByUserId] = RechargeSelectByUserIdHttp
	Apis[RechargeSelectByExecutionUserId] = RechargeSelectByExecutionUserIdHttp
	Apis[RechargeSelectByTime] = RechargeSelectByTimeHttp

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
