package userapi

import (
	"crypto/rsa"
	"encoding/base64"
	"encoding/json"
	"fmt"
	"server/common"
	"server/sign"
)

// sqlStr := "SELECT * FROM `user_tab` WHERE `index`=?"

func CheckLogin(LoginForm common.UserForm) common.Feedback {
	var feedback common.Feedback
	feedback.Flag = true //查库 return true
	fmt.Println("接受登录请求的ID和密码", LoginForm.UserID, LoginForm.Password)
	db := common.DB
	var u common.DBUser
	sqlStr := "SELECT * FROM `user_tab` WHERE userid=?"
	// rows, err := db.Query(sqlStr, payload.UserD)
	err := db.QueryRow(sqlStr, LoginForm.UserID).Scan(&u.Index, &u.UserID, &u.Password)
	if err != nil {
		fmt.Println(err)
		feedback.Flag = false
	}
	if LoginForm.Password != u.Password {
		fmt.Println("密码错误")
		feedback.Flag = false
	}

	return feedback
}

func CheckRegedit(LoginForm common.UserForm) common.Feedback {
	var feedback common.Feedback
	fmt.Println("接受注册请求的ID和密码", LoginForm.UserID, LoginForm.Password)

	db := common.DB
	var u common.DBUser
	var sqlStr = "select * from user_tab where userid=?" // rows, err := db.Query(sqlStr, payload.UserD)
	err := db.QueryRow(sqlStr, LoginForm.UserID).Scan(&u.Index, &u.UserID, &u.Password)
	fmt.Println(err)
	if err == nil { //有该用户

		feedback.Flag = false
	} else {
		//goland:noinspection SqlDialectInspection
		sqlStr = "INSERT INTO `chatroom`.`user_tab` (`userid`, `password`) VALUES (?, ?)"
		_, err = db.Exec(sqlStr, LoginForm.UserID, LoginForm.Password)
		if err != nil {
			fmt.Println(err)
			feedback.Flag = false
		}
		feedback.Flag = true //查库 return true

	}
	return feedback

}
func CheckInsertMsg(publicKey *rsa.PublicKey, chatmsg common.ChatMsg) common.Feedback {
	data, _ := base64.StdEncoding.DecodeString(chatmsg.Msg)
	s, _ := base64.StdEncoding.DecodeString(chatmsg.Signature)

	if sign.RsaVerify(publicKey, s, data) {
		fmt.Println("认证成功")
	}
	var feedback common.Feedback
	var msg common.Comment
	_ = json.Unmarshal(data, &msg)
	db := common.DB
	sqlStr := "INSERT INTO liuyan(name,content,time,mail )VALUES(?,?,?,?)"
	_, err := db.Exec(sqlStr, msg.Name, msg.Content, msg.Time, msg.Mail)
	if err != nil {
		fmt.Printf("insert failed, err:%v\n", err)
		feedback.Flag = false
	}
	feedback.Flag = true
	return feedback
}
