package dbops

import (
	_ "database/sql"
	"gitee.com/DengAnbang/PrivateChatService/src/bean"
	"gitee.com/DengAnbang/goutils/dbutils"
	"gitee.com/DengAnbang/goutils/loge"
	"gitee.com/DengAnbang/goutils/utils"
)

func UserRegister(userBean bean.UserBean) (bean.UserBean, error) {
	var user bean.UserBean
	if len(userBean.Account) == 0 {
		return user, bean.NewErrorMessage("账号不能为空")
		//return user, httpUtils.NewResultError(code.NormalErr, "账号不能为空")
	}
	if len(userBean.Pwd) == 0 {
		return user, bean.NewErrorMessage("密码不能为空")
	}
	user, err = UserSelectByAccount(userBean.Account)
	if len(user.UserId) > 0 {
		return user, bean.NewErrorMessage("账号已经存在了!")
	}
	userBean.UserId = utils.NewUUID()
	stmtIn, err := dbConn.Prepare("INSERT INTO table_user (user_id,account,pwd,name,head_portrait,vip_time)VALUES(?,?,?,?,?,?)")
	if err != nil {
		return user, err
	}
	_, err = stmtIn.Exec(userBean.UserId, userBean.Account, userBean.Pwd, userBean.UserName, userBean.HeadPortrait, userBean.VipTime)
	_ = stmtIn.Close()
	if err != nil {
		return user, err
	}
	loge.WD(bean.NewSucceedMessage(userBean).GetJson())
	return UserSelectByAccount(userBean.Account)
}

func UserLogin(account, pwd string) (user bean.UserBean, err error) {
	if len(account) == 0 {
		return user, bean.NewErrorMessage("登陆名字不能为空")
	}
	if len(pwd) == 0 {
		return user, bean.NewErrorMessage("密码不能为空")
	}
	user, err = UserSelectByAccount(account)
	if err != nil && len(user.UserId) <= 0 {
		return user, bean.NewErrorMessage("用户不存在或密码错误!")
	}
	if user.Pwd != pwd {
		return user, bean.NewErrorMessage("用户不存在或密码错误!")
	}
	return
}
func UserSelectByAccount(account string) (bean.UserBean, error) {
	var user bean.UserBean
	stmtOut, err := dbConn.Prepare("SELECT * FROM table_user WHERE account = ?")
	if err != nil {
		return user, err
	}
	defer stmtOut.Close()
	rows, err := stmtOut.Query(account)
	if err != nil {
		return user, err
	}
	defer rows.Close()
	if rows.Next() {
		mapStrings, err := dbutils.GetRowsMap(rows)
		if err != nil {
			return user, err
		}
		user = *bean.NewUserBean(mapStrings)
	}
	return user, err
}
