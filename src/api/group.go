package api

import (
	"gitee.com/DengAnbang/PrivateChatService/src/bean"
	"gitee.com/DengAnbang/PrivateChatService/src/dbops"
	"gitee.com/DengAnbang/goutils/httpUtils"
	"net/http"
)

//创建组
func GroupRegisterHttp(_ http.ResponseWriter, r *http.Request) error {
	user_id := httpUtils.GetValueFormRequest(r, "user_id")
	group, err := dbops.GroupRegister(user_id)
	if err != nil {
		return err
	}
	return bean.NewSucceedMessage(group)
}

//往组里面加人
func GroupAddUserHttp(_ http.ResponseWriter, r *http.Request) error {
	user_id := httpUtils.GetValueFormRequest(r, "user_id")
	group_id := httpUtils.GetValueFormRequest(r, "group_id")
	err := dbops.GroupAddUser(group_id, user_id)
	if err != nil {
		return err
	}
	return bean.NewSucceedMessage("添加成功")
}

//将人移除组
func GroupRemoveUserHttp(_ http.ResponseWriter, r *http.Request) error {
	user_id := httpUtils.GetValueFormRequest(r, "user_id")
	group_id := httpUtils.GetValueFormRequest(r, "group_id")
	err := dbops.GroupRemoveUser(group_id, user_id)
	if err != nil {
		return err
	}
	return bean.NewSucceedMessage("删除成功")
}

//查询用户的组列表
func GroupSelectListHttp(_ http.ResponseWriter, r *http.Request) error {
	user_id := httpUtils.GetValueFormRequest(r, "user_id")
	groups, err := dbops.GroupSelectList(user_id)
	if err != nil {
		return err
	}
	return bean.NewSucceedMessage(groups)
}

//查询群里的人的id.没有信息
func GroupSelectUserHttp(_ http.ResponseWriter, r *http.Request) error {
	group_id := httpUtils.GetValueFormRequest(r, "group_id")
	groups, err := dbops.GroupSelectUser(group_id)
	if err != nil {
		return err
	}
	return bean.NewSucceedMessage(groups)
}

//查询群里的人的信息
func GroupSelectUserMsgHttp(_ http.ResponseWriter, r *http.Request) error {
	group_id := httpUtils.GetValueFormRequest(r, "group_id")
	beans, err := dbops.GroupSelectUserMsg(group_id)
	if err != nil {
		return err
	}
	return bean.NewSucceedMessage(beans)
}
