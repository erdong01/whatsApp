package models

import "gorm.io/gorm"

type User struct {
	gorm.Model
	Name  string `gorm:"column:name" json:"Name"`   //type:string       comment:名称                version:2023-08-22 09:52
	Phone string `gorm:"column:phone" json:"Phone"` //type:string       comment:手机号              version:2023-08-22 09:52
}

func (User) TableName() string {
	return "whats_app_user"
}
