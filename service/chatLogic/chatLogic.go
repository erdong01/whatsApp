package chatLogic

import (
	"fmt"
	"strings"
	"whatsApp/core"
	"whatsApp/models"
	"whatsApp/service"

	"go.mau.fi/whatsmeow/binary/proto"
)

// 发送消息存储
func SendMessageStore(userId uint, receiverPhone string, content string, msgId string) (err error) {
	OtherUser, err := service.ServiceApp.WhatsAppUserService.FindByPhone(receiverPhone)
	if err != nil {
		var user = models.WhatsAppUser{
			Phone: receiverPhone,
		}
		OtherUser, err = service.ServiceApp.WhatsAppUserService.Carete(user)
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
	} else {
		chatUser.OtherUserId++
		service.ServiceApp.ChatUserService.Update(chatUser)
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

// 接收消息存储
func ReceiverMessageStore(userId uint, sendPhone string, content string, msgId string) (err error) {
	OtherUser, err := service.ServiceApp.WhatsAppUserService.FindByPhone(sendPhone)
	if err != nil {
		var user = models.WhatsAppUser{
			Phone: sendPhone,
		}
		OtherUser, err = service.ServiceApp.WhatsAppUserService.Carete(user)
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
	} else {
		chatUser.OtherUserId++
		service.ServiceApp.ChatUserService.Update(chatUser)
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

// 消息同步
func HistorySync(userId uint, Conversations []*proto.Conversation) {
	if userId == 0 {
		return
	}
	for _, v := range Conversations {
		var otherPhone string
		parts := strings.Split(*v.Id, "@")
		// 获取@符号前的字符串
		otherPhone = parts[0]
		OtherUser, err := service.ServiceApp.WhatsAppUserService.FindByPhone(otherPhone)
		if err != nil {
			var user = models.WhatsAppUser{
				Phone: otherPhone,
			}
			OtherUser, err = service.ServiceApp.WhatsAppUserService.Carete(user)
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
		if err != nil {
			continue
		}

		if len(v.Messages) == 0 {
			continue
		}

		chatMsgData, err := service.ServiceApp.ChatMsgService.LastMsg(chatUser.ChatId)
		if err != nil {
			chatUser.MsgOrderId = 0
		}

		if err == nil {
			cIndex := (len(v.Messages) - 1) - chatUser.MsgOrderId
			fmt.Println("*v.Messages[cIndex].Message.Key.Id", *v.Messages[cIndex].Message.Key.Id)
			fmt.Println("chatMsgData.WsMsgId ", chatMsgData.WsMsgId)
			if *v.Messages[cIndex].Message.Key.Id != chatMsgData.WsMsgId {
				chatUser.MsgOrderId = 0
			}
		}
		var chatMsg []models.ChatMsg
		index := (len(v.Messages) - 1) - chatUser.MsgOrderId
		fmt.Println("index", index)
		for index > 0 {
			index--
			if *v.Messages[index].MsgOrderId == 1 {
				continue
			}
			chatUser.MsgOrderId++
			if *v.Messages[index].Message.Key.FromMe {
				chatMsg = append(chatMsg, models.ChatMsg{
					Content:    v.Messages[index].Message.Message.GetConversation(),
					SenderId:   userId,
					ReceiverId: OtherUser.ID,
					State:      1,
					ChatId:     chatUser.ChatId,
					WsMsgId:    *v.Messages[index].Message.Key.Id,
				})
			} else {
				chatMsg = append(chatMsg, models.ChatMsg{
					Content:    v.Messages[index].Message.Message.GetConversation(),
					SenderId:   OtherUser.ID,
					ReceiverId: userId,
					State:      1,
					ChatId:     chatUser.ChatId,
					WsMsgId:    *v.Messages[index].Message.Key.Id,
				})
			}
		}
		core.New().Db.Create(&chatMsg)
		service.ServiceApp.ChatUserService.Update(chatUser)
	}
}
