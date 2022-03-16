package main

import (
	http "github.com/Gofr-VMS/http/Users"
	service "github.com/Gofr-VMS/service/Users"
	store "github.com/Gofr-VMS/store/Users"

	"developer.zopsmart.com/go/gofr/pkg/gofr"
)

func main() {
	g := gofr.New()
	g.Server.ValidateHeaders = false

	st := store.New()
	sv := service.New(st)
	h := http.New(sv)

	g.GET("/vehicle/{id}", h.Get)
	g.GET("/vehicle", h.GetAll)
	g.POST("/vehicle", h.Create)
	g.PUT("/vehicle/{id}", h.Update)
	g.DELETE("/vehicle/{id}", h.Remove)

	g.Server.HTTP.Port = 8080

	g.Start()
}
