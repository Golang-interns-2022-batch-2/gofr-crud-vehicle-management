package model

import "time"

type Vehicles struct {
	ID        int       `json:"id"`
	Model     string    `json:"model"`
	Color     string    `json:"color"`
	NumPlate  string    `json:"numplate"`
	CreatedAt time.Time `json:"craetedAt"`
	UpdatedAt time.Time `json:"updatedAt"`
	Name      string    `json:"name"`
	Launched  bool      `json:"launched"`
}
type RespnseData struct {
	Data interface{} `json:"vehicle"`
}
