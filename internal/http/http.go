package http

import (
	"database/sql"
	"fmt"
	"strconv"

	"developer.zopsmart.com/go/gofr/pkg/gofr"
	"developer.zopsmart.com/go/gofr/pkg/gofr/types"
	"github.com/ishankochar09/go_pro/gofrTutorial/internal/models"
	"github.com/ishankochar09/go_pro/gofrTutorial/internal/service"

	gerror "developer.zopsmart.com/go/gofr/pkg/errors"
)

type Handler struct {
	serviceHandler service.Vehicle
}

func New(vehicle service.Vehicle) Handler {
	return Handler{serviceHandler: vehicle}
}

type vehicleResponse struct {
	Data interface{} `json:"vehicle"`
}

func (h Handler) GetByIDVehicle(ctx *gofr.Context) (interface{}, error) {
	vehID := ctx.PathParam("id")

	vehicleID, err := strconv.Atoi(vehID)

	if err != nil {
		return nil, gerror.MissingParam{Param: []string{"id"}}
	}

	res, err := h.serviceHandler.GetIDVehicle(ctx, vehicleID)

	if err != nil {
		return nil, err
	}

	resp := vehicleResponse{
		Data: res,
	}

	return types.Response{
		Data: resp,
	}, nil
}

func (h Handler) CreateVehicle(ctx *gofr.Context) (interface{}, error) {
	var vehicle models.Vehicle
	err := ctx.Bind(&vehicle)

	if err != nil {
		return nil, gerror.InvalidParam{Param: []string{"body"}}
	}

	res, err := h.serviceHandler.Create(ctx, &vehicle)

	if err != nil {
		return nil, err
	}

	resp := vehicleResponse{
		Data: res,
	}

	return types.Response{
		Data: resp,
	}, nil
}

func (h Handler) GetAllVehicles(ctx *gofr.Context) (interface{}, error) {
	res, err := h.serviceHandler.AllVehicles(ctx)
	fmt.Println(err)

	if err == sql.ErrNoRows {
		return nil, err
	}

	resp := vehicleResponse{
		Data: res,
	}

	return types.Response{
		Data: resp,
	}, nil
}

func (h Handler) DeleteVehicle(ctx *gofr.Context) (interface{}, error) {
	vehID := ctx.PathParam("id")
	vehicleID, err := strconv.Atoi(vehID)

	if err != nil {
		return nil, gerror.MissingParam{Param: []string{"id"}}
	}

	err = h.serviceHandler.DeleteIDVehicle(ctx, vehicleID)
	fmt.Print("Err : ", err)

	if err == sql.ErrNoRows {
		return nil, err
	}

	return "Deleted successfully", nil
}

func (h Handler) UpdateVehicle(ctx *gofr.Context) (interface{}, error) {
	vehID := ctx.PathParam("id")
	vehicleID, err := strconv.Atoi(vehID)

	if err != nil {
		return nil, gerror.MissingParam{Param: []string{"id"}}
	}

	var updatedVehicle models.Vehicle

	err = ctx.Bind(&updatedVehicle)
	if err != nil {
		return nil, gerror.InvalidParam{Param: []string{"body"}}
	}

	res, err := h.serviceHandler.UpdateIDVehicle(ctx, vehicleID, &updatedVehicle)

	if err != nil {
		return nil, err
	}

	resp := vehicleResponse{
		Data: res,
	}

	return types.Response{
		Data: resp,
	}, nil
}
