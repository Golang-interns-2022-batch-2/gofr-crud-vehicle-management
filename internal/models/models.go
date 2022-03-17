package models

import "time"

type Vehicle struct {
	ID          int       `json:"id"`
	Model       string    `json:"model"`
	Color       string    `json:"color"`
	NumberPlate string    `json:"numberPlate"`
	UpdatedAt   time.Time `json:"updatedAt"`
	CreatedAt   time.Time `json:"createdAt"`
	Name        string    `json:"name"`
	Launched    bool      `json:"launched"`
}
