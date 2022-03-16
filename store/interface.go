package store

import (
	"developer.zopsmart.com/go/gofr/pkg/gofr"
	"github.com/Gofr-VMS/model"
)

type Store interface {
	GetByID(ctx *gofr.Context, id int) (*model.Vehicles, error)
	Insert(ctx *gofr.Context, v *model.Vehicles) (*model.Vehicles, error)
	GetAll(ctx *gofr.Context) ([]*model.Vehicles, error)
	Update(ctx *gofr.Context, v *model.Vehicles) (*model.Vehicles, error)
	Delete(ctx *gofr.Context, id int) error
}
