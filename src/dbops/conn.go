package dbops

import (
	"database/sql"
	"fmt"
	"gitee.com/DengAnbang/PrivateChatService/src/config"
	"gitee.com/DengAnbang/goutils/loge"
	"github.com/go-sql-driver/mysql"
)

var (
	dbConn *sql.DB
	err    error
)

func init() {
	dataSource := fmt.Sprintf("%s@tcp(%s)", config.ConfigBean.DatabaseUserName+":"+config.ConfigBean.DatabasePassWord, config.ConfigBean.DatabaseAddress+":"+config.ConfigBean.DatabasePort)
	databaseName := config.ConfigBean.DatabaseName
	dataSourceName := fmt.Sprintf("%s/%s?charset=utf8", dataSource, databaseName)
	dbConn, err = sql.Open("mysql", dataSourceName)
	if err != nil {
		loge.W(err)
		panic(err)
	}
	err = dbConn.Ping()
	if err == nil {
		//数据库连接正常,初始化数据库
		err = initDatabase(dbConn)
		if err != nil {
			loge.W(err)
			panic(err)
		}
	}
	if v, ok := err.(*mysql.MySQLError); ok && v.Number == 1049 {
		dbConn, err = sql.Open("mysql", dataSource+"/")
		if err != nil {
			loge.W(err)
			panic(err)
		}
		err = dbConn.Ping()
		if err != nil {
			loge.W(err)
			panic(err)
		}
		_, err = dbConn.Exec(fmt.Sprintf("create database if not exists %s character set utf8;", databaseName))
		if err != nil {
			loge.W(err)
			panic(err)
		}
		dbConn, err = sql.Open("mysql", dataSourceName)
		if err != nil {
			loge.W(err)
			panic(err)
		}
	}
}
func initDatabase(db *sql.DB) error {
	if err = createTable(db); err != nil {
		return err
	}
	if err = updateTable(db); err != nil {
		return err
	}
	if err = InsertData(db); err != nil {
		return err
	}
	return nil
}
func InsertData(db *sql.DB) error {
	//加入账号
	//_, err = db.Exec("INSERT INTO table_user (id,user_id,account,pwd,name)VALUES(?,?,?,?,?)", 1, "10000", code.AdministratorAccount, code.AdministratorPwd, code.AdministratorName)
	////加入根节点
	//_, err = db.Exec("INSERT INTO table_tree_structure (id,node_id,node_name,node_parent_id,depth)VALUES(?,?,?,?,?)", 1, "0", "根节点", "-1", "0")
	////账号加入组织结构
	//_, err = db.Exec("REPLACE INTO table_tree_structure_user (node_id,user_id)VALUES(?,?)", "0", "10000")
	return nil

}

func createTable(db *sql.DB) error {
	if err = CreateUserTable(db); err != nil {
		return err
	}
	if err = CreateUserFriendTable(db); err != nil {
		return err
	}
	if err = CreateFriendListTable(db); err != nil {
		return err
	}
	return nil
}

func updateTable(db *sql.DB) error {
	//if err = UpTreeTable(db); err != nil {
	//	return err
	//}
	return nil
}
