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
func GetComment() []byte {
	sqlStr := "select * from liuyan ORDER BY id DESC LIMIT 20"
	rows, err := DB.Query(sqlStr)
	if err != nil {
		panic("fail to connect databse,err:")
	}
	defer rows.Close()
	var liuyanData []*Comment
	for rows.Next() {
		var com Comment
		err := rows.Scan(&com.Id, &com.Name, &com.Content, &com.Mail, &com.Time)
		if err != nil {
			fmt.Printf("scan failed, err:%v\n", err)
		}
		liuyanData = append(liuyanData, &com)
	}
	PackedLiuyanData, _ := json.Marshal(liuyanData)
	//b64PackedLiuyan := base64.StdEncoding.EncodeToString(PackedLiuyanData)
	//fmt.Println(b64PackedLiuyan)
	return PackedLiuyanData
	//PackedLiuyanData, _ = base64.StdEncoding.DecodeString(b64PackedLiuyan)
	//_ = json.Unmarshal(PackedLiuyanData, &liuyanData)
	//fmt.Println(liuyanData)
	//for key, value := range liuyanData {
	//	fmt.Println(key, value)
}
