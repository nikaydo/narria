package userDb

import (
	"narria/backend/models"

	"github.com/google/uuid"
)

func (u *UserDb) GetUserSystem(uuid uuid.UUID) (models.System, error) {
	var system models.System
	err := u.Dbase.QueryRow(`SELECT 
	language, 
	theme, 
	timezone, 
	date_format, 
	time_format, 
	currency, 
	measurement_system, 
	locale 
	FROM usersSystem 
	WHERE userUUID = ?;`,
		uuid.String()).Scan(
		&system.Language,
		&system.Theme,
		&system.Timezone,
		&system.DateFormat,
		&system.TimeFormat,
		&system.Currency,
		&system.MeasurementSystem,
		&system.Locale)
	if err != nil {
		return system, err
	}
	return system, nil
}

func (u *UserDb) SetUserSystem(uuid uuid.UUID, system models.System) error {
	_, err := u.Dbase.Exec(`UPDATE usersSystem SET 
		language = ?,
		theme = ?,
		timezone = ?,
		date_format = ?,
		time_format = ?,
		currency = ?,
		measurement_system = ?,
		locale = ?
		WHERE userUUID = ?;`,
		system.Language,
		system.Theme,
		system.Timezone,
		system.DateFormat,
		system.TimeFormat,
		system.Currency,
		system.MeasurementSystem,
		system.Locale,
		uuid.String())
	if err != nil {
		return err
	}
	return nil
}
