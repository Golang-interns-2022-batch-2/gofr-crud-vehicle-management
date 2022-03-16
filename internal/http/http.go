package http

import (
	"strconv"

	"developer.zopsmart.com/go/gofr/pkg/errors"
	"developer.zopsmart.com/go/gofr/pkg/gofr"
	"developer.zopsmart.com/go/gofr/pkg/gofr/types"
	"github.com/SukantArora/CRUD_Gofr/internal/models"
	"github.com/SukantArora/CRUD_Gofr/internal/service"
)

type response struct {
	Data interface{} `json:"Vehicle"`
}

type APIHandler struct {
	serviceHandler service.Vehicle
}

func New(vehicle service.Vehicle) *APIHandler {
	return &APIHandler{serviceHandler: vehicle}
}

func (h *APIHandler) Create(ctx *gofr.Context) (interface{}, error) {
	var Vehicle *models.Vehicle

	var respVehicle *models.Vehicle

	err := ctx.Bind(&Vehicle)

	if err != nil {
		return nil, errors.InvalidParam{Param: []string{"Body"}}
	}

	respVehicle, err = h.serviceHandler.Create(ctx, Vehicle)

	if err != nil {
		return nil, err
	}

	resp := types.Response{
		Data: response{
			Data: respVehicle,
		},
	}

	return resp, nil
}

func (h *APIHandler) GetByID(ctx *gofr.Context) (interface{}, error) {
	vehID := ctx.PathParam("id")
	vehIDint, converr := strconv.Atoi(vehID)

	if converr != nil {
		return nil, errors.InvalidParam{Param: []string{"id"}}
	}

	respVehicle, err := h.serviceHandler.GetByID(ctx, vehIDint)

	if err == nil {
		resp := types.Response{
			Data: response{
				Data: respVehicle,
			},
		}
		return resp, nil
	}

	return nil, err
}

func (h *APIHandler) Get(ctx *gofr.Context) (interface{}, error) {
	data, err := h.serviceHandler.Get(ctx)
	if err == nil {
		resp := types.Response{
			Data: response{
				Data: data,
			},
		}
		return resp, nil
	}

	return nil, err
}

func (h *APIHandler) Delete(ctx *gofr.Context) (interface{}, error) {
	empID := ctx.PathParam("id")
	empIDint, converr := strconv.Atoi(empID)

	if converr != nil {
		return nil, errors.InvalidParam{Param: []string{"id"}}
	}

	err := h.serviceHandler.Delete(ctx, empIDint)

	if err != nil {
		return nil, err
	}

	return nil, nil
}

func (h *APIHandler) Update(ctx *gofr.Context) (interface{}, error) {
	vehID := ctx.PathParam("id")
	vehIDint, converr := strconv.Atoi(vehID)

	if converr != nil {
		return nil, errors.InvalidParam{Param: []string{"id"}}
	}

	var updatedVehicle *models.Vehicle

	var respVehicle *models.Vehicle

	err := ctx.Bind(&updatedVehicle)

	if err != nil {
		return nil, errors.InvalidParam{Param: []string{"Body"}}
	}

	respVehicle, err = h.serviceHandler.Update(ctx, vehIDint, updatedVehicle)

	if err != nil {
		return nil, err
	}

	resp := types.Response{
		Data: response{
			Data: respVehicle,
		},
	}

	return resp, nil
}
