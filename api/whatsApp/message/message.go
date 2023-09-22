package message

import (
	"net/http"
	"strconv"
	"whatsApp/api/whatsApp/types"
	"whatsApp/core/memory"
	"whatsApp/lib/whatsmeow"
	"whatsApp/service/chatLogic"

	"github.com/gin-gonic/gin"
)

func Send(ctx *gin.Context) {
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
	var sendReq types.SendReq
	err = ctx.BindJSON(&sendReq)
	if err != nil {
		ctx.JSON(http.StatusExpectationFailed,
			map[string]interface{}{
				"status":  "error",
				"message": "参数错误",
			})
		return
	}
	user, ok := memory.GetUser(user_id)
	if !ok {
		ctx.JSON(http.StatusExpectationFailed,
			map[string]interface{}{
				"status":  "error",
				"message": "未登录",
			})
		return
	}
	msgId, err := whatsmeow.Send(user.WhatsClient, sendReq.ReceiverPhone, sendReq.Content)
	chatLogic.MessageStore(uint(user_id), sendReq.ReceiverPhone, sendReq.Content, msgId)
	if err != nil {
		ctx.JSON(http.StatusExpectationFailed,
			map[string]interface{}{
				"status":  "error",
				"message": "消息发送失败 err：" + err.Error(),
			})
		return
	}
	if err != nil {
		ctx.JSON(http.StatusExpectationFailed,
			map[string]interface{}{
				"status":     "ok",
				"message_id": msgId,
			})
		return
	}
}
