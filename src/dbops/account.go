package dbops

import (
	"database/sql"
	_ "database/sql"
	"fmt"
	"gitee.com/DengAnbang/PrivateChatService/src/bean"
	"gitee.com/DengAnbang/PrivateChatService/src/socket/push"
	"gitee.com/DengAnbang/goutils/dbutils"
	"gitee.com/DengAnbang/goutils/loge"
	"gitee.com/DengAnbang/goutils/timeUtils"
	"gitee.com/DengAnbang/goutils/utils"
	"strconv"
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
	stmtOut, err := dbConn.Prepare("SELECT * ,UNIX_TIMESTAMP(table_user.vip_time) as vip_time FROM table_user WHERE account = ?")
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
func UserSelectByFuzzySearch(word string) ([]bean.UserBean, error) {
	var userBeans = make([]bean.UserBean, 0)
	if len(word) == 0 {
		return userBeans, bean.NewErrorMessage("搜索信息不能为空")
	}
	userBean, err := UserSelectByAccount(word)
	if err != nil {
		return userBeans, err
	}
	if userBean.Account != "" {
		userBeans = append(userBeans, userBean)
	}

	stmtOut, err := dbConn.Prepare("SELECT *  ,UNIX_TIMESTAMP(table_user.vip_time) as vip_time FROM table_user WHERE user_name like CONCAT('%',?,'%')")
	if err != nil {
		return userBeans, err
	}
	defer stmtOut.Close()
	rows, err := stmtOut.Query(word)
	if err != nil {
		return userBeans, err
	}
	defer rows.Close()
	for rows.Next() {
		mapStrings, err := dbutils.GetRowsMap(rows)
		if err != nil {
			return userBeans, err
		}
		i := *bean.NewUserBean(mapStrings)
		i.Online = push.CheckOnline(i.UserId)
		userBeans = append(userBeans, i)

	}
	return userBeans, err
}
func UserSelectByFuzzySearchAll(word string) ([]bean.UserBean, error) {
	var userBeans = make([]bean.UserBean, 0)
	if len(word) != 0 {
		beans, err := UserSelectByFuzzySearch(word)
		if err != nil {
			return beans, err
		}
		return beans, err
	} else {
		stmtOut, err := dbConn.Prepare("SELECT *  ,UNIX_TIMESTAMP(table_user.vip_time) as vip_time FROM table_user")
		if err != nil {
			return userBeans, err
		}
		defer stmtOut.Close()
		rows, err := stmtOut.Query()
		if err != nil {
			return userBeans, err
		}
		defer rows.Close()
		for rows.Next() {
			mapStrings, err := dbutils.GetRowsMap(rows)
			if err != nil {
				return userBeans, err
			}
			i := *bean.NewUserBean(mapStrings)
			i.Online = push.CheckOnline(i.UserId)
			userBeans = append(userBeans, i)
		}
	}

	return userBeans, err
}
func UserSelectById(user_id, my_user_id string) (bean.UserBean, error) {
	var user bean.UserBean
	stmtOut, err := dbConn.Prepare(`SELECT *
     ,UNIX_TIMESTAMP(table_user.vip_time) as vip_time
     ,table_user.user_id as table_user_id
FROM table_user
    LEFT OUTER JOIN table_user_friend_comment ON table_user.user_id= table_user_friend_comment.to_user_id AND ?= table_user_friend_comment.user_id
    LEFT OUTER JOIN table_user_friend ON ((table_user.user_id=table_user_friend.user_id AND table_user_friend.to_user_id=?) OR (table_user.user_id=table_user_friend.to_user_id AND table_user_friend.user_id=?))
WHERE table_user.user_id = ?`)
	if err != nil {
		return user, err
	}
	defer stmtOut.Close()
	rows, err := stmtOut.Query(my_user_id, my_user_id, my_user_id, user_id)
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
		user.UserId = user_id
		if len(mapStrings["table_user_id"]) != 0 {
			user.UserId = mapStrings["table_user_id"]
		}
		user.Online = push.CheckOnline(user.UserId)
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
	vipTime, err := strconv.ParseInt(user.VipTime, 10, 64)
	if err != nil {
		return user, err
	}
	format := timeUtils.GetTimeFormat(vipTime, timeUtils.DATE_TIME_FMT)
	_, err = stmtIn.Exec(user.Pwd, user.UserName, user.HeadPortrait, format, user.Account)
	if err != nil {
		return user, err
	}
	_ = stmtIn.Close()
	return UserSelectByAccount(user.Account)
}
func UserRecharge(user_id, pay_id, execution_user_id, recharge_type string) error {
	if len(user_id) == 0 {
		return bean.NewErrorMessage("账号不能为空")
	}
	if len(pay_id) == 0 {
		return bean.NewErrorMessage("充值类型不能为空")
	}
	if pay_id == "1" { //如果是首次充值,验证账号是否是首次充值
		securityBean, err := UserSelectSecurityByAccount(user_id)
		if err != nil {
			return err
		}
		//如果已经首次充值了
		if securityBean.RechargeType == "1" {
			return bean.NewErrorMessage("此充值类型,仅限首次充值账号充值!")
		}
	}

	userBean, err := UserSelectById(user_id, user_id)
	if err != nil {
		return err
	}
	priceBean, err := PriceSelectById(pay_id)
	if err != nil {
		return err
	}
	day, err := strconv.ParseInt(priceBean.Day, 10, 32)
	if err != nil {
		return err
	}
	givingDay, err := strconv.ParseInt(priceBean.GivingDay, 10, 32)
	if err != nil {
		return err
	}
	addTime := givingDay + day
	uTime, _ := strconv.ParseInt(userBean.VipTime, 10, 64)
	if uTime < timeUtils.GetTimestamp() {
		uTime = timeUtils.GetTimestamp()
	}
	finallyTime := uTime + (addTime * 24 * 60 * 60)
	stmtIn, err := dbConn.Prepare("UPDATE table_user SET pwd=?,user_name=?,head_portrait=?,vip_time=? WHERE account=?")
	if err != nil {
		return err
	}

	format := timeUtils.GetTimeFormat(finallyTime, timeUtils.DATE_TIME_FMT)
	_, err = stmtIn.Exec(userBean.Pwd, userBean.UserName, userBean.HeadPortrait, format, userBean.Account)
	if err != nil {
		return err
	}
	_ = stmtIn.Close()
	if pay_id == "1" {
		err = UserFirstRechargeUpdate(user_id, "1")
	}
	err = RechargeAdd(user_id, execution_user_id, priceBean.Money, fmt.Sprint(addTime), recharge_type)

	return err
}

func UserSelectSecurityByAccount(user_id string) (bean.SecurityBean, error) {
	var user bean.SecurityBean
	stmtOut, err := dbConn.Prepare("SELECT * FROM table_user_extension WHERE user_id = ?")
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
		user = *bean.NewSecurityBean(mapStrings)
	}
	return user, err
}

func UserSecurityUpdate(user_id, q1, a1, q2, a2 string) error {
	if len(q1) == 0 || len(a1) == 0 || len(q2) == 0 || len(a2) == 0 {
		return bean.NewErrorMessage("账号不能为空")
		//return user, httpUtils.NewResultError(code.NormalErr, "账号不能为空")
	}
	user, err := UserSelectById(user_id, user_id)
	if err != nil {
		return err
	}
	if len(user.UserId) < 0 {
		return bean.NewErrorMessage("用户不存在!")
	}
	//INSERT INTO table_user (user_id,account,pwd,user_name,head_portrait,vip_time)VALUES(?,?,?,?,?,?)
	stmtIn, err := dbConn.Prepare(`REPLACE INTO table_user_extension(question1,answer1,question2,answer2,user_id)VALUES(?,?,?,?,?)`)
	//stmtIn, err := dbConn.Prepare("UPDATE table_user_extension SET question1=?,answer1=?,question2=?,answer2=? WHERE user_id=?")
	if err != nil {
		return err
	}
	_, err = stmtIn.Exec(q1, a1, q2, a2, user_id)
	_ = stmtIn.Close()
	return err
}
func UserFirstRechargeUpdate(user_id, recharge_type string) error {
	if len(user_id) == 0 {
		return bean.NewErrorMessage("用户id不能为空!")
	}
	if len(recharge_type) == 0 {
		return bean.NewErrorMessage("recharge_type不能为空!")
	}
	user, err := UserSelectById(user_id, user_id)
	if err != nil {
		return err
	}
	if len(user.UserId) < 0 {
		return bean.NewErrorMessage("用户不存在!")
	}
	stmtIn, err := dbConn.Prepare(`REPLACE INTO table_user_extension(recharge_type,user_id)VALUES(?,?)`)
	//stmtIn, err := dbConn.Prepare("UPDATE table_user_extension SET recharge_type=? WHERE user_id=?")
	if err != nil {
		return err
	}
	_, err = stmtIn.Exec(recharge_type, user_id)
	_ = stmtIn.Close()
	return err
}
func UserAddFriend(user_id, to_user_id, friend_type, chat_pwd string) error {
	if len(user_id) == 0 || len(to_user_id) == 0 {
		return bean.NewErrorMessage("好友不能为空")
	}
	if user_id == to_user_id {
		return bean.NewErrorMessage("不能添加自己")
	}
	s, err := UserSelectFriendType(user_id, to_user_id)
	if s == "1" {
		return bean.NewErrorMessage("已经是好友关系了")
	}
	if len(friend_type) == 0 {
		friend_type = "0"
	}
	stmtIn, err := dbConn.Prepare("REPLACE INTO table_user_friend (user_id,to_user_id,friend_type,chat_pwd)VALUES(?,?,?,?)")
	if err != nil {
		return err
	}
	_, err = stmtIn.Exec(user_id, to_user_id, friend_type, chat_pwd)
	_ = stmtIn.Close()
	if err != nil {
		return err
	}
	_ = stmtIn.Close()
	return err
}
func UserFriendCommentSet(user_id, to_user_id, nickname string) error {
	if len(user_id) == 0 || len(to_user_id) == 0 {
		return bean.NewErrorMessage("好友不能为空")
	}
	if user_id == to_user_id {
		return bean.NewErrorMessage("设置自己的备注")
	}

	stmtIn, err := dbConn.Prepare("REPLACE INTO table_user_friend_comment (user_id,to_user_id,nickname)VALUES(?,?,?)")
	if err != nil {
		return err
	}
	_, err = stmtIn.Exec(user_id, to_user_id, nickname)
	_ = stmtIn.Close()
	if err != nil {
		return err
	}
	_ = stmtIn.Close()
	return err
}
func UserRemoveFriend(user_id, to_user_id string) error {
	if len(user_id) == 0 || len(to_user_id) == 0 {
		return bean.NewErrorMessage("好友不能为空")
	}
	stmtIn, err := dbConn.Prepare("DELETE FROM table_user_friend WHERE to_user_id=? AND user_id=?")
	if err != nil {
		return err
	}
	_, _ = stmtIn.Exec(to_user_id, user_id)
	_ = stmtIn.Close()
	stmtIn, err = dbConn.Prepare("DELETE FROM table_user_friend WHERE to_user_id=? AND user_id=?")
	if err != nil {
		return err
	}
	_, _ = stmtIn.Exec(user_id, to_user_id)
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
	var stmtOut *sql.Stmt
	if friend_type == "1" {
		stmtOut, err = dbConn.Prepare(`SELECT *  , 
 table_user.user_id as table_user_id,
UNIX_TIMESTAMP(table_user.vip_time) as vip_time,
table_user_friend.user_id as table_user_friend_user_id,
       table_user_friend_comment.user_id as table_user_friend_comment_user_id
FROM table_user_friend
    LEFT OUTER JOIN table_user ON (if(table_user_friend.to_user_id = ?, table_user_friend.user_id,table_user_friend.to_user_id)) = table_user.user_id
    LEFT OUTER JOIN table_user_friend_comment ON (if(table_user_friend.to_user_id = ?, table_user_friend.user_id,table_user_friend.to_user_id)) = table_user_friend_comment.to_user_id AND table_user_friend_comment.user_id=?
WHERE (table_user_friend.user_id = ? OR table_user_friend.to_user_id = ?) AND friend_type = ?`)
	}
	if friend_type == "2" {
		stmtOut, err = dbConn.Prepare(`SELECT *  , 
 table_user.user_id as table_user_id,
UNIX_TIMESTAMP(table_user.vip_time) as vip_time,
table_user_friend.user_id as table_user_friend_user_id,
       table_user_friend_comment.user_id as table_user_friend_comment_user_id
FROM table_user_friend
    LEFT OUTER JOIN table_user ON (if(table_user_friend.to_user_id = ?, table_user_friend.user_id,table_user_friend.to_user_id)) = table_user.user_id
    LEFT OUTER JOIN table_user_friend_comment ON (if(table_user_friend.to_user_id = ?, table_user_friend.user_id,table_user_friend.to_user_id)) = table_user_friend_comment.to_user_id AND table_user_friend_comment.user_id=?
WHERE (table_user_friend.to_user_id = ? OR table_user_friend.to_user_id = ?) AND friend_type = ?`)
	}
	if err != nil {
		return userBeans, err
	}
	defer stmtOut.Close()
	rows, err := stmtOut.Query(user_id, user_id, user_id, user_id, user_id, friend_type)
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
		newUserBean.UserId = mapStrings["table_user_id"]
		newUserBean.Online = push.CheckOnline(newUserBean.UserId)
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
