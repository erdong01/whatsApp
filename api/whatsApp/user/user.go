package user

import (
	"net/http"
	"strconv"
	"whatsApp/core/memory"
	"whatsApp/lib/whatsmeow"

	"github.com/gin-gonic/gin"
)

func Login(ctx *gin.Context) {
	var user_id int
	user_id, err := strconv.Atoi(ctx.GetHeader("Authorization"))
	if err != nil || user_id <= 0 {
		ctx.JSON(http.StatusExpectationFailed,
			map[string]interface{}{
				"status":  "error",
				"message": "用户id不能为空",
			})
		return
	}
	client, err := whatsmeow.NewConn(uint(user_id))
	if err != nil {
		ctx.JSON(http.StatusExpectationFailed,
			map[string]interface{}{
				"status":  "error",
				"message": err.Error(),
			})
		return
	}
	memory.SetUser(user_id, memory.User{WhatsClient: client})
	png, err := whatsmeow.GetQRChannel(client, ctx)
	if err == nil {
		ctx.String(http.StatusOK, string(png))

	}
}
