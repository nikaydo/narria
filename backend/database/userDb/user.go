package userDb

import (
	"database/sql"
	"encoding/json"
	"narria/backend/encrypt"
	"narria/backend/models"

	"github.com/google/uuid"
)

type UserDb struct {
	Dbase *sql.DB
}

func (u *UserDb) InsertUser(user models.UserData, security encrypt.Security) (uuid.UUID, error) {
	newUUID := uuid.New()
	securityJson, err := json.Marshal(security)
	if err != nil {
		return uuid.Nil, err
	}
	_, err = u.Dbase.Exec(`INSERT INTO users (userUUID, username, security) VALUES (?, ?, ?);`, newUUID.String(), user.Username, securityJson)
	return newUUID, err
}

func (u *UserDb) GetUserByUUID(uuid uuid.UUID) (models.UserData, error) {
	var user models.UserData
	err := u.Dbase.QueryRow(`SELECT userUUID, username FROM users WHERE userUUID = ?;`, uuid.String()).Scan(&user.Uuid, &user.Password)
	if err != nil {
		return user, err
	}
	return user, nil
}

func (u *UserDb) SelectUserAuthData(user models.UserData) (models.UserData, encrypt.Security, error) {
	var userData models.UserData
	var securityRaw []byte
	err := u.Dbase.QueryRow(`SELECT userUUID, username, security FROM users WHERE username = ?;`, user.Username).Scan(&userData.Uuid, &userData.Username, &securityRaw)
	if err != nil {
		return userData, encrypt.Security{}, err
	}
	var security encrypt.Security
	err = json.Unmarshal(securityRaw, &security)
	if err != nil {
		return userData, encrypt.Security{}, err
	}
	return userData, security, nil
}

func (u *UserDb) UpdateUserSecurity(uuid uuid.UUID, security encrypt.Security) error {
	securityJson, err := json.Marshal(security)
	if err != nil {
		return err
	}
	_, err = u.Dbase.Exec(`UPDATE users SET security = ? WHERE userUUID = ?;`, securityJson, uuid.String())
	if err != nil {
		return err
	}
	return nil
}
