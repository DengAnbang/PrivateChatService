package dbops

import (
	"database/sql"
	"fmt"
)

//创建用户表
func CreateUserTable(db *sql.DB) error {
	_, err := db.Exec(`CREATE TABLE IF NOT EXISTS table_user (
	    id INT(64) NOT NULL AUTO_INCREMENT COMMENT '数据id',
	    user_id VARCHAR(64) NOT NULL  COMMENT '用户id',
	    account VARCHAR(64) NULL DEFAULT NULL COMMENT '用户账号',
	    pwd VARCHAR(64) NULL DEFAULT NULL COMMENT '用户密码',
		user_name VARCHAR(64) NULL DEFAULT '未设置名称',
		head_portrait VARCHAR(128) NULL DEFAULT '' COMMENT '头像',
		vip_time TIMESTAMP  NOT NULL DEFAULT current_timestamp COMMENT 'vip到期时间',
	    created TIMESTAMP NOT NULL DEFAULT current_timestamp COMMENT '创建时间',
	    PRIMARY KEY (id)
	)AUTO_INCREMENT = 0`)
	_ = AddColumn(db, "table_user", "permissions", "0", "权限代码,0:默认,1:超级管理员")
	return err
}

//创建用户扩展表
func CreateUserExtensionTable(db *sql.DB) error {
	_, err := db.Exec(`CREATE TABLE IF NOT EXISTS table_user_extension (
	    user_id VARCHAR(64) NOT NULL  COMMENT '用户id',
	    question1 VARCHAR(64) NULL DEFAULT ''  COMMENT '问题1',
	    answer1 VARCHAR(64) NULL DEFAULT ''  COMMENT '答案1',
	    question2 VARCHAR(64) NULL DEFAULT ''  COMMENT '问题2',
	    answer2 VARCHAR(64) NULL DEFAULT ''  COMMENT '答案2',
	    recharge_type VARCHAR(64) NULL DEFAULT '0'  COMMENT '充值type,0表示默认,1表示首次充值',
	    created TIMESTAMP NOT NULL DEFAULT current_timestamp COMMENT '创建时间',
	    PRIMARY KEY (user_id)
	)AUTO_INCREMENT = 0`)
	return err
}

//创建好友表
func CreateUserFriendTable(db *sql.DB) error {
	_, err := db.Exec(`CREATE TABLE IF NOT EXISTS table_user_friend (
	    id INT(64) NOT NULL AUTO_INCREMENT COMMENT '数据id',
	    user_id VARCHAR(64) NOT NULL  COMMENT '用户id',
	    to_user_id VARCHAR(64) NOT NULL  COMMENT '好友id',
	    friend_type VARCHAR(64) NOT NULL  COMMENT '好友类型(0:默认,1:好友,2:未同意,3:拒绝)',
	    chat_pwd VARCHAR(64) NULL DEFAULT ''  COMMENT '聊天密码',
	    created TIMESTAMP NOT NULL DEFAULT current_timestamp COMMENT '创建时间',
	    PRIMARY KEY (id),UNIQUE KEY (user_id,to_user_id)
	)AUTO_INCREMENT = 0`)
	return err
}

//
func CreateUserFriendCommentTable(db *sql.DB) error {
	_, err := db.Exec(`CREATE TABLE IF NOT EXISTS table_user_friend_comment (
	    id INT(64) NOT NULL AUTO_INCREMENT COMMENT '数据id',
	    user_id VARCHAR(64) NOT NULL  COMMENT '用户id',
	    to_user_id VARCHAR(64) NOT NULL  COMMENT '好友id',
	    nickname VARCHAR(64) NOT NULL  COMMENT '昵称',
	    created TIMESTAMP NOT NULL DEFAULT current_timestamp COMMENT '创建时间',
	    PRIMARY KEY (id),UNIQUE KEY (user_id,to_user_id)
	)AUTO_INCREMENT = 0`)
	return err
}

//群的表
func CreateFriendListTable(db *sql.DB) error {
	_, err := db.Exec(`CREATE TABLE IF NOT EXISTS table_friend_group (
	    id INT(64) NOT NULL AUTO_INCREMENT COMMENT '数据id',
	    group_id VARCHAR(64) NOT NULL  COMMENT '组的id',
	    group_name VARCHAR(64) NOT NULL  COMMENT '组的名称',
	    group_portrait VARCHAR(128) NOT NULL  COMMENT '组的图片',
	    user_id VARCHAR(64) NOT NULL  COMMENT '好友id',
	    user_type VARCHAR(64) NOT NULL  COMMENT '用户类型(0:群员,1:群主)',
	    chat_pwd VARCHAR(64) NULL DEFAULT ''  COMMENT '聊天密码',
	    created TIMESTAMP NOT NULL DEFAULT current_timestamp COMMENT '创建时间',
	    PRIMARY KEY (id)
	)AUTO_INCREMENT = 0`)
	return err
}

//充值的表
func CreateRechargeTable(db *sql.DB) error {
	_, err := db.Exec(`CREATE TABLE IF NOT EXISTS table_recharge (
	    id INT(64) NOT NULL AUTO_INCREMENT COMMENT '数据id',
	    user_id VARCHAR(64) NOT NULL  COMMENT '用户id',
	    execution_user_id VARCHAR(64) NOT NULL  COMMENT '执行的人的id',
	    money VARCHAR(64) NOT NULL  COMMENT '充值的金额',
	    day VARCHAR(64) NOT NULL  COMMENT '充值天数',
	    recharge_type VARCHAR(64) NOT NULL  COMMENT '充值类型(zfb,gly,wx)',
	    created TIMESTAMP NOT NULL DEFAULT current_timestamp COMMENT '创建时间',
	    PRIMARY KEY (id)
	)AUTO_INCREMENT = 0`)
	return err
}

//充值的价格表
func CreatePriceTable(db *sql.DB) error {
	_, err := db.Exec(`CREATE TABLE IF NOT EXISTS table_price (
	    id INT(64) NOT NULL AUTO_INCREMENT COMMENT '数据id',
	    money VARCHAR(64) NOT NULL  COMMENT '充值的金额',
	    day VARCHAR(64) NOT NULL  COMMENT '充值天数',
	    giving_day VARCHAR(64) NOT NULL  COMMENT '赠送天数',
	    pay_image VARCHAR(128) NULL DEFAULT '' COMMENT '支付图片',
	    created TIMESTAMP NOT NULL DEFAULT current_timestamp COMMENT '创建时间',
	    PRIMARY KEY (id)
	)AUTO_INCREMENT = 0`)
	return err
}

func AddColumn(db *sql.DB, tableName, columnName string, defaultValue string, comment string) error {
	result, err := db.Prepare(`SELECT count(*) FROM information_schema.columns WHERE table_name = ? AND column_name = ?`)
	if err != nil {
		return err
	}
	defer result.Close()
	rows, err := result.Query(tableName, columnName)
	if err != nil {
		return err
	}
	defer rows.Close()
	if rows.Next() {
		var num int
		err := rows.Scan(&num)
		if err != nil {
			return err
		}
		if num <= 0 {
			_, err := db.Exec(fmt.Sprintf("ALTER TABLE %v ADD %v VARCHAR(64) NULL DEFAULT '%v' COMMENT '%v'", tableName, columnName, defaultValue, comment))
			if err != nil {
				return err
			}
		}
	}

	return nil
}
