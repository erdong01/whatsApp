package userService

import (
	"whatsApp/core"
	"whatsApp/models"
)

type UserService struct{}

func (u *UserService) Carete(user models.User) (userData models.User, err error) {
	err = core.New().Db.Where("phone = ?", user.Phone).First(&userData).Error
	if err == nil && userData.ID > 0 {
		return userData, err
	}
	err = core.New().Db.Create(&user).Error
	if err != nil {
		return userData, err
	}

	return user, nil
}

func (u *UserService) FindByPhone(phone string) (user models.User, err error) {
	err = core.New().Db.Where("phone = ?", phone).First(&user).Error
	return
}
