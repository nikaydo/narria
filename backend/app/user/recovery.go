package user

import (
	"encoding/hex"
	"narria/backend/encrypt"
	"narria/backend/models"
)

func (n *User) RecoveryUser(recoveryCode, passNew, username string) (models.UserData, error) {
	// проверяем пользователя и возвращаем его данные(uuid,username), ошибку
	userData, security, err := n.DBase.User.SelectUserAuthData(models.UserData{Username: username})
	if err != nil {
		return models.UserData{}, err
	}
	recCode, err := hex.DecodeString(recoveryCode)
	if err != nil {
		return models.UserData{}, err
	}
	// проверяем код восстановления и получаем DEK
	dek, err := encrypt.GetDekUser(recCode, security.Recovery)
	if err != nil {
		return models.UserData{}, err
	}
	// создаем новые salt wrapped и nonce для нового пароля используя DEK
	securityNew, err := encrypt.MakeEncrypt(dek, []byte(passNew))
	if err != nil {
		return models.UserData{}, err
	}
	// в структуре security заменяем старые данные пороля на новые
	security.Main = securityNew
	// обновляем структуру security в базе данных
	if err := n.DBase.User.UpdateUserSecurity(userData.Uuid, security); err != nil {
		return models.UserData{}, err
	}
	return userData, nil
}
