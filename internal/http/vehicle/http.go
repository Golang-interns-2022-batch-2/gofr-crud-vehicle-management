package vehicle

import (
	"strconv"

	"developer.zopsmart.com/go/gofr/pkg/errors"
	"developer.zopsmart.com/go/gofr/pkg/gofr"
	"github.com/SN786/gofr_vms/internal/model"
	"github.com/SN786/gofr_vms/internal/service"
)

type Handler struct {
	ServiceHandler service.VehicleManager
}

func New(s service.VehicleManager) *Handler {
	return &Handler{ServiceHandler: s}
}

func (s Handler) Get(c *gofr.Context) (interface{}, error) {
	i := c.PathParam("id")
	id, _ := strconv.Atoi(i)
	vehResponse, err := s.ServiceHandler.Get(c, id)

	if err != nil {
		return nil, err
	}

	df := model.DataFields{
		Vehicle: vehResponse,
	}
	resp := model.VehicleModelResponse{
		Code:   200,
		Status: "SUCCESS",
		Data:   df,
	}

	return resp, nil
}

func (s Handler) Post(ctx *gofr.Context) (interface{}, error) {
	var veh model.Vehicle

	if err := ctx.Bind(&veh); err != nil {
		ctx.Logger.Errorf("error in binding: %v", err)
		return nil, errors.InvalidParam{Param: []string{"body"}}
	}

	resp, err := s.ServiceHandler.Post(ctx, &veh)

	if err != nil {
		return nil, err
	}

	df := model.DataFields{
		Vehicle: resp,
	}
	vehResponse := model.VehicleModelResponse{
		Code:   200,
		Status: "SUCCESS",
		Data:   df,
	}

	return vehResponse, nil
}

func (s Handler) Delete(c *gofr.Context) (interface{}, error) {
	type response struct {
		code    int
		status  string
		message string
	}

	i := c.PathParam("id")
	id, _ := strconv.Atoi(i)
	err := s.ServiceHandler.Delete(c, id)

	if err != nil {
		return nil, err
	}

	return response{code: 200, status: "success", message: "successfully deleted"}, nil
}

func (s Handler) Update(ctx *gofr.Context) (interface{}, error) {
	var veh model.Vehicle

	i := ctx.PathParam("id")
	id, _ := strconv.Atoi(i)

	if err := ctx.Bind(&veh); err != nil {
		ctx.Logger.Errorf("error in binding: %v", err)
		return nil, errors.InvalidParam{Param: []string{"body"}}
	}

	vehResponse, err := s.ServiceHandler.Update(ctx, id, &veh)

	if err != nil {
		return nil, err
	}

	df := model.DataFields{
		Vehicle: vehResponse,
	}
	resp := model.VehicleModelResponse{
		Code:   200,
		Status: "SUCCESS",
		Data:   df,
	}

	return resp, nil
}

func (s Handler) GetAll(ctx *gofr.Context) (interface{}, error) {
	vehResponse, err := s.ServiceHandler.GetAll(ctx)

	if err != nil {
		return nil, err
	}

	dataField := model.DataField{
		Vehicle: vehResponse,
	}
	resp := model.AllVehicleModelResponse{
		Code:   200,
		Status: "SUCCESS",
		Data:   dataField,
	}

	return resp, nil
}
