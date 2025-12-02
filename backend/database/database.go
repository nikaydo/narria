package database

import (
	"database/sql"
	"log"
	"narria/backend/database/userDb"

	_ "modernc.org/sqlite"
)

type Database struct {
	Dbase *sql.DB
	User  userDb.UserDb
}

func (d *Database) CreteTables() error {
	_, err := d.Dbase.Exec(`CREATE TABLE IF NOT EXISTS users (
		userUUID TEXT PRIMARY KEY,
		username TEXT NOT NULL UNIQUE,
		security TEXT NOT NULL
	)`)
	if err != nil {
		return err
	}
	_, err = d.Dbase.Exec(`CREATE TABLE IF NOT EXISTS usersSystem (
	userUUID TEXT PRIMARY KEY,	
	language TEXT,
	theme TEXT default 'light',
	timezone TEXT,
	date_format TEXT,
	time_format TEXT,
	currency TEXT,
	measurement_system TEXT,
	plugins TEXT
	)`)
	if err != nil {
		return err
	}
	_, err = d.Dbase.Exec(`CREATE TABLE IF NOT EXISTS usersBio (
	userUUID TEXT PRIMARY KEY,	
	name TEXT,
	surname TEXT,
	birthday DATETIME,
	address TEXT,
	country TEXT,
	city TEXT,
	gender TEXT
	)`)
	if err != nil {
		return err
	}
	_, err = d.Dbase.Exec(`
	CREATE TRIGGER IF NOT EXISTS create_user_data
	AFTER INSERT ON users
	FOR EACH ROW
	BEGIN
		INSERT INTO usersSystem (userUUID, language, theme, timezone, date_format, time_format, currency, measurement_system)
		VALUES (NEW.userUUID, '', '', '', '', '', '', '');
		INSERT INTO usersBio (userUUID, name, surname, birthday, address, country, city, gender) 
		VALUES (NEW.userUUID, '', '', '', '', '', '', '');
	END;`)
	if err != nil {
		return err
	}
	return nil
}

func InitBD(path string) (*Database, error) {
	db, err := sql.Open("sqlite", path)
	if err != nil {
		return nil, err
	}
	var dbase = &Database{Dbase: db}
	err = dbase.CreteTables()
	if err != nil {
		return nil, err
	}
	log.Println("Initialized app database successfully")
	dbase.User = userDb.UserDb{Dbase: dbase.Dbase}
	return dbase, nil
}
