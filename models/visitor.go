package models

import (
	"time"
	"whatsApp/core"

	"gorm.io/gorm"
)

type Visitor struct {
	gorm.Model
	Name         string    `json:"name" form:"name"`
	RegisterTime time.Time ` json:"register_time" form:"register_time"` // 注册时间
}

func (v *Visitor) TableName() string {
	return "visitor"
}

func (v *Visitor) Create() (err error) {
	err = core.New().Db.Create(&v).Error
	return
}

func (v *Visitor) Find() (visitor []Visitor, err error) {
	err = core.New().Db.Find(&visitor).Error
	return
}
