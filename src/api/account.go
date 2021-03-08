package api

import (
	"gitee.com/DengAnbang/PrivateChatService/src/bean"
	"gitee.com/DengAnbang/PrivateChatService/src/dbops"
	"gitee.com/DengAnbang/goutils/timeUtils"

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
	name := httpUtils.GetValueFormRequest(r, "name")
	headPortrait := httpUtils.GetValueFormRequest(r, "headPortrait")
	userBean := bean.UserBean{
		UserName:     name,
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
	return bean.NewSucceedMessage(user)

}
