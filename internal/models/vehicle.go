package models

import (
	"time"

	"gopkg.in/guregu/null.v4"
)

type Vehicle struct {
	ID          int       `json:"id"`
	Model       string    `json:"model"`
	Color       string    `json:"color"`
	NumberPlate string    `json:"numberPlate"`
	UpdatedAt   time.Time `json:"updatedAt"`
	CreatedAt   time.Time `json:"createdAt"`
	Name        string    `json:"name"`
	Launched    null.Bool `json:"launched"`
}
