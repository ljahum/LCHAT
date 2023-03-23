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

type UerRegeditRouter struct {
}

func (router *UerRegeditRouter) Router(ginServer *gin.Engine) {
	ginServer.GET("/regedit", router.register)
	ginServer.POST("/regedit", router.checkRegister)

}
func (router *UerRegeditRouter) register(context *gin.Context) {
	context.HTML(http.StatusOK, "regedit.html", nil)
}

func (uc UerRegeditRouter) checkRegister(context *gin.Context) {
	var StatuedFlow common.StatusFlow
	var feedBack common.Feedback
	//发
	data := common.UserForm{context.PostForm("userID"), context.PostForm("password")}
	fmt.Println("data", data)
	jsondata, _ := json.Marshal(data)
	b64Jdata := base64.StdEncoding.EncodeToString(jsondata)

	StatuedFlow.Status = common.StatusRegedit
	StatuedFlow.Payload = b64Jdata
	fmt.Println("StatuedFlow", StatuedFlow)
	JsonStatuedFlow, _ := json.Marshal(StatuedFlow)
	wsHandler.Send(JsonStatuedFlow)
	//收
	feedBackBytes := <-wsHandler.Recvdata
	_ = json.Unmarshal(feedBackBytes, &feedBack)

	fmt.Println("feedBack.Flag", feedBack.Flag)
	fmt.Println("feedBack.MsgList", feedBack.MsgList) //以后存到全局消息列表
	//每次刷新都是直接全部刷新哈哈哈
	response := map[string]interface{}{
		"regeditCheck": feedBack.Flag,
	}
	context.JSON(http.StatusOK, response)
}
