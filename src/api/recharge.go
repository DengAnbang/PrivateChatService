package api

import (
	"gitee.com/DengAnbang/PrivateChatService/src/bean"
	"gitee.com/DengAnbang/PrivateChatService/src/dbops"
	"gitee.com/DengAnbang/goutils/httpUtils"
	"net/http"
)

func RechargeAddHttp(_ http.ResponseWriter, r *http.Request) error {
	user_id := httpUtils.GetValueFormRequest(r, "user_id")
	execution_user_id := httpUtils.GetValueFormRequest(r, "execution_user_id")
	money := httpUtils.GetValueFormRequest(r, "money")
	day := httpUtils.GetValueFormRequest(r, "day")
	recharge_type := httpUtils.GetValueFormRequest(r, "recharge_type")
	err := dbops.RechargeAdd(user_id, execution_user_id, money, day, recharge_type)
	if err != nil {
		return err
	}
	return bean.NewSucceedMessage("提交成功!")
}
func RechargeSelectByTypeHttp(_ http.ResponseWriter, r *http.Request) error {
	recharge_type := httpUtils.GetValueFormRequest(r, "recharge_type")
	beans, err := dbops.RechargeSelectByType(recharge_type)
	if err != nil {
		return err
	}
	return bean.NewSucceedMessage(beans)
}
func RechargeSelectByUserIdHttp(_ http.ResponseWriter, r *http.Request) error {
	userId := httpUtils.GetValueFormRequest(r, "user_id")
	beans, err := dbops.RechargeSelectByUserId(userId)
	if err != nil {
		return err
	}
	return bean.NewSucceedMessage(beans)
}
func RechargeSelectByExecutionUserIdHttp(_ http.ResponseWriter, r *http.Request) error {
	execution_user_id := httpUtils.GetValueFormRequest(r, "execution_user_id")
	beans, err := dbops.RechargeSelectByExecutionUserId(execution_user_id)
	if err != nil {
		return err
	}
	return bean.NewSucceedMessage(beans)
}
func RechargeSelectByTimeHttp(_ http.ResponseWriter, r *http.Request) error {
	startTime := httpUtils.GetValueFormRequest(r, "startTime")
	endTime := httpUtils.GetValueFormRequest(r, "endTime")
	beans, err := dbops.RechargeSelectByTime(startTime, endTime)
	if err != nil {
		return err
	}
	return bean.NewSucceedMessage(beans)
}
func RechargeSelectAllHttp(_ http.ResponseWriter, r *http.Request) error {
	beans, err := dbops.RechargeSelectAll()
	if err != nil {
		return err
	}
	return bean.NewSucceedMessage(beans)
}
