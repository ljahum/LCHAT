package router

import (
	"client/common"
	"github.com/gin-gonic/gin"
	"net/http"
)

type PlatfromRouter struct {
}

func (router *PlatfromRouter) Router(ginServer *gin.Engine) {
	ginServer.GET("/logout", router.logout)
	ginServer.GET("/home", router.home)

}

func (router *PlatfromRouter) home(context *gin.Context) {
	context.HTML(http.StatusOK, "home.html", nil)
}

func (router *PlatfromRouter) logout(context *gin.Context) {
	context.SetCookie("user_cookie", "nil", -1, "/", common.Client_domain, false, true)
	context.Redirect(http.StatusFound, "/login")
	//context.HTML(http.StatusOK, "login.html", nil)
}
