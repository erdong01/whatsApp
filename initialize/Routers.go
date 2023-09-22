package initialize

import (
	"whatsApp/api/whatsApp/user"

	"github.com/gin-gonic/gin"
)

func Routers() *gin.Engine {
	Router := gin.Default()
	userGroup := Router.Group("user")
	userGroup.GET("/whats_login", user.Login)

	msgGroup := Router.Group("msg")
	msgGroup.POST("/send", user.Login)
	return Router
}
