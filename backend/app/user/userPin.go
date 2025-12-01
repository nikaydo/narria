package user

import (
	"crypto/rand"
	"encoding/hex"
	"errors"
	"narria/backend/encrypt"
	"narria/backend/models"

	"github.com/google/uuid"
)

func (n *User) InitPinCode(pinCode string, userUuid uuid.UUID, dek []byte) error {
	if len(pinCode) > 4 || len(pinCode) < 12 {
		return errors.New("pin code length is not valid")
	}
	_, userData, err := n.DBase.User.SelectUserAuthData(models.UserData{Uuid: userUuid})
	if err != nil {
		return err
	}
	pinToken := make([]byte, 16)
	rand.Read(pinToken)
	SecurityPin, err := encrypt.MakeEncrypt(pinToken, []byte(pinCode))
	if err != nil {
		return err
	}
	SecurityPinDek, err := encrypt.MakeEncrypt(pinToken, dek)
	if err != nil {
		return err
	}
	userData.Pin = SecurityPin
	userData.MainForPin = SecurityPinDek
	if err := n.DBase.User.UpdateUserSecurity(userUuid, userData); err != nil {
		return err
	}
	return nil
}

func (n *User) CheckPinCode(pinCode string, userUuid uuid.UUID) ([]byte, error) {
	if len(pinCode) > 4 || len(pinCode) < 12 {
		return nil, errors.New("pin code length is not valid")
	}
	_, userData, err := n.DBase.User.SelectUserAuthData(models.UserData{Uuid: userUuid})
	if err != nil {
		return nil, err
	}
	key, err := hex.DecodeString(pinCode)
	if err != nil {
		return nil, err
	}
	pinToken, err := encrypt.GetDekUser(key, userData.MainForPin)
	if err != nil {
		return nil, err
	}
	dek, err := encrypt.GetDekUser(pinToken, userData.Pin)
	if err != nil {
		return nil, err
	}

	return dek, nil
}
