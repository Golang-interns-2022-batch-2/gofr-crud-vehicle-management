package vehicle

import (
	"net/http"
	"strconv"

	"github.com/Gofr-VMS/model"
	"github.com/Gofr-VMS/service"

	"developer.zopsmart.com/go/gofr/pkg/errors"
	"developer.zopsmart.com/go/gofr/pkg/gofr"
	"developer.zopsmart.com/go/gofr/pkg/gofr/types"
)

type Handler struct {
	serv service.Service
}

func New(s service.Service) Handler {
	return Handler{
		serv: s,
	}
}

const (
	badReqCode = 400
)

func (h Handler) Get(ctx *gofr.Context) (interface{}, error) {
	ID := ctx.PathParam("id")

	if ID == "" {
		return nil, errors.MissingParam{Param: []string{"id"}}
	}

	id, err := strconv.Atoi(ID)
	if err != nil {
		return nil, errors.InvalidParam{Param: []string{"id"}}
	}

	resp, err := h.serv.GetByIdService(ctx, id)
	if err != nil {
		return errors.Response{
			StatusCode:    badReqCode,
			Code:  "ERROR",
			Reason: err.Error(),
		}, nil
	}

	return types.Response{
		Data: model.RespnseData{Data: resp},
	}, nil
}
func (h Handler) GetAll(ctx *gofr.Context) (interface{}, error) {
	resp, err := h.serv.GetAllService(ctx)
	if err != nil {
		return errors.Response{
			StatusCode:    badReqCode,
			Code:  "ERROR",
			Reason: err.Error(),
		}, nil
	}

	return types.Response{
		Data: model.RespnseData{Data: resp},
	}, nil
}
func (h Handler) Remove(ctx *gofr.Context) (interface{}, error) {
	ID := ctx.PathParam("id")

	if ID == "" {
		return nil, errors.MissingParam{Param: []string{"id"}}
	}

	id, err := strconv.Atoi(ID)
	if err != nil {
		return nil, errors.InvalidParam{Param: []string{"id"}}
	}

	err = h.serv.RemoveService(ctx, id)

	if err != nil {
		return errors.Response{
			StatusCode:    badReqCode,
			Code:  "ERROR",
			Reason: err.Error(),
		}, nil
	}
	return errors.Response{
		StatusCode:    http.StatusOK,
		Code:  "Success",
		Reason: "vehicle deleted successfully",
	}, nil

	
}
func (h Handler) Create(ctx *gofr.Context) (interface{}, error) {
	var veh *model.Vehicles

	if err := ctx.Bind(&veh); err != nil {
		return nil, errors.InvalidParam{Param: []string{"body"}}
	}
	r, e := h.serv.InsertService(ctx, veh)

	if e != nil {
		return errors.Response{
			StatusCode:    badReqCode,
			Code:  "ERROR",
			Reason: e.Error(),
		}, nil
	}
	return types.Response{
		Data: model.RespnseData{Data: r},
	}, nil

}
func (h Handler) Update(ctx *gofr.Context) (interface{}, error) {
	ID := ctx.PathParam("id")

	if ID == "" {
		return nil, errors.MissingParam{Param: []string{"id"}}
	}

	id, err := strconv.Atoi(ID)
	if err != nil {
		return nil, errors.InvalidParam{Param: []string{"id"}}
	}

	var veh model.Vehicles
	veh.ID = id

	if err := ctx.Bind(&veh); err != nil {
		return nil, err
	}

	r, e := h.serv.UpdateService(ctx, &veh)
	if e != nil {
		return errors.Response{
			StatusCode:    badReqCode,
			Code:  "ERROR",
			Reason: e.Error(),
		}, nil
	}

	return types.Response{
		Data: model.RespnseData{Data: r},
	}, nil
}
