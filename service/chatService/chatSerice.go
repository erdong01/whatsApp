package chatService

import (
	"whatsApp/core"
	"whatsApp/models"
)

type ChatService struct {
}

func (c *ChatService) Carete(chat models.Chat) (chatData models.Chat, err error) {
	err = core.New().Db.Create(&chat).Error
	return chat, err
}
