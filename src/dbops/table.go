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
		head_portrait VARCHAR(64) NULL DEFAULT '' COMMENT '头像',
		vip_time TIMESTAMP  NOT NULL DEFAULT current_timestamp COMMENT 'vip到期时间',
	    created TIMESTAMP NOT NULL DEFAULT current_timestamp COMMENT '创建时间',
	    PRIMARY KEY (id)
	)AUTO_INCREMENT = 0`)
	_ = AddColumn(db, "table_user", " question1", "", "问题1")
	_ = AddColumn(db, "table_user", " answer1", "", "答案1")
	_ = AddColumn(db, "table_user", " question2", "", "问题2")
	_ = AddColumn(db, "table_user", " answer2", "", "答案2")
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
