//go:generate mockgen -destination=service_mock.go -package=service github.com/ishankochar09/go_pro/gofrTutorial/service VehicleInterface
package service

import (
	"developer.zopsmart.com/go/gofr/pkg/gofr"
	"github.com/ishankochar09/go_pro/gofrTutorial/internal/models"
)

type VehicleInterface interface {
	GetIDVehicle(*gofr.Context, int) (models.Vehicle, error)
	Create(*gofr.Context, *models.Vehicle) (models.Vehicle, error)
	AllVehicles(*gofr.Context) ([]models.Vehicle, error)
	DeleteIDVehicle(*gofr.Context, int) error
	UpdateIDVehicle(*gofr.Context, int, *models.Vehicle) (models.Vehicle, error)
}
