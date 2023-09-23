package chatUserService

import (
	"whatsApp/core"
	"whatsApp/models"
)

type ChatUserService struct {
}

func (c *ChatUserService) Create(userId uint, otherUserId uint) (chatUserData models.ChatUser, err error) {
	chat := models.Chat{
		Name: "",
	}
	core.New().Db.Create(&chat)
	var chatUser = models.ChatUser{
		WhatsUserId: userId,
		OtherUserId: otherUserId,
		ChatId:      chat.ID,
		MsgOrderId:  1,
	}
	err = core.New().Db.Create(&chatUser).Error
	if err != nil {
		return chatUserData, err
	}
	chatUserData = chatUser
	var otherUser = models.ChatUser{
		WhatsUserId: otherUserId,
		OtherUserId: userId,
		ChatId:      chat.ID,
	}
	err = core.New().Db.Create(&otherUser).Error
	return
}

func (c *ChatUserService) FindByUserId(userId uint) (chatUser models.ChatUser, err error) {
	err = core.New().Db.Where("user_id = ?", userId).First(&chatUser).Error
	return
}

func (c *ChatUserService) FindByOtherUserId(userId uint, otherUserId uint) (chatUser models.ChatUser, err error) {
	err = core.New().Db.Where("user_id = ?", userId).Where("other_user_id = ?", otherUserId).
		First(&chatUser).Error
	return
}
func (c *ChatUserService) Update(hatUserData models.ChatUser) {
	core.New().Db.Save(&hatUserData)
}
