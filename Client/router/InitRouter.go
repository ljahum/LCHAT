package router

import (
	"fmt"
	"github.com/gin-gonic/gin"
)

func MiddleWare() gin.HandlerFunc {
	return func(context *gin.Context) {

		//fmt.Println("中间件开始执行了")

		cookieValue, _ := context.Cookie("user_cookie")
		fmt.Println("中间件检查cookieValue", cookieValue)

		//fmt.Println("中间件执行完毕", status)

	}
}

func InitRouter() *gin.Engine {
	ginServer := gin.Default()
	ginServer.Use(MiddleWare())
	ginServer.Static("/assets", "./assets")
	ginServer.LoadHTMLGlob("templates/*")
	//ginServer.SetHTMLTemplate()
	registRouter(ginServer)
	return ginServer
	//log.Print("http://localhost:" + port + "/home")
}

func registRouter(ginengine *gin.Engine) {
	new(UerLoginRouter).Router(ginengine)
	new(PlatfromRouter).Router(ginengine)
	new(UerRegeditRouter).Router(ginengine)
	new(PubController).Router(ginengine)

}
