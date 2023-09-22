package models

import "gorm.io/gorm"

// ChatUser  聊天用户。
// 说明:
// 表名:chat_user
// group: ChatUser
// obsolete:
// appliesto:go 1.8+;
// namespace:hongmouer.his.models.ChatUser
// assembly: hongmouer.his.models.go
// class:HongMouer.HIS.Models.ChatUser
// version:2023-08-22 09:56
type ChatUser struct {
	gorm.Model
	ChatId      uint `gorm:"column:chat_id" json:"ChatId"`            //type:*int         comment:聊天id              version:2023-08-22 09:56
	UserId      uint `gorm:"column:user_id" json:"UserId"`            //type:*int         comment:用户id              version:2023-08-22 09:56
	OtherUserId uint `gorm:"column:other_user_id" json:"OtherUserId"` //type:*int         comment:对方用户id          version:2023-08-22 14:48
	MsgOrderId  int  `gorm:"column:msg_order_id" json:"MsgOrderId"`   //type:*int         comment:消息当前顺序        version:2023-08-23 00:30
}

// TableName 表名:chat_user，聊天用户。
// 说明:
func (ChatUser) TableName() string {
	return "chat_user"
}
