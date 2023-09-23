package whatsAppUserService

import (
	"whatsApp/core"
	"whatsApp/models"
)

type WhatsAppUserService struct{}

func (u *WhatsAppUserService) Carete(userWhatsApp models.WhatsAppUser) (userWhatsAppData models.WhatsAppUser, err error) {
	err = core.New().Db.Where("phone = ?", userWhatsApp.Phone).First(&userWhatsAppData).Error
	if err == nil && userWhatsAppData.ID > 0 {
		return userWhatsAppData, err
	}
	err = core.New().Db.Create(&userWhatsApp).Error
	if err != nil {
		return userWhatsAppData, err
	}

	return userWhatsApp, nil
}

func (u *WhatsAppUserService) FindByPhone(phone string) (userWhatsAppData models.WhatsAppUser, err error) {
	err = core.New().Db.Where("phone = ?", phone).First(&userWhatsAppData).Error
	return
}

func (u *WhatsAppUserService) FindById(userId uint) (userWhatsAppData models.WhatsAppUser, err error) {
	err = core.New().Db.Where("id = ?", userId).First(&userWhatsAppData).Error
	return
}
