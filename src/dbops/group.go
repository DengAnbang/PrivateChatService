package dbops

import (
	"gitee.com/DengAnbang/PrivateChatService/src/bean"
	"gitee.com/DengAnbang/goutils/dbutils"
	"gitee.com/DengAnbang/goutils/utils"
	"strings"
)

func GroupRegister(user_id, group_name string) (bean.ChatGroup, error) {
	var chatGroup bean.ChatGroup
	if len(user_id) == 0 {
		return chatGroup, bean.NewErrorMessage("用户id不能为空")
	}
	if len(group_name) == 0 {
		return chatGroup, bean.NewErrorMessage("群名字不能为空")
	}

	group_id := utils.NewUUID()
	stmtIn, err := dbConn.Prepare("INSERT INTO table_friend_group (group_id,group_name,group_portrait,user_id,user_type,chat_pwd)VALUES(?,?,?,?,?,?)")
	if err != nil {
		return chatGroup, err
	}

	_, err = stmtIn.Exec(group_id, group_name, "", user_id, "1", "")
	_ = stmtIn.Close()
	if err != nil {
		return chatGroup, err
	}
	chatGroup = bean.ChatGroup{
		GroupId:       group_id,
		UserId:        user_id,
		UserType:      "1",
		ChatPwd:       "",
		GroupName:     group_name,
		GroupPortrait: "",
	}
	return chatGroup, err
}
func GroupAddUser(group_id, user_ids string) error {
	if len(user_ids) == 0 || len(group_id) == 0 {
		return bean.NewErrorMessage("信息不能为空")
	}
	group, err := GroupSelectMsg(group_id)
	if err != nil {
		return err
	}
	//stmtOut, err := dbConn.Prepare("SELECT * FROM table_friend_group  WHERE user_id = ? AND group_id = ?")
	//if err != nil {
	//	return err
	//}
	//rows, err := stmtOut.Query(user_ids, group_id)
	//if err != nil {
	//	return err
	//}
	//if rows.Next() {
	//	return bean.NewErrorMessage("已经在这个群里了!")
	//}
	userIds := strings.Split(user_ids, "#")
	if len(userIds) == 0 {
		//return bean.NewErrorMessage("人员ID不能为空")
		return nil
	}
	stmtIn, err := dbConn.Prepare("REPLACE INTO table_friend_group (group_id,user_id,user_type,chat_pwd,group_name,group_portrait)VALUES(?,?,?,?,?,?)")
	if err != nil {
		return err
	}
	for _, id := range userIds {
		if len(id) != 0 {
			_, err = stmtIn.Exec(group_id, id, "0", "", group.GroupName, group.GroupPortrait)
			if err != nil {
				return err
			}
		}
	}

	_ = stmtIn.Close()
	return err
}
func GroupRemoveUser(group_id, user_id string) error {
	if len(user_id) == 0 || len(group_id) == 0 {
		return bean.NewErrorMessage("信息不能为空")
	}

	stmtIn, err := dbConn.Prepare("DELETE FROM table_friend_group WHERE group_id=? AND user_id = ?")
	if err != nil {
		return err
	}
	_, err = stmtIn.Exec(group_id, user_id)
	_ = stmtIn.Close()
	if err != nil {
		return err
	}
	return err
}
func GroupRemoveUserAll(group_id string) error {
	if len(group_id) == 0 {
		return bean.NewErrorMessage("信息不能为空")
	}

	stmtIn, err := dbConn.Prepare("DELETE FROM table_friend_group WHERE group_id=? ")
	if err != nil {
		return err
	}
	_, err = stmtIn.Exec(group_id)
	_ = stmtIn.Close()
	if err != nil {
		return err
	}
	return err
}

func GroupSelectList(user_id string) ([]bean.ChatGroup, error) {
	var chatGroup = make([]bean.ChatGroup, 0)
	if len(user_id) == 0 {
		return chatGroup, bean.NewErrorMessage("用户id不能为空")
	}
	stmtOut, err := dbConn.Prepare("SELECT * FROM table_friend_group  WHERE user_id = ?")
	if err != nil {
		return chatGroup, err
	}
	defer stmtOut.Close()
	rows, err := stmtOut.Query(user_id)

	if err != nil {
		return chatGroup, err
	}
	defer rows.Close()
	for rows.Next() {
		mapStrings, err := dbutils.GetRowsMap(rows)
		if err != nil {
			return chatGroup, err
		}
		newUserBean := *bean.NewChatGroup(mapStrings)
		chatGroup = append(chatGroup, newUserBean)
	}
	return chatGroup, err
}
func GroupSelectUser(group_id string) ([]bean.ChatGroup, error) {
	var chatGroup = make([]bean.ChatGroup, 0)
	if len(group_id) == 0 {
		return chatGroup, bean.NewErrorMessage("用户id不能为空")
	}
	stmtOut, err := dbConn.Prepare("SELECT * FROM table_friend_group  WHERE group_id = ?")
	if err != nil {
		return chatGroup, err
	}
	defer stmtOut.Close()
	rows, err := stmtOut.Query(group_id)
	if err != nil {
		return chatGroup, err
	}
	defer rows.Close()
	for rows.Next() {
		mapStrings, err := dbutils.GetRowsMap(rows)
		if err != nil {
			return chatGroup, err
		}
		newUserBean := *bean.NewChatGroup(mapStrings)
		chatGroup = append(chatGroup, newUserBean)
	}
	return chatGroup, err
}
func GroupSelectMsg(group_id string) (bean.ChatGroup, error) {
	var chatGroup bean.ChatGroup
	if len(group_id) == 0 {
		return chatGroup, bean.NewErrorMessage("用户id不能为空")
	}
	stmtOut, err := dbConn.Prepare("SELECT * FROM table_friend_group  WHERE group_id = ? AND user_type = 1")
	if err != nil {
		return chatGroup, err
	}
	defer stmtOut.Close()
	rows, err := stmtOut.Query(group_id)
	if err != nil {
		return chatGroup, err
	}
	defer rows.Close()
	if rows.Next() {
		mapStrings, err := dbutils.GetRowsMap(rows)
		if err != nil {
			return chatGroup, err
		}
		chatGroup = *bean.NewChatGroup(mapStrings)
	}
	return chatGroup, err
}
func GroupSelectUserMsg(group_id string) ([]bean.UserBean, error) {
	var userBeans = make([]bean.UserBean, 0)
	if len(group_id) == 0 {
		return userBeans, bean.NewErrorMessage("用户id不能为空")
	}
	stmtOut, err := dbConn.Prepare("SELECT * FROM table_friend_group LEFT OUTER JOIN table_user ON table_friend_group.user_id=table_user.user_id WHERE table_friend_group.group_id = ? ")
	if err != nil {
		return userBeans, err
	}
	defer stmtOut.Close()
	rows, err := stmtOut.Query(group_id)
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
