package models

import "gorm.io/gorm"

// ChatMsg  聊天消息。
// 说明:
// 表名:chat_msg
// group: ChatMsg
// obsolete:
// appliesto:go 1.8+;
// namespace:hongmouer.his.models.ChatMsg
// assembly: hongmouer.his.models.go
// class:HongMouer.HIS.Models.ChatMsg
// version:2023-08-22 09:54
type ChatMsg struct {
	gorm.Model
	Content    string `gorm:"column:content" json:"Content"`        //type:string       comment:消息内容               version:2023-08-22 09:54
	SenderId   uint   `gorm:"column:sender_id" json:"SenderId"`     //type:*int         comment:发送者id               version:2023-08-22 09:54
	ReceiverId uint   `gorm:"column:receiver_id" json:"ReceiverId"` //type:*int         comment:接收者id               version:2023-08-22 09:54
	State      int    `gorm:"column:state" json:"State"`            //type:string       comment:状态 1 已读  2 未读    version:2023-08-22 09:54
	ChatId     uint   `gorm:"column:chat_id" json:"ChatId"`         //type:*int         comment:聊天id                 version:2023-08-22 14:21
	WsMsgId    string `gorm:"column:ws_msg_id" json:"WsMsgId"`      //type:string       comment:whatsApp消息Id         version:2023-08-22 19:55
}

// TableName 表名:chat_msg，聊天消息。
// 说明:
func (ChatMsg) TableName() string {
	return "chat_msg"
}
