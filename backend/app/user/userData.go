package user

import (
	"narria/backend/models"

	"github.com/google/uuid"
)

// GetUser получает данные пользователя по uuid
func (n *User) GetUser(uuid uuid.UUID) (models.UserData, error) {
	userData, err := n.DBase.User.GetUserByUUID(uuid)
	if err != nil {
		return models.UserData{}, err
	}
	return userData, nil
}

// GetUserBio получает данные профиля пользователя по uuid
func (n *User) GetUserBio(uuid uuid.UUID) (models.Bio, error) {
	userBio, err := n.DBase.User.GetUserBio(uuid)
	if err != nil {
		return models.Bio{}, err
	}
	return userBio, nil
}

// SetUserBio обновляет данные профиля пользователя по uuid
func (n *User) SetUserBio(uuid uuid.UUID, bio models.Bio) error {
	return n.DBase.User.SetUserBio(uuid, bio)
}
