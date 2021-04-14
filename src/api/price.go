package api

import (
	"gitee.com/DengAnbang/PrivateChatService/src/bean"
	"gitee.com/DengAnbang/PrivateChatService/src/dbops"
	"gitee.com/DengAnbang/goutils/httpUtils"
	"net/http"
)

func PriceAddHttp(_ http.ResponseWriter, r *http.Request) error {
	money := httpUtils.GetValueFormRequest(r, "money")
	day := httpUtils.GetValueFormRequest(r, "day")
	giving_day := httpUtils.GetValueFormRequest(r, "giving_day")
	err := dbops.PriceAdd(money, day, giving_day)
	if err != nil {
		return err
	}
	return bean.NewSucceedMessage("添加成功!")
}
func PriceDeleteHttp(_ http.ResponseWriter, r *http.Request) error {
	id := httpUtils.GetValueFormRequest(r, "id")
	err := dbops.PriceDelete(id)
	if err != nil {
		return err
	}
	return bean.NewSucceedMessage("删除成功!")
}
func PriceUpdateHttp(_ http.ResponseWriter, r *http.Request) error {
	money := httpUtils.GetValueFormRequest(r, "money")
	day := httpUtils.GetValueFormRequest(r, "day")
	giving_day := httpUtils.GetValueFormRequest(r, "giving_day")
	id := httpUtils.GetValueFormRequest(r, "id")
	err := dbops.PriceUpdate(money, day, giving_day, id)
	if err != nil {
		return err
	}
	return bean.NewSucceedMessage("修改成功!")
}
func PriceSelectByIdHttp(_ http.ResponseWriter, r *http.Request) error {
	id := httpUtils.GetValueFormRequest(r, "id")
	priceBean, err := dbops.PriceSelectById(id)
	if err != nil {
		return err
	}
	return bean.NewSucceedMessage(priceBean)
}
func PriceSelectAllHttp(_ http.ResponseWriter, r *http.Request) error {
	beans, err := dbops.PriceSelectAll()
	if err != nil {
		return err
	}
	return bean.NewSucceedMessage(beans)
}
