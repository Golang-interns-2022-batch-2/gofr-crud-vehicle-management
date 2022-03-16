package model

import "gopkg.in/guregu/null.v4"

type Vehicle struct {
	ID          int64     `json:"id"`
	Model       string    `json:"model"`
	Color       string    `json:"color"`
	NumberPlate string    `json:"numberPlate"`
	UpdatedAt   string    `json:"updateAt"`
	CreatedAt   string    `json:"createdAt"`
	Name        string    `json:"name"`
	Launched    null.Bool `json:"launched"`
}

type DataField struct {
	Vehicle []*Vehicle `json:"vehicle"`
}
type DataFields struct {
	Vehicle *Vehicle `json:"vehicle"`
}
type VehicleDataField struct {
	Vehicle []*Vehicle `json:"vehicle"`
}

type AllVehicleModelResponse struct {
	Code   int       `json:"code"`
	Status string    `json:"status"`
	Data   DataField `json:"data"`
}
type VehicleModelResponse struct {
	Code   int        `json:"code"`
	Status string     `json:"status"`
	Data   DataFields `json:"data"`
}
