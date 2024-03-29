package dbops

import (
	"gitee.com/DengAnbang/PrivateChatService/src/bean"
	"gitee.com/DengAnbang/goutils/dbutils"
)

func PriceAdd(money, day, giving_day, pay_image string) error {
	if len(money) == 0 {
		return bean.NewErrorMessage("充值的金额不能为空")
	}
	if len(day) == 0 {
		return bean.NewErrorMessage("充值的天数不能为空")
	}
	if len(giving_day) == 0 {
		giving_day = "0"
	}
	stmtIn, err := dbConn.Prepare("INSERT INTO table_price (money,day,giving_day,pay_image)VALUES(?,?,?,?)")
	if err != nil {
		return err
	}
	_, err = stmtIn.Exec(money, day, giving_day, pay_image)
	_ = stmtIn.Close()
	if err != nil {
		return err
	}
	return nil
}
func PriceDelete(id string) error {
	if len(id) == 0 {
		return bean.NewErrorMessage("id不能为空")
	}
	if id == "1" {
		return bean.NewErrorMessage("预设的充值不能删除,只能修改")
	}
	price, err := PriceSelectById(id)
	if err != nil {
		return err
	}
	if len(price.Money) < 0 {
		return bean.NewErrorMessage("id对应的价格不存在!")
	}
	stmtIn, err := dbConn.Prepare("DELETE FROM table_price  WHERE id=?")
	if err != nil {
		return err
	}
	_, err = stmtIn.Exec(id)
	if err != nil {
		return err
	}
	_ = stmtIn.Close()
	return nil
}
func PriceUpdate(money, day, giving_day, id, pay_image string) error {
	if len(id) == 0 {
		return bean.NewErrorMessage("id不能为空")
	}
	if id == "1" {
		giving_day = "0"
	}
	price, err := PriceSelectById(id)
	if err != nil {
		return err
	}
	if len(price.Money) < 0 {
		return bean.NewErrorMessage("id对应的价格不存在!")
	}
	stmtIn, err := dbConn.Prepare("UPDATE table_price SET money=?,day=?,giving_day=?,pay_image=? WHERE id=?")
	if err != nil {
		return err
	}
	_, err = stmtIn.Exec(money, day, giving_day, pay_image, id)
	if err != nil {
		return err
	}
	_ = stmtIn.Close()
	return nil
}
func PriceSelectById(id string) (bean.PriceBean, error) {
	var priceBean bean.PriceBean
	stmtOut, err := dbConn.Prepare("SELECT *  FROM table_price WHERE id = ?")
	if err != nil {
		return priceBean, err
	}
	defer stmtOut.Close()
	rows, err := stmtOut.Query(id)
	if err != nil {
		return priceBean, err
	}
	defer rows.Close()
	if rows.Next() {
		mapStrings, err := dbutils.GetRowsMap(rows)
		if err != nil {
			return priceBean, err
		}
		priceBean = *bean.NewPriceBean(mapStrings)
	}
	return priceBean, err
}
func PriceSelectAll(user_id string) ([]bean.PriceBean, error) {
	beans := make([]bean.PriceBean, 0)
	securityBean, err := UserSelectSecurityByAccount(user_id)
	if err != nil {
		return beans, err
	}
	firstRecharge := securityBean.RechargeType != "1"
	userBean, err := UserSelectById(user_id, user_id)
	if err != nil {
		return beans, err
	}

	stmtOut, err := dbConn.Prepare(`SELECT * FROM table_price`)
	if err != nil {
		return beans, err
	}
	defer stmtOut.Close()
	rows, err := stmtOut.Query()
	if err != nil {
		return beans, err
	}
	defer rows.Close()
	for rows.Next() {
		mapStrings, err := dbutils.GetRowsMap(rows)
		if err != nil {
			return beans, err
		}
		newUserBean := *bean.NewPriceBean(mapStrings)
		if newUserBean.Id == "1" {
			if userBean.Permissions == "1" || firstRecharge {
				beans = append(beans, newUserBean)
			}
		} else {
			beans = append(beans, newUserBean)
		}

	}
	return beans, nil
}
