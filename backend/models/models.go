package models

import (
	"time"

	"github.com/google/uuid"
)

type UserData struct {
	Uuid     uuid.UUID `json:"uuid"`
	Username string    `json:"username"`
	Password string    `json:"password"`
	Recovery string    `json:"recovery"`
}

type System struct {
	Language          string `json:"language"`
	Theme             string `json:"theme"`
	Timezone          string `json:"timezone"`
	DateFormat        string `json:"date_format"`
	TimeFormat        string `json:"time_format"`
	Currency          string `json:"currency"`
	MeasurementSystem string `json:"measurement_system"`
	Locale            string `json:"locale"`
}

type Bio struct {
	Name     string    `json:"name"`
	Surname  string    `json:"surname"`
	Birthday time.Time `json:"birthday"`
	Address  string    `json:"address"`
	Country  string    `json:"country"`
	City     string    `json:"city"`
	Gender   string    `json:"gender"`
}
