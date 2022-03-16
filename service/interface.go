package service

import (
	"github.com/Gofr-VMS/model"

	"developer.zopsmart.com/go/gofr/pkg/gofr"
)

type Service interface {
	GetByIdService(ctx *gofr.Context, id int) (*model.Vehicles, error)
	GetAllService(ctx *gofr.Context) ([]*model.Vehicles, error)
	InsertService(ctx *gofr.Context, vh *model.Vehicles) (*model.Vehicles, error)
	UpdateService(ctx *gofr.Context, v *model.Vehicles) (*model.Vehicles, error)
	RemoveService(ctx *gofr.Context, id int) error
}
