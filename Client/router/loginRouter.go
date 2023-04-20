package router

import (
	"client/common"
	"client/wsHandler"
	"encoding/base64"
	"encoding/json"
	"fmt"
	"github.com/gin-gonic/gin"
	"net/http"
)

type UerLoginRouter struct {
}

func (router *UerLoginRouter) Router(ginServer *gin.Engine) {
	ginServer.GET("/login", router.login)

	ginServer.POST("/login", router.checkLogin)

}
func (router *UerLoginRouter) login(context *gin.Context) {
	context.HTML(http.StatusOK, "login.html", nil)
}

func (router *UerLoginRouter) checkLogin(context *gin.Context) { //post
	var StatuedFlow common.StatusFlow

	//发
	data := common.UserForm{context.PostForm("userID"), context.PostForm("password")}
	jsondata, _ := json.Marshal(data)
	b64Jdata := base64.StdEncoding.EncodeToString(jsondata)

	StatuedFlow.Status = common.StatusLogin
	StatuedFlow.Payload = b64Jdata
	JsonStatuedFlow, _ := json.Marshal(StatuedFlow)
	wsHandler.Send(JsonStatuedFlow)
	//收
	var feedBack common.Feedback
	feedBackBytes := <-wsHandler.Recvdata
	_ = json.Unmarshal(feedBackBytes, &feedBack)

	fmt.Println("feedBack.Flag", feedBack.Flag)
	fmt.Println("feedBack.MsgList", feedBack.MsgList) //以后存到全局消息列表
	if feedBack.Flag == true {
		wsHandler.UserName = data.UserID
		PackedLiuyanData, _ := base64.StdEncoding.DecodeString(feedBack.MsgList)
		_ = json.Unmarshal(PackedLiuyanData, &common.LiuyanData)
		for k, v := range common.LiuyanData {
			fmt.Println(k, v)
		}
	}
	//每次刷新都是直接全部刷新哈哈哈
	response := map[string]interface{}{
		"loginCheck": feedBack.Flag,
	}
	context.JSON(http.StatusOK, response)
}
