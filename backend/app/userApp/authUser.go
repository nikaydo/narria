package userApp

import (
	"encoding/hex"
	"narria/backend/encrypt"
	"narria/backend/models"

	"github.com/google/uuid"
)

// AuthUser проверяет пользователя и возвращает его данные(uuid,username), ошибку
func (n *UserApi) AuthUser(userData models.UserData) (models.UserData, []byte, error) {
	userDataFromDb, security, err := n.DBase.User.SelectUserAuthData(userData)
	if err != nil {
		return models.UserData{}, nil, err
	}
	pass, err := hex.DecodeString(userData.Password)
	if err != nil {
		return models.UserData{}, nil, err
	}
	dek, err := encrypt.GetDekUser(pass, security.Main)
	if err != nil {
		return models.UserData{}, nil, err
	}
	return userDataFromDb, dek, nil
}

// CreateUser создаёт новый пользователя и возвращает его UUID,ключ для востановления пороля, ошибку
func (n *UserApi) CreateUser(userData models.UserData) (uuid.UUID, string, []byte, error) {
	Security, dek, err := encrypt.InitUser(userData.Password)
	if err != nil {
		return uuid.Nil, "", nil, err
	}
	userData.Password, userData.Recovery = Security.Main.Wrapped, Security.Recovery.Wrapped
	recovery := Security.Recovery.Key
	Security.Main.Key, Security.Recovery.Key = "", ""
	userUUID, err := n.DBase.User.InsertUser(userData, Security)
	if err != nil {
		return uuid.Nil, "", nil, err
	}
	return userUUID, recovery, dek, nil
}
