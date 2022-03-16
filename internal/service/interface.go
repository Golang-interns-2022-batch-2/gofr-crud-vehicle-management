package service

//go:generate mockgen -destination=interface_mock.go -package=service github.com/SN786/gofr_vms/internal/service VehicleManager
import (
	"developer.zopsmart.com/go/gofr/pkg/gofr"
	"github.com/SN786/gofr_vms/internal/model"
)

type VehicleManager interface {
	Get(*gofr.Context, int) (*model.Vehicle, error)
	Post(*gofr.Context, *model.Vehicle) (*model.Vehicle, error)
	Delete(ctx *gofr.Context, id int) error
	Update(*gofr.Context, int, *model.Vehicle) (*model.Vehicle, error)
	GetAll(*gofr.Context) ([]*model.Vehicle, error)
}
