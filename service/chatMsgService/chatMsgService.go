package chatMsgService

import (
	"whatsApp/core"
	"whatsApp/models"
)

type ChatMsgService struct{}

func (c *ChatMsgService) Carete(chatMsg models.ChatMsg) (err error) {
	err = core.New().Db.Create(&chatMsg).Error
	return
}
