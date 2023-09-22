package chatLogic

import (
	"strings"
	"whatsApp/models"
	"whatsApp/service"

	"go.mau.fi/whatsmeow/binary/proto"
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
		chatUser, err = service.ServiceApp.ChatUserService.Create(userId, OtherUser.ID)
		if err != nil {
			return
		}
	}
	service.ServiceApp.ChatMsgService.Carete(models.ChatMsg{
		Content:    content,
		SenderId:   OtherUser.ID,
		ReceiverId: userId,
		State:      1,
		ChatId:     chatUser.ChatId,
		WsMsgId:    msgId,
	})
	return
}

func HistorySync(userId uint, Conversations []*proto.Conversation) {

	for _, v := range Conversations {
		var otherPhone string
		parts := strings.Split(*v.Id, "@")
		// 获取@符号前的字符串
		otherPhone = parts[0]
		OtherUser, err := service.ServiceApp.UserService.FindByPhone(otherPhone)
		if err != nil {
			var user = models.User{
				Phone: otherPhone,
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
		if err != nil {
			continue
		}
		if len(v.Messages) == 0 {
			continue
		}
		if chatUser.MsgOrderId+1 > len(v.Messages) {
			chatUser.MsgOrderId = 0
		}
		var chatMsg []models.ChatMsg
		for i := chatUser.MsgOrderId; i < len(v.Messages); i++ {
			if i == 0 {
				continue
			}
			if *v.Messages[i].Message.Key.FromMe {
				chatMsg = append(chatMsg, models.ChatMsg{
					Content:    v.Messages[i].Message.Message.GetConversation(),
					SenderId:   userId,
					ReceiverId: OtherUser.ID,
					State:      1,
					ChatId:     chatUser.ChatId,
					WsMsgId:    *v.Messages[i].Message.Key.Id,
				})
			} else {
				chatMsg = append(chatMsg, models.ChatMsg{
					Content:    v.Messages[i].Message.Message.GetConversation(),
					SenderId:   OtherUser.ID,
					ReceiverId: userId,
					State:      1,
					ChatId:     chatUser.ChatId,
					WsMsgId:    *v.Messages[i].Message.Key.Id,
				})
			}
		}
	}
}
