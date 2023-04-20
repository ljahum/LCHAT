package router

import (
	"client/common"
	"client/sign"
	"client/wsHandler"
	"encoding/base64"
	"encoding/json"
	"fmt"
	"github.com/gin-gonic/gin"
	"net/http"
	"time"
)

type PubController struct {
}

//功能 get 向serverclient请求sql中前30条留言加密传输会clientserver
//聊天室get到的消息不需要认证

// Router post 向server发送消息及其签名 由服务器认证和决定要不要存入数据库
// 入库只存放 userid 和 消息 签名丢弃
func (router PubController) Router(context *gin.Engine) {

	context.GET("/square", router.room)
	context.POST("/square", router.sendMsg)
	context.GET("/room", router.MailBox)

}
func (router *PubController) sendMail(c *gin.Context) {
	//var feedback string
	var msg common.Mail
	var estimate string
	var feedBack common.Feedback

	Msg := c.PostForm("message")
	if len(Msg) > 128 {
		estimate = "消息太长了"
	} else {
		msg.Mail = "nil"
		msg.Name = wsHandler.UserName
		msg.Content = Msg
		msg.To = "admin"
		msg.Time = time.Now().Format("2006-01-02 15:04:05")
		packedmsg, _ := json.Marshal(msg)
		signature := sign.RsaSign(wsHandler.PriKey, packedmsg)

		var chatmsg common.ChatMsg
		var statuedPlain common.StatusFlow

		//
		chatmsg.Msg = base64.StdEncoding.EncodeToString(packedmsg)
		chatmsg.Signature = base64.StdEncoding.EncodeToString(signature)
		//
		var payload []byte
		payload, _ = json.Marshal(chatmsg)
		b64payload := base64.StdEncoding.EncodeToString(payload)
		statuedPlain.Payload = b64payload
		//statuedPlain.Status = common.StatusBroadcast
		statuedPlain.Status = common.StatusChat
		PackedStatuedplain, _ := json.Marshal(statuedPlain)
		//

		wsHandler.Send(PackedStatuedplain)

		//收

		feedBackBytes := <-wsHandler.Recvdata
		_ = json.Unmarshal(feedBackBytes, &feedBack)

		fmt.Println("feedBack.Flag", feedBack.Flag)
		fmt.Println("feedBack.MsgList", feedBack.MsgList) //以后存到全局消息列表
		if feedBack.Flag == true {
			estimate = "发表成功"
			PackedLiuyanData, _ := base64.StdEncoding.DecodeString(feedBack.MsgList)
			_ = json.Unmarshal(PackedLiuyanData, &common.LiuyanData)
			for k, v := range common.LiuyanData {
				fmt.Println(k, v)
			}
		} else {
			fmt.Println("insert failed")
			estimate = "数据库插入出错"
		}
	}
	///feedback改改
	//check := true
	SendJsonBack(estimate, feedBack.Flag, c)
}

func (router *PubController) room(context *gin.Context) {
	context.HTML(http.StatusOK, "arr.html", gin.H{
		"title": "Gin",
		// "stuArr": [1]*model.Comment{com1},
		"stuArr":   common.LiuyanData,
		"username": wsHandler.UserName,
	})
}
func (router PubController) sendMsg(c *gin.Context) {
	//var feedback string
	var msg common.Comment
	var estimate string
	var feedBack common.Feedback

	Msg := c.PostForm("message")
	Touser := c.PostForm("touser")
	if len(Msg) > 128 {
		estimate = "消息太长了"
	} else {
		if Touser == "" {
			msg.To = "admin"
		} else {
			msg.To = Touser
		}
		msg.Mail = "nil"
		msg.Name = wsHandler.UserName
		msg.Content = Msg
		msg.Time = time.Now().Format("2006-01-02 15:04:05")
		packedmsg, _ := json.Marshal(msg)
		signature := sign.RsaSign(wsHandler.PriKey, packedmsg)

		var chatmsg common.ChatMsg
		var statuedPlain common.StatusFlow

		//
		chatmsg.Msg = base64.StdEncoding.EncodeToString(packedmsg)
		chatmsg.Signature = base64.StdEncoding.EncodeToString(signature)
		//
		var payload []byte
		payload, _ = json.Marshal(chatmsg)
		b64payload := base64.StdEncoding.EncodeToString(payload)
		statuedPlain.Payload = b64payload
		statuedPlain.Status = common.StatusBroadcast
		PackedStatuedplain, _ := json.Marshal(statuedPlain)
		//

		wsHandler.Send(PackedStatuedplain)

		//收

		feedBackBytes := <-wsHandler.Recvdata
		_ = json.Unmarshal(feedBackBytes, &feedBack)

		fmt.Println("feedBack.Flag", feedBack.Flag)
		fmt.Println("feedBack.MsgList", feedBack.MsgList) //以后存到全局消息列表
		if feedBack.Flag == true {
			estimate = "发表成功"
			PackedLiuyanData, _ := base64.StdEncoding.DecodeString(feedBack.MsgList)
			_ = json.Unmarshal(PackedLiuyanData, &common.LiuyanData)
			for k, v := range common.LiuyanData {
				fmt.Println(k, v)
			}
		} else {
			fmt.Println("insert failed")
			estimate = "数据库插入出错"
		}
	}
	///feedback改改
	//check := true
	SendJsonBack(estimate, feedBack.Flag, c)
}

func (router *PubController) MailBox(context *gin.Context) {
	fmt.Println(wsHandler.UserName)
	context.HTML(http.StatusOK, "room.html", gin.H{
		"title": "Gin",
		// "stuArr": [1]*model.Comment{com1},
		"stuArr":   common.LiuyanData,
		"username": wsHandler.UserName,
	})
}

func SendJsonBack(feedback string, check bool, c *gin.Context) {
	messageMap := map[string]interface{}{
		"msg":   feedback,
		"check": check,
	}
	c.JSON(http.StatusOK, messageMap)
}
