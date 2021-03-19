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
	stmtIn, err := dbConn.Prepare("INSERT INTO table_user (user_id,account,pwd,user_name,head_portrait,vip_time)VALUES(?,?,?,?,?,?)")
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
func UserSelectById(user_id string) (bean.UserBean, error) {
	var user bean.UserBean
	stmtOut, err := dbConn.Prepare("SELECT * FROM table_user WHERE user_id = ?")
	if err != nil {
		return user, err
	}
	defer stmtOut.Close()
	rows, err := stmtOut.Query(user_id)
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
func UserUpdate(userBean bean.UserBean) (bean.UserBean, error) {
	var user bean.UserBean
	if len(userBean.Account) == 0 {
		return user, bean.NewErrorMessage("账号不能为空")
		//return user, httpUtils.NewResultError(code.NormalErr, "账号不能为空")
	}
	user, err := UserSelectByAccount(userBean.Account)
	if err != nil {
		return user, err
	}
	if len(user.UserId) < 0 {
		return user, bean.NewErrorMessage("用户不存在!")
	}
	user.Modify(userBean)
	stmtIn, err := dbConn.Prepare("UPDATE table_user SET pwd=?,user_name=?,head_portrait=?,vip_time=? WHERE account=?")
	if err != nil {
		return user, err
	}
	_, err = stmtIn.Exec(user.Pwd, user.UserName, user.HeadPortrait, user.VipTime, user.Account)
	_ = stmtIn.Close()
	return UserSelectByAccount(user.Account)
}

func UserSelectSecurityByAccount(account string) (bean.SecurityBean, error) {
	var user bean.SecurityBean
	stmtOut, err := dbConn.Prepare("SELECT question1,answer1,question2,answer2 FROM table_user WHERE account = ?")
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
		user = *bean.NewSecurityBean(mapStrings)
	}
	return user, err
}

func UserSecurityUpdate(account, q1, a1, q2, a2 string) error {
	if len(q1) == 0 || len(a1) == 0 || len(q2) == 0 || len(a2) == 0 {
		return bean.NewErrorMessage("账号不能为空")
		//return user, httpUtils.NewResultError(code.NormalErr, "账号不能为空")
	}
	user, err := UserSelectByAccount(account)
	if err != nil {
		return err
	}
	if len(user.UserId) < 0 {
		return bean.NewErrorMessage("用户不存在!")
	}
	stmtIn, err := dbConn.Prepare("UPDATE table_user SET question1=?,answer1=?,question2=?,answer2=? WHERE account=?")
	if err != nil {
		return err
	}
	_, err = stmtIn.Exec(q1, a1, q2, a2, account)
	_ = stmtIn.Close()
	return err
}
func UserAddFriend(user_id, to_user_id, friend_type string) error {
	if len(user_id) == 0 || len(to_user_id) == 0 {
		return bean.NewErrorMessage("好友不能为空")
	}
	s, err := UserSelectFriendType(user_id, to_user_id)
	if s == "1" {
		return bean.NewErrorMessage("已经是好友关系了")
	} else if len(s) == 0 {
		if len(friend_type) == 0 {
			friend_type = "0"
		} else if friend_type == "1" { //1表示直接添加好友
			stmtIn, err := dbConn.Prepare("INSERT INTO table_user_friend (user_id,to_user_id,friend_type)VALUES(?,?,?)")
			if err != nil {
				return err
			}
			_, err = stmtIn.Exec(to_user_id, user_id, friend_type)
			_ = stmtIn.Close()
			if err != nil {
				return err
			}
		}
		stmtIn, err := dbConn.Prepare("INSERT INTO table_user_friend (user_id,to_user_id,friend_type)VALUES(?,?,?)")
		if err != nil {
			return err
		}
		_, err = stmtIn.Exec(user_id, to_user_id, friend_type)
		_ = stmtIn.Close()
		if err != nil {
			return err
		}
	} else {
		return bean.NewErrorMessage("对方不允许你添加为好友")
	}

	return err
}
func UserRemoveFriend(user_id, to_user_id string) error {
	if len(user_id) == 0 || len(to_user_id) == 0 {
		return bean.NewErrorMessage("好友不能为空")
	}
	stmtIn, err := dbConn.Prepare("DELETE FROM table_user_friend WHERE to_user_id=?")
	if err != nil {
		return err
	}
	_, err = stmtIn.Exec(to_user_id)
	_ = stmtIn.Close()
	if err != nil {
		return err
	}

	return err
}

func UserSelectFriend(user_id, friend_type string) ([]bean.UserBean, error) {
	var userBeans = make([]bean.UserBean, 0)
	if len(user_id) == 0 {
		return userBeans, bean.NewErrorMessage("用户id不能为空")
	}
	if len(friend_type) == 0 {
		friend_type = "1"
	}

	stmtOut, err := dbConn.Prepare("SELECT * FROM table_user_friend LEFT OUTER JOIN table_user ON table_user_friend.to_user_id=table_user.user_id WHERE table_user_friend.user_id = ? AND table_user_friend.friend_type = ?")
	if err != nil {
		return userBeans, err
	}
	defer stmtOut.Close()
	rows, err := stmtOut.Query(user_id, "1")
	if err != nil {
		return userBeans, err
	}
	defer rows.Close()
	for rows.Next() {
		mapStrings, err := dbutils.GetRowsMap(rows)
		if err != nil {
			return userBeans, err
		}
		newUserBean := *bean.NewUserBean(mapStrings)
		userBeans = append(userBeans, newUserBean)
	}

	return userBeans, err
}
func UserSelectFriendType(user_id, to_user_id string) (string, error) {

	if len(user_id) == 0 {
		return "", bean.NewErrorMessage("用户id不能为空")
	}
	stmtOut, err := dbConn.Prepare("SELECT * FROM table_user_friend  WHERE table_user_friend.user_id = ? AND table_user_friend.to_user_id = ?")
	if err != nil {
		return "", err
	}
	defer stmtOut.Close()
	rows, err := stmtOut.Query(user_id, to_user_id)
	if err != nil {
		return "", err
	}
	defer rows.Close()
	if rows.Next() {
		mapStrings, err := dbutils.GetRowsMap(rows)
		if err != nil {
			return "", err
		}
		friend_type := mapStrings["friend_type"]
		return friend_type, nil
	}

	return "", err
}