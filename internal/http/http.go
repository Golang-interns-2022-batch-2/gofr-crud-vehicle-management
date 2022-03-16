package http

import (
	"database/sql"
	"fmt"
	"strconv"

	gerror "developer.zopsmart.com/go/gofr/pkg/errors"
	"developer.zopsmart.com/go/gofr/pkg/gofr"
	"github.com/ishankochar09/go_pro/gofrTutorial/internal/models"
	"github.com/ishankochar09/go_pro/gofrTutorial/internal/service"
)

type Handler struct {
	reqService service.VehicleInterface
}

func New(user service.VehicleInterface) Handler {
	return Handler{reqService: user}
}

type vehicleDetails struct {
	Data models.Vehicle `json:"Vehicle"`
}

type vehiclesList struct {
	Data []models.Vehicle `json:"Vehicle"`
}

type vehicleResponse struct {
	Code   int            `json:"Code"`
	Status string         `json:"Status"`
	Data   vehicleDetails `json:"Data"`
}

type allVehicleResponse struct {
	Code   int          `json:"Code"`
	Status string       `json:"Status"`
	Data   vehiclesList `json:"Data"`
}
type messageError struct {
	Code    int    `json:"code"`
	Status  string `json:"status"`
	Message string `json:"message"`
}

func (hn Handler) GetByIDVehicle(ctx *gofr.Context) (interface{}, error) {
	vehID := ctx.PathParam("id")

	vehicleID, err := strconv.Atoi(vehID)

	if err != nil {
		return nil, gerror.MissingParam{Param: []string{"id"}}
	}

	vehicle, err := hn.reqService.GetIDVehicle(ctx, vehicleID)

	if err != nil {
		resp := messageError{
			Code:    400,
			Status:  "ERROR",
			Message: err.Error(),
		}

		return resp, err
	}

	data := vehicleDetails{
		Data: vehicle,
	}

	resp := vehicleResponse{
		Code:   200,
		Status: "SUCCESS",
		Data:   data,
	}

	return resp, nil
}

func (hn Handler) CreateVehicle(ctx *gofr.Context) (interface{}, error) {
	var vehicle models.Vehicle
	err := ctx.Bind(&vehicle)

	if err != nil {
		return nil, gerror.InvalidParam{Param: []string{"body"}}
	}

	var ans models.Vehicle

	ans, err = hn.reqService.Create(ctx, &vehicle)

	if err != nil {
		resp := messageError{
			Code:    400,
			Status:  "ERROR",
			Message: err.Error(),
		}

		return resp, err
	}

	data := vehicleDetails{
		Data: ans,
	}

	resp := vehicleResponse{
		Code:   200,
		Status: "SUCCESS",
		Data:   data,
	}

	return resp, nil
}

func (hn Handler) GetAllVehicles(ctx *gofr.Context) (interface{}, error) {
	res, err := hn.reqService.AllVehicles(ctx)
	fmt.Println(err)

	if err == sql.ErrNoRows {
		resp := messageError{
			Code:    400,
			Status:  "ERROR",
			Message: err.Error(),
		}

		return resp, err
	}

	data := vehiclesList{
		Data: res,
	}

	resp := allVehicleResponse{
		Code:   200,
		Status: "SUCCESS",
		Data:   data,
	}

	return resp, nil
}

func (hn Handler) DeleteVehicle(ctx *gofr.Context) (interface{}, error) {
	vehID := ctx.PathParam("id")
	vehicleID, err := strconv.Atoi(vehID)

	if err != nil {
		return nil, gerror.MissingParam{Param: []string{"id"}}
	}

	err = hn.reqService.DeleteIDVehicle(ctx, vehicleID)
	fmt.Print("Err : ", err)

	if err == sql.ErrNoRows {
		resp := messageError{
			Code:    400,
			Status:  "ERROR",
			Message: err.Error(),
		}

		return resp, err
	}

	return "Deleted successfully", nil
}

func (hn Handler) UpdateVehicle(ctx *gofr.Context) (interface{}, error) {
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

	res, err := hn.reqService.UpdateIDVehicle(ctx, vehicleID, &updatedVehicle)

	if err != nil {
		resp := messageError{
			Code:    400,
			Status:  "ERROR",
			Message: err.Error(),
		}

		return resp, err
	}

	data := vehicleDetails{
		Data: res,
	}

	resp := vehicleResponse{
		Code:   200,
		Status: "SUCCESS",
		Data:   data,
	}

	return resp, nil
}
