package app

import (
	"narria/backend/app/userApp"
	"narria/backend/database"
)

type NarriaApi struct {
	DBase *database.Database
	User  *userApp.UserApi
	Dek   *[]byte
}
