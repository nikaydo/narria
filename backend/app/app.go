package app

import (
	"narria/backend/app/userApp"
	"narria/backend/database"
	"narria/backend/plugins"
)

type NarriaApi struct {
	//Структура для работы с базой данных
	DBase *database.Database

	//Структура для работы с пользователями
	User *userApp.UserApi

	//Структура для работы с плагинами
	Plugins *plugins.Plugins

	//Ключь для расшифровки данных пользователя
	Dek *[]byte
}
