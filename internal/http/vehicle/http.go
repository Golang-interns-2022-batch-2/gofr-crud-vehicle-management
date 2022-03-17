package vehicle

import (
	"strconv"

	"developer.zopsmart.com/go/gofr/pkg/errors"
	"developer.zopsmart.com/go/gofr/pkg/gofr"
	"developer.zopsmart.com/go/gofr/pkg/gofr/types"
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
	idStr := c.PathParam("id")
	id, _ := strconv.Atoi(idStr)
	vehResponse, err := s.ServiceHandler.Get(c, id)

	if err != nil {
		return nil, err
	}

	res := types.Response{
		Data: model.DataFields{
			Vehicle: vehResponse,
		},
	}

	return res, nil
}

func (s Handler) Post(ctx *gofr.Context) (interface{}, error) {
	var vehicle model.Vehicle

	if err := ctx.Bind(&vehicle); err != nil {
		ctx.Logger.Errorf("error in binding: %v", err)
		return nil, errors.InvalidParam{}
	}

	vehicleData, err := s.ServiceHandler.Post(ctx, &vehicle)

	if err != nil {
		return nil, err
	}

	res := types.Response{
		Data: model.DataFields{
			Vehicle: vehicleData,
		},
	}

	return res, nil
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
	var vehicle model.Vehicle

	i := ctx.PathParam("id")
	id, _ := strconv.Atoi(i)

	if err := ctx.Bind(&vehicle); err != nil {
		ctx.Logger.Errorf("error in binding: %v", err)
		return nil, errors.InvalidParam{Param: []string{"body"}}
	}

	vehicleData, err := s.ServiceHandler.Update(ctx, id, &vehicle)

	if err != nil {
		return nil, err
	}

	res := types.Response{
		Data: model.DataFields{
			Vehicle: vehicleData,
		},
	}

	return res, nil
}

func (s Handler) GetAll(ctx *gofr.Context) (interface{}, error) {
	vehicleData, err := s.ServiceHandler.GetAll(ctx)

	if err != nil {
		return nil, err
	}

	res := types.Response{
		Data: model.DataField{
			Vehicle: vehicleData,
		},
	}

	return res, nil
}
