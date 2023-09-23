package service

import (
	"whatsApp/service/chatMsgService"
	"whatsApp/service/chatService"
	"whatsApp/service/chatUserService"
	"whatsApp/service/whatsAppUserService"
)

type ServiceGroup struct {
	WhatsAppUserService whatsAppUserService.WhatsAppUserService
	ChatService         chatService.ChatService
	ChatUserService     chatUserService.ChatUserService
	ChatMsgService      chatMsgService.ChatMsgService
}

var ServiceApp = new(ServiceGroup)
