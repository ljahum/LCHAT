package common

import (
	"database/sql"
	"encoding/json"
	"fmt"

	_ "github.com/go-sql-driver/mysql"
)

var DB *sql.DB

type DBUser struct {
	Index    int64
	UserID   string
	Password string
}

var Main_domain string = "localhost"

// var DB *sql.DB

func InitDB() *sql.DB {
	db, _ := sql.Open("mysql", "root:123321@tcp(127.0.0.1:3306)/chatroom?charset=utf8")

	sqlStr := "SELECT * FROM `user_tab` WHERE userid=?"
	rows, _ := db.Query(sqlStr, "ljahum")

	defer rows.Close()
	for rows.Next() {
		var u DBUser
		err := rows.Scan(&u.Index, &u.UserID, &u.Password)
		if err != nil {
			fmt.Printf("scan failed, err:%v\n", err)
		}
		fmt.Printf("name:%s Password:%s \n", u.UserID, u.Password)
	}

	DB = db

	return db
}

//	func GetDB() *sql.DB {
//		return DB
//	}
func GetComment(SessionID string) ([]byte, []byte) {
	//获取广播的消息
	sqlStr := "select * from mail_table ORDER BY id DESC LIMIT 30"
	//sqlStr = "SELECT * FROM mail_table WHERE `to` = '" + SessionID + "' ORDER BY id DESC LIMIT 20;"
	rows, err := DB.Query(sqlStr)
	if err != nil {
		panic("fail to connect databse,err:")
	}
	defer rows.Close()
	var liuyanData []*Mail
	for rows.Next() {
		var com Mail
		err := rows.Scan(&com.Id, &com.Name, &com.Content, &com.Mail, &com.Time, &com.To)
		if err != nil {
			fmt.Printf("scan failed, err:%v\n", err)
		}
		liuyanData = append(liuyanData, &com)
	}
	PackedLiuyanData, _ := json.Marshal(liuyanData)
	return PackedLiuyanData, nil

}
