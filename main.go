package main

import (
	"developer.zopsmart.com/go/gofr/pkg/gofr"
	api "github.com/SukantArora/CRUD_Gofr/internal/http"
	vehicleService "github.com/SukantArora/CRUD_Gofr/internal/service/vehicle"
	vehicleStore "github.com/SukantArora/CRUD_Gofr/internal/store/vehicle"
)

func main() {
	dataStore := vehicleStore.New()
	vehServ := vehicleService.New(dataStore)
	handler := api.New(vehServ)
	app := gofr.New()
	app.Server.ValidateHeaders = false
	app.GET("/vehicles/{id}", handler.GetByID)
	app.POST("/vehicles", handler.Create)
	app.DELETE("/vehicles/{id}", handler.Delete)
	app.PUT("/vehicles/{id}", handler.Update)
	app.GET("/vehicles", handler.Get)
	app.Start()
}
