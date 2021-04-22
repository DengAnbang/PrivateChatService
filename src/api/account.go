package api

import (
	"fmt"
	"gitee.com/DengAnbang/PrivateChatService/src/bean"
	"gitee.com/DengAnbang/PrivateChatService/src/dbops"
	"gitee.com/DengAnbang/PrivateChatService/src/socket/push"
	"gitee.com/DengAnbang/PrivateChatService/src/socket/socketConst"
	"gitee.com/DengAnbang/goutils/timeUtils"
	"math/rand"
	"strconv"

	"gitee.com/DengAnbang/goutils/httpUtils"
	"net/http"
)

/**
* showdoc
* @catalog 接口文档/用户相关
* @title 用户注册
* @description 用户注册的接口
* @method socket type=101
* @url /pc/user/register
* @param account 必选 string 登陆名
* @param pwd 必选 string 密码
* @param name 必选 string 用户昵称
* @param jobNumber 选填 string 用户工号
* @param phoneNumber 选填 string 用户电话
* @param headPortrait 选填 string 用户头像
* @param credentialType 选填 string 证件类型
* @param credentialNumber 选填 string 证件号码
* @return {"code":0,"type":0,"message":"","debug_message":"","data":UserBean}
* @remark 最后修改时间:2018.12.07 10.30
* @number 1
 */
func UserRegisterHttp(_ http.ResponseWriter, r *http.Request) error {
	account := httpUtils.GetValueFormRequest(r, "account")
	pwd := httpUtils.GetValueFormRequest(r, "pwd")
	user_name := httpUtils.GetValueFormRequest(r, "user_name")
	headPortrait := httpUtils.GetValueFormRequest(r, "headPortrait")
	userBean := bean.UserBean{
		UserName:     user_name,
		Account:      account,
		Pwd:          pwd,
		VipTime:      timeUtils.GetCurrentTimeFormat(timeUtils.DATE_TIME_FMT),
		HeadPortrait: headPortrait,
	}
	user, err := dbops.UserRegister(userBean)
	if err != nil {
		return err
	}
	return bean.NewSucceedMessage(user)
}

/**
* showdoc
* @catalog 接口文档/用户相关
* @title 用户登陆
* @description 用户登陆的接口
* @method socket type=102
* @url /pc/user/login
* @param account 必选 string 登陆名
* @param pwd 必选 string 密码
* @return {"code":0,"type":0,"message":"","debug_message":"","data":UserBean}
* @remark http请求登录成功后，需要用socket发送一个 type为100，参数为userId的消息，使用socket登录则不需要
* @number 2
 */
func UserLoginHttp(_ http.ResponseWriter, r *http.Request) error {
	account := httpUtils.GetValueFormRequest(r, "account")
	pwd := httpUtils.GetValueFormRequest(r, "pwd")
	user, err := dbops.UserLogin(account, pwd)
	if err != nil {
		return err
	}
	timestamp := timeUtils.GetTimestamp()
	vipTime, err := strconv.ParseInt(user.VipTime, 10, 64)
	if err != nil {
		return err
	}
	if timestamp >= vipTime {
		return bean.NewSucceedMessage(user).SetCode("2").SetMsg("账号已经过期了")
	}
	return bean.NewSucceedMessage(user)

}
func UserSelectByFuzzySearchHttp(_ http.ResponseWriter, r *http.Request) error {
	word := httpUtils.GetValueFormRequest(r, "word")
	user, err := dbops.UserSelectByFuzzySearch(word)
	if err != nil {
		return err
	}
	return bean.NewSucceedMessage(user)

}
func UserSelectByFuzzySearchAllHttp(_ http.ResponseWriter, r *http.Request) error {
	word := httpUtils.GetValueFormRequest(r, "word")
	user, err := dbops.UserSelectByFuzzySearchAll(word)
	if err != nil {
		return err
	}
	return bean.NewSucceedMessage(user)

}
func UserUpdateHttp(_ http.ResponseWriter, r *http.Request) error {
	account := httpUtils.GetValueFormRequest(r, "account")
	pwd := httpUtils.GetValueFormRequest(r, "pwd")
	user_name := httpUtils.GetValueFormRequest(r, "user_name")
	headPortrait := httpUtils.GetValueFormRequest(r, "headPortrait")
	vip_time := httpUtils.GetValueFormRequest(r, "vip_time")
	userBean := bean.UserBean{
		UserName:     user_name,
		UserId:       "",
		Account:      account,
		HeadPortrait: headPortrait,
		VipTime:      vip_time,
		Pwd:          pwd,
	}
	if len(vip_time) > 0 {
		_, err := strconv.ParseInt(vip_time, 10, 32)
		if err != nil {
			return bean.NewErrorMessage("vip时间格式化错误,vip时间只能是数字表示天数").SetDeBugMessage(err.Error())
		}

	}
	user, err := dbops.UserUpdate(userBean)
	if err != nil {
		return err
	}
	return bean.NewSucceedMessage(user)
}

//账号充值
func UserRechargeHttp(_ http.ResponseWriter, r *http.Request) error {
	user_id := httpUtils.GetValueFormRequest(r, "user_id")
	pay_id := httpUtils.GetValueFormRequest(r, "pay_id")
	err := dbops.UserRecharge(user_id, pay_id)
	if err != nil {
		return err
	}
	return bean.NewSucceedMessage("充值成功!")
}

//安全问题接口
func UserSecurityUpdateHttp(_ http.ResponseWriter, r *http.Request) error {
	account := httpUtils.GetValueFormRequest(r, "account")
	q1 := httpUtils.GetValueFormRequest(r, "q1")
	a1 := httpUtils.GetValueFormRequest(r, "a1")
	q2 := httpUtils.GetValueFormRequest(r, "q2")
	a2 := httpUtils.GetValueFormRequest(r, "a2")
	err := dbops.UserSecurityUpdate(account, q1, a1, q2, a2)
	if err != nil {
		return err
	}
	return bean.NewSucceedMessage("修改成功")
}
func UserSelectSecurityByAccountHttp(_ http.ResponseWriter, r *http.Request) error {
	account := httpUtils.GetValueFormRequest(r, "account")
	securityBean, err := dbops.UserSelectSecurityByAccount(account)
	if err != nil {
		return err
	}
	return bean.NewSucceedMessage(securityBean)
}
func UserSelectByAccountHttp(_ http.ResponseWriter, r *http.Request) error {
	account := httpUtils.GetValueFormRequest(r, "account")
	userBean, err := dbops.UserSelectByAccount(account)
	if err != nil {
		return err
	}
	return bean.NewSucceedMessage(userBean)
}
func UserSelectByIdHttp(_ http.ResponseWriter, r *http.Request) error {
	user_id := httpUtils.GetValueFormRequest(r, "user_id")
	my_user_id := httpUtils.GetValueFormRequest(r, "my_user_id")
	userBean, err := dbops.UserSelectById(user_id, my_user_id)
	if err != nil {
		return err
	}
	return bean.NewSucceedMessage(userBean)
}

func UserFriendAddHttp(_ http.ResponseWriter, r *http.Request) error {
	user_id := httpUtils.GetValueFormRequest(r, "user_id")
	to_user_id := httpUtils.GetValueFormRequest(r, "to_user_id")
	friend_type := httpUtils.GetValueFormRequest(r, "friend_type")
	chat_pwd := httpUtils.GetValueFormRequest(r, "chat_pwd")
	err := dbops.UserRemoveFriend(user_id, to_user_id)
	//添加好友的时候
	if friend_type == "2" {
		chat_pwd = fmt.Sprint(1000 + rand.Intn(9999-1000))
	}
	err = dbops.UserAddFriend(user_id, to_user_id, friend_type, chat_pwd)
	if err != nil {
		return err
	}
	//申请好友,通知对方
	if friend_type == "2" {
		data := bean.SocketData{
			TargetId: to_user_id,
			SenderId: user_id,
			Type:     socketConst.TYPE_FRIEND_ADD,
			Msg:      "",
			DebugMsg: "",
			Data:     "",
		}
		push.PushSocket(&data)
		return bean.NewSucceedMessage(chat_pwd)
	}
	//通过对方好友,提醒对方更新好友列表
	if friend_type == "1" {
		data := bean.SocketData{
			TargetId: to_user_id,
			SenderId: user_id,
			Type:     socketConst.TYPE_FRIEND_CHANGE,
			Msg:      "",
			DebugMsg: "",
			Data:     "",
		}
		push.PushSocket(&data)
	}
	if friend_type == "1" {
		data := bean.SocketData{
			TargetId: user_id,
			SenderId: to_user_id,
			Type:     socketConst.TYPE_FRIEND_CHANGE,
			Msg:      "",
			DebugMsg: "",
			Data:     "",
		}
		push.PushSocket(&data)
	}

	return bean.NewSucceedMessage("添加成功!")
}
func UserFriendCommentSetHttp(_ http.ResponseWriter, r *http.Request) error {
	user_id := httpUtils.GetValueFormRequest(r, "user_id")
	to_user_id := httpUtils.GetValueFormRequest(r, "to_user_id")
	nickname := httpUtils.GetValueFormRequest(r, "nickname")

	err := dbops.UserFriendCommentSet(user_id, to_user_id, nickname)
	if err != nil {
		return err
	}
	return bean.NewSucceedMessage("申请成功!")
}
func UserFriendDeleteHttp(_ http.ResponseWriter, r *http.Request) error {
	user_id := httpUtils.GetValueFormRequest(r, "user_id")
	to_user_id := httpUtils.GetValueFormRequest(r, "to_user_id")
	err := dbops.UserRemoveFriend(user_id, to_user_id)
	if err != nil {
		return err
	}
	//通过对方,提醒对方更新好友列表
	data := bean.SocketData{
		TargetId: to_user_id,
		SenderId: user_id,
		Type:     socketConst.TYPE_FRIEND_CHANGE,
		Msg:      "",
		DebugMsg: "",
		Data:     "",
	}
	push.PushSocket(&data)

	return bean.NewSucceedMessage("删除成功!")
}

func UserSelectFriendHttp(_ http.ResponseWriter, r *http.Request) error {
	user_id := httpUtils.GetValueFormRequest(r, "user_id")
	friend_type := httpUtils.GetValueFormRequest(r, "friend_type")
	beans, err := dbops.UserSelectFriend(user_id, friend_type)
	if err != nil {
		return err
	}

	return bean.NewSucceedMessage(beans)
}
