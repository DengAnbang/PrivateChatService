package dbops

import (
	"gitee.com/DengAnbang/PrivateChatService/src/bean"
	"gitee.com/DengAnbang/goutils/dbutils"
)

func RechargeAdd(user_id, execution_user_id, money, day, recharge_type string) error {
	if len(user_id) == 0 {
		return bean.NewErrorMessage("充值的对象不能为空")
	}
	if len(execution_user_id) == 0 {
		return bean.NewErrorMessage("充值的执行人不能为空")
	}
	if len(money) == 0 {
		return bean.NewErrorMessage("充值的金额不能为空")
	}
	if len(day) == 0 {
		return bean.NewErrorMessage("充值的天数不能为空")
	}
	if len(recharge_type) == 0 {
		return bean.NewErrorMessage("充值的类型不能为空")
	}
	stmtIn, err := dbConn.Prepare("INSERT INTO table_recharge (user_id,execution_user_id,money,day,recharge_type)VALUES(?,?,?,?,?)")
	if err != nil {
		return err
	}
	_, err = stmtIn.Exec(user_id, execution_user_id, money, day, recharge_type)
	_ = stmtIn.Close()
	if err != nil {
		return err
	}
	return nil
}

func RechargeSelectByType(recharge_type string) ([]bean.RechargeBean, error) {
	beans := make([]bean.RechargeBean, 0)
	if len(recharge_type) == 0 {
		return beans, bean.NewErrorMessage("充值的类型不能为空")
	}
	stmtOut, err := dbConn.Prepare(`SELECT *,
       tu1.user_name as user_name,
       tu1.account as user_account,
       tu2.account as execution_user_account,
       tu2.user_name as execution_user_name,
       UNIX_TIMESTAMP(table_recharge.created) as created
FROM table_recharge
         LEFT OUTER JOIN table_user tu1 ON table_recharge.user_id = tu1.user_id
         left outer join table_user tu2 ON table_recharge.execution_user_id = tu2.user_id
WHERE table_recharge.recharge_type = ? order by table_recharge.created desc`)
	if err != nil {
		return beans, err
	}
	defer stmtOut.Close()
	rows, err := stmtOut.Query(recharge_type)
	if err != nil {
		return beans, err
	}
	defer rows.Close()
	for rows.Next() {
		mapStrings, err := dbutils.GetRowsMap(rows)
		if err != nil {
			return beans, err
		}
		newUserBean := *bean.NewRechargeBean(mapStrings)
		beans = append(beans, newUserBean)
	}
	return beans, nil
}
func RechargeSelectByUserId(userId string) ([]bean.RechargeBean, error) {
	beans := make([]bean.RechargeBean, 0)
	if len(userId) == 0 {
		return beans, bean.NewErrorMessage("充值的人的id不能为空")
	}
	stmtOut, err := dbConn.Prepare(`SELECT *,
       tu1.user_name as user_name,
       tu1.account as user_account,
       tu2.account as execution_user_account,
       tu2.user_name as execution_user_name,
       UNIX_TIMESTAMP(table_recharge.created) as created
FROM table_recharge
         LEFT OUTER JOIN table_user tu1 ON table_recharge.user_id = tu1.user_id
         left outer join table_user tu2 ON table_recharge.execution_user_id = tu2.user_id
WHERE table_recharge.user_id = ? order by table_recharge.created desc`)
	if err != nil {
		return beans, err
	}
	defer stmtOut.Close()
	rows, err := stmtOut.Query(userId)
	if err != nil {
		return beans, err
	}
	defer rows.Close()
	for rows.Next() {
		mapStrings, err := dbutils.GetRowsMap(rows)
		if err != nil {
			return beans, err
		}
		newUserBean := *bean.NewRechargeBean(mapStrings)
		beans = append(beans, newUserBean)
	}
	return beans, nil
}
func RechargeSelectByExecutionUserId(execution_user_id string) ([]bean.RechargeBean, error) {
	beans := make([]bean.RechargeBean, 0)
	if len(execution_user_id) == 0 {
		return beans, bean.NewErrorMessage("执行的人的类型不能为空")
	}
	stmtOut, err := dbConn.Prepare(`SELECT *,
        tu1.user_name as user_name,
       tu1.account as user_account,
       tu2.account as execution_user_account,
       tu2.user_name as execution_user_name,
       UNIX_TIMESTAMP(table_recharge.created) as created
FROM table_recharge
         LEFT OUTER JOIN table_user tu1 ON table_recharge.user_id = tu1.user_id
         left outer join table_user tu2 ON table_recharge.execution_user_id = tu2.user_id
WHERE table_recharge.execution_user_id = ? order by table_recharge.created desc`)
	if err != nil {
		return beans, err
	}
	defer stmtOut.Close()
	rows, err := stmtOut.Query(execution_user_id)
	if err != nil {
		return beans, err
	}
	defer rows.Close()
	for rows.Next() {
		mapStrings, err := dbutils.GetRowsMap(rows)
		if err != nil {
			return beans, err
		}
		newUserBean := *bean.NewRechargeBean(mapStrings)
		beans = append(beans, newUserBean)
	}
	return beans, nil
}
func RechargeSelectByTime(startTime, endTime string) ([]bean.RechargeBean, error) {
	beans := make([]bean.RechargeBean, 0)
	if len(startTime) == 0 || len(endTime) == 0 {
		return beans, bean.NewErrorMessage("时间不能为空!")
	}
	stmtOut, err := dbConn.Prepare(`SELECT *,
        tu1.user_name as user_name,
       tu1.account as user_account,
       tu2.account as execution_user_account,
       tu2.user_name as execution_user_name,
       UNIX_TIMESTAMP(table_recharge.created) as created
FROM table_recharge
         LEFT OUTER JOIN table_user tu1 ON table_recharge.user_id = tu1.user_id
         left outer join table_user tu2 ON table_recharge.execution_user_id = tu2.user_id
WHERE date_format(table_recharge.created,'%Y-%m-%d') between ? and ? order by table_recharge.created desc`)
	if err != nil {
		return beans, err
	}

	defer stmtOut.Close()
	rows, err := stmtOut.Query(startTime, endTime)
	if err != nil {
		return beans, err
	}
	defer rows.Close()
	for rows.Next() {
		mapStrings, err := dbutils.GetRowsMap(rows)
		if err != nil {
			return beans, err
		}
		newUserBean := *bean.NewRechargeBean(mapStrings)
		beans = append(beans, newUserBean)
	}
	return beans, nil
}
func RechargeSelectAll() ([]bean.RechargeBean, error) {
	beans := make([]bean.RechargeBean, 0)
	stmtOut, err := dbConn.Prepare(`SELECT *,
       tu1.user_name as user_name,
       tu1.account as user_account,
       tu2.account as execution_user_account,
       tu2.user_name as execution_user_name,
       UNIX_TIMESTAMP(table_recharge.created) as created
FROM table_recharge
         LEFT OUTER JOIN table_user tu1 ON table_recharge.user_id = tu1.user_id
         left outer join table_user tu2 ON table_recharge.execution_user_id = tu2.user_id
 order by table_recharge.created desc`)
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
		newUserBean := *bean.NewRechargeBean(mapStrings)
		beans = append(beans, newUserBean)
	}
	return beans, nil
}
