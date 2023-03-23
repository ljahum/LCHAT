package common

import (
	"encoding/base64"
	"encoding/json"
	"github.com/gin-gonic/gin"
)

// used in broadcast
func SetUserCookie(u UserForm, c *gin.Context) {
	value, _ := json.Marshal(u)
	cookieStr := base64.StdEncoding.EncodeToString(value)

	c.SetCookie("user_cookie", cookieStr, 3600, "/", Client_domain, false, true)
}
