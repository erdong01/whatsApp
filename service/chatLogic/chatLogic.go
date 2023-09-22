package chatLogic

import (
	"whatsApp/models"
	"whatsApp/service"
)

func SendMessageStore(userId uint, receiverPhone string, content string, msgId string) (err error) {
	OtherUser, err := service.ServiceApp.UserService.FindByPhone(receiverPhone)
	if err != nil {
		var user = models.User{
			Phone: receiverPhone,
		}
		OtherUser, err = service.ServiceApp.UserService.Carete(user)
		if err != nil {
			return
		}
	}
	chatUser, err := service.ServiceApp.ChatUserService.FindByOtherUserId(userId, OtherUser.ID)
	if err != nil {
		chatUser, err = service.ServiceApp.ChatUserService.Create(userId, chatUser.ID)
		if err != nil {
			return
		}
	}
	service.ServiceApp.ChatMsgService.Carete(models.ChatMsg{
		Content:    content,
		SenderId:   userId,
		ReceiverId: OtherUser.ID,
		State:      1,
		ChatId:     chatUser.ChatId,
		WsMsgId:    msgId,
	})
	return
}

func ReceiverMessageStore(userId uint, sendPhone string, content string, msgId string) (err error) {
	OtherUser, err := service.ServiceApp.UserService.FindByPhone(sendPhone)
	if err != nil {
		var user = models.User{
			Phone: sendPhone,
		}
		OtherUser, err = service.ServiceApp.UserService.Carete(user)
		if err != nil {
			return
		}
	}
	chatUser, err := service.ServiceApp.ChatUserService.FindByOtherUserId(userId, OtherUser.ID)
	if err != nil {
		chatUser, err = service.ServiceApp.ChatUserService.Create(userId, chatUser.ID)
		if err != nil {
			return
		}
	}
	service.ServiceApp.ChatMsgService.Carete(models.ChatMsg{
		Content:    content,
		SenderId:   userId,
		ReceiverId: OtherUser.ID,
		State:      1,
		ChatId:     chatUser.ChatId,
		WsMsgId:    msgId,
	})
	return
}
