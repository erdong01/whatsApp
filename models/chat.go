package models

import "gorm.io/gorm"

type Chat struct {
	gorm.Model
	Name string `gorm:"column:name" json:"Name"` //type:string       comment:名称                version:2023-08-22 09:54
}

// TableName 表名:chat，聊天。
// 说明:
func (Chat) TableName() string {
	return "chat"
}
