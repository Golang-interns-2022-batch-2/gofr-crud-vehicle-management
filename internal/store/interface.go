//go:generate mockgen -destination=interface_mock.go -package=store github.com/ishankochar09/go_pro/gofrTutorial/store VehicleStore
package store

import (
	"developer.zopsmart.com/go/gofr/pkg/gofr"
	"github.com/ishankochar09/go_pro/gofrTutorial/internal/models"
)

type VehicleStore interface {
	GetVehicleByID(*gofr.Context, int) (models.Vehicle, error)
	InsertVehicle(*gofr.Context, *models.Vehicle) (models.Vehicle, error)
	GetAll(*gofr.Context) ([]models.Vehicle, error)
	DeleteVehicle(*gofr.Context, int) error
	UpdateVehByID(*gofr.Context, int, *models.Vehicle) (models.Vehicle, error)
}
