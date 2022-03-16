package main

import (
	_ "github.com/go-sql-driver/mysql"

	"developer.zopsmart.com/go/gofr/pkg/gofr"

	hthandler "github.com/SN786/gofr_vms/internal/http/vehicle"
	serviceVehicle "github.com/SN786/gofr_vms/internal/service/vehicle"
	StoreHandler "github.com/SN786/gofr_vms/internal/store/vehicle"
)

func main() {
	datastore := StoreHandler.New()

	serviceHandler := serviceVehicle.New(datastore)

	httpHandler := hthandler.New(serviceHandler)

	k := gofr.New()
	k.Server.ValidateHeaders = false

	k.GET("/vehicles/{id}", httpHandler.Get)
	k.GET("/vehicles", httpHandler.GetAll)
	k.POST("/vehicles", httpHandler.Post)
	k.DELETE("/vehicles/{id}", httpHandler.Delete)
	k.PUT("/vehicles/{id}", httpHandler.Update)
	k.Start()
}
