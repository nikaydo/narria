package userDb

import (
	"narria/backend/models"

	"github.com/google/uuid"
)

func (u *UserDb) GetUserBio(uuid uuid.UUID) (models.Bio, error) {
	var bio models.Bio
	err := u.Dbase.QueryRow(`SELECT 
	name, 
	surname, 
	birthday, 
	address, 
	country, 
	city, 
	gender 
	FROM usersBio 
	WHERE userUUID = ?;`,
		uuid).Scan(
		&bio.Name,
		&bio.Surname,
		&bio.Birthday,
		&bio.Address,
		&bio.Country,
		&bio.City,
		&bio.Gender)
	if err != nil {
		return bio, err
	}
	return bio, nil
}

func (u *UserDb) SetUserBio(uuid uuid.UUID, bio models.Bio) error {
	_, err := u.Dbase.Exec(`UPDATE usersBio SET 
		name = ?,
		surname = ?,
		birthday = ?,
		address = ?,
		country = ?,
		city = ?,
		gender = ?
		WHERE userUUID = ?;`,
		bio.Name,
		bio.Surname,
		bio.Birthday,
		bio.Address,
		bio.Country,
		bio.City,
		bio.Gender,
		uuid.String())
	if err != nil {
		return err
	}
	return nil
}
