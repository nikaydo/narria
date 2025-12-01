package app

import (
	"narria/backend/app/user"
	"narria/backend/database"
)

type NarriaApi struct {
	DBase *database.Database
	User  *user.User
	Dek   *[]byte
}
