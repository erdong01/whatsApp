package service

import (
	"whatsApp/service/chatMsgService"
	"whatsApp/service/chatService"
	"whatsApp/service/chatUserService"
	"whatsApp/service/userService"
)

type ServiceGroup struct {
	UserService     userService.UserService
	ChatService     chatService.ChatService
	ChatUserService chatUserService.ChatUserService
	ChatMsgService  chatMsgService.ChatMsgService
}

var ServiceApp = new(ServiceGroup)
