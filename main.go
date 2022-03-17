package main

import (
	"developer.zopsmart.com/go/gofr/pkg/gofr"

	hn "github.com/ishankochar09/go_pro/gofrTutorial/internal/http"
	src "github.com/ishankochar09/go_pro/gofrTutorial/internal/service"
	st "github.com/ishankochar09/go_pro/gofrTutorial/internal/store"
)

func main() {
	dataStore := st.New()
	service := src.New(dataStore)
	handler := hn.New(service)

	api := gofr.New()
	api.Server.ValidateHeaders = false
	api.GET("/vehicle/{id}", handler.GetByIDVehicle)
	api.POST("/vehicle", handler.CreateVehicle)
	api.GET("/vehicle", handler.GetAllVehicles)
	api.DELETE("/vehicle/{id}", handler.DeleteVehicle)
	api.PUT("/vehicle/{id}", handler.UpdateVehicle)
	api.Start()
}
