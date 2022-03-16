package http

import (
	"database/sql"
	"encoding/json"
	"net/http"
	"strconv"

	"developer.zopsmart.com/go/gofr/pkg/errors"
	"developer.zopsmart.com/go/gofr/pkg/gofr"
	"github.com/SukantArora/CRUD_Gofr/internal/models"
	"github.com/SukantArora/CRUD_Gofr/internal/service"
)

type errorMsg struct {
	Code    int    `json:"code"`
	Status  string `json:"status"`
	Message string `jsong:"message"`
}

type response struct {
	Code   int         `json:"code"`
	Status string      `json:"status"`
	Data   interface{} `json:"data"`
}

func writeErrorMsg(code int, err error) []byte {
	resp := errorMsg{
		Code:    code,
		Status:  "ERROR",
		Message: err.Error(),
	}
	encodedData, _ := json.Marshal(resp)

	return encodedData
}

func writeResponse(result interface{}) response {
	resp := response{
		Code:   http.StatusOK,
		Status: "SUCCESS",
		Data:   result,
	}

	return resp
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
		encodedData := writeErrorMsg(http.StatusBadRequest, err)
		return encodedData, errors.InvalidParam{Param: []string{"Body"}}
	}

	respVehicle, err = h.serviceHandler.Create(ctx, Vehicle)

	if err != nil {
		if err.Error() == "invalid data format" {
			encodedData := writeErrorMsg(http.StatusBadRequest, err)
			return encodedData, err
		}

		encodedData := writeErrorMsg(http.StatusInternalServerError, err)

		return encodedData, err
	}

	resp := writeResponse(respVehicle)

	return resp, nil
}

func (h *APIHandler) GetByID(ctx *gofr.Context) (interface{}, error) {
	vehID := ctx.PathParam("id")
	vehIDint, converr := strconv.Atoi(vehID)

	if converr != nil {
		encodedData := writeErrorMsg(http.StatusBadRequest, errors.InvalidParam{Param: []string{"id"}})
		return encodedData, errors.InvalidParam{Param: []string{"id"}}
	}

	respVehicle, err := h.serviceHandler.GetByID(ctx, vehIDint)

	if err == nil {
		resp := writeResponse(respVehicle)
		return resp, nil
	}

	if err == sql.ErrNoRows {
		encodedData := writeErrorMsg(http.StatusNotFound, err)
		return encodedData, err
	}

	encodedData := writeErrorMsg(http.StatusInternalServerError, err)

	return encodedData, err
}

func (h *APIHandler) Get(ctx *gofr.Context) (interface{}, error) {
	data, err := h.serviceHandler.Get(ctx)
	if err == nil {
		resp := writeResponse(data)
		return resp, nil
	}

	if err == sql.ErrNoRows {
		encodedData := writeErrorMsg(http.StatusNotFound, err)
		return encodedData, err
	}

	encodedData := writeErrorMsg(http.StatusInternalServerError, errors.Error("internal server error"))

	return encodedData, err
}

func (h *APIHandler) Delete(ctx *gofr.Context) (interface{}, error) {
	empID := ctx.PathParam("id")
	empIDint, converr := strconv.Atoi(empID)

	if converr != nil {
		encodedData := writeErrorMsg(http.StatusBadRequest, errors.InvalidParam{Param: []string{"id"}})
		return encodedData, errors.InvalidParam{Param: []string{"id"}}
	}

	err := h.serviceHandler.Delete(ctx, empIDint)

	if err != nil {
		if err.Error() == "error Invalid ID" {
			encodedData := writeErrorMsg(http.StatusBadRequest, err)
			return encodedData, err
		} else if err == sql.ErrNoRows {
			encodedData := writeErrorMsg(http.StatusNotFound, err)
			return encodedData, err
		} else {
			encodedData := writeErrorMsg(http.StatusInternalServerError, err)
			return encodedData, err
		}
	} else {
		resp := writeResponse("Vehicle deleted successfully.")
		return resp, nil
	}
}

func (h *APIHandler) Update(ctx *gofr.Context) (interface{}, error) {
	vehID := ctx.PathParam("id")
	vehIDint, converr := strconv.Atoi(vehID)

	if converr != nil {
		encodedData := writeErrorMsg(http.StatusBadRequest, errors.InvalidParam{Param: []string{"id"}})
		return encodedData, errors.InvalidParam{Param: []string{"id"}}
	}

	var updatedVehicle *models.Vehicle

	var respVehicle *models.Vehicle

	err := ctx.Bind(&updatedVehicle)
	if err != nil {
		encodedData := writeErrorMsg(http.StatusBadRequest, err)
		return encodedData, errors.InvalidParam{Param: []string{"Body"}}
	}

	respVehicle, err = h.serviceHandler.Update(ctx, vehIDint, updatedVehicle)
	if err != nil {
		if err == sql.ErrNoRows {
			encodedData := writeErrorMsg(http.StatusNotFound, err)
			return encodedData, err
		}

		encodedData := writeErrorMsg(http.StatusInternalServerError, err)

		return encodedData, err
	}

	resp := writeResponse(respVehicle)

	return resp, nil
}
