//go:generate mockgen -destination=interface_mock.go -package=store github.com/SukantArora/CRUD_Gofr/internal/store Vehicle
package store

import (
	"developer.zopsmart.com/go/gofr/pkg/gofr"
	"github.com/SukantArora/CRUD_Gofr/internal/models"
)

type Vehicle interface {
	GetByID(*gofr.Context, int) (*models.Vehicle, error)
	Get(*gofr.Context) ([]*models.Vehicle, error)
	Delete(*gofr.Context, int) error
	Update(*gofr.Context, int, *models.Vehicle) (*models.Vehicle, error)
	Create(*gofr.Context, *models.Vehicle) (*models.Vehicle, error)
}
