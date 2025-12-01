package app

import (
	"errors"
	"narria/backend/models"

	"github.com/google/uuid"
)

func (n *NarriaApi) ensureUserService() error {
	if n == nil || n.User == nil {
		return errors.New("user service is not configured")
	}
	return nil
}

// AuthUser авторизация пользователя.
func (n *NarriaApi) AuthUser(userData models.UserData) (models.UserData, error) {
	if err := n.ensureUserService(); err != nil {
		return models.UserData{}, err
	}
	userData, dek, err := n.User.AuthUser(userData)
	if err != nil {
		return models.UserData{}, err
	}
	n.Dek = &dek
	return userData, nil
}

// CreateUser создание нового пользователя.
func (n *NarriaApi) CreateUser(userData models.UserData) (string, string, error) {
	if err := n.ensureUserService(); err != nil {
		return "", "", err
	}
	userUUID, recovery, dek, err := n.User.CreateUser(userData)
	if err != nil {
		return "", "", err
	}
	n.Dek = &dek
	return userUUID.String(), recovery, nil
}

// RecoveryUser смена пароля пользователя по коду восстановления.
// принимает код восстановления, новый пароль и юзернейм пользователя.
func (n *NarriaApi) RecoveryUser(recoveryCode, passNew, username string) (string, error) {
	if err := n.ensureUserService(); err != nil {
		return "", err
	}
	userData, err := n.User.RecoveryUser(recoveryCode, passNew, username)
	if err != nil {
		return "", err
	}
	return userData.Uuid.String(), nil
}

func (n *NarriaApi) SetupPinCode(pinCode string, userUuid uuid.UUID) error {
	if err := n.ensureUserService(); err != nil {
		return err
	}
	if n.Dek == nil {
		return errors.New("dek is not configured")
	}
	if err := n.User.InitPinCode(pinCode, userUuid, *n.Dek); err != nil {
		return err
	}
	return nil
}

func (n *NarriaApi) CheckPinCode(pinCode string, userUuid uuid.UUID) error {
	if err := n.ensureUserService(); err != nil {
		return err
	}
	Dek, err := n.User.CheckPinCode(pinCode, userUuid)
	if err != nil {
		return err
	}
	n.Dek = &Dek
	return err
}
