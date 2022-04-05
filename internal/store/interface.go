//go:generate mockgen -destination=interface_mock.go -package=store github.com/SN786/gofr_vms/internal/store VehicleManager
package store

import (
	"developer.zopsmart.com/go/gofr/pkg/gofr"
	"github.com/SN786/gofr_vms/internal/model"
)

type VehicleManager interface {
	GetDetailsByID(*gofr.Context, int) (*model.Vehicle, error)
	InsertVehicle(*gofr.Context, *model.Vehicle) (*model.Vehicle, error)
	DeleteVehicleByID(ctx *gofr.Context, id int) error
	UpdateVehicleByID(ctx *gofr.Context, vehicle *model.Vehicle) error
	GetAll(*gofr.Context) ([]*model.Vehicle, error)
}
