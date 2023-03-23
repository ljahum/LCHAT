package db

import (
	"database/sql"
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

// func GetDB() *sql.DB {
// 	return DB
// }
