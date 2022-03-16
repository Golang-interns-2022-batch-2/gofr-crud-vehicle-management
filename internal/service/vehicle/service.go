package vehicle

import (
	"developer.zopsmart.com/go/gofr/pkg/errors"
	"developer.zopsmart.com/go/gofr/pkg/gofr"
	"github.com/SukantArora/CRUD_Gofr/internal/models"
	"github.com/SukantArora/CRUD_Gofr/internal/store"
)

type Service struct {
	VehicleStoreHandler store.Vehicle
}

func New(st store.Vehicle) *Service {
	return &Service{
		VehicleStoreHandler: st,
	}
}

func (vs *Service) GetByID(ctx *gofr.Context, id int) (*models.Vehicle, error) {
	if !validateID(id) {
		return nil, errors.InvalidParam{Param: []string{"id"}}
	}

	return vs.VehicleStoreHandler.GetByID(ctx, id)
}

func (vs *Service) Get(ctx *gofr.Context) ([]*models.Vehicle, error) {
	vehicles, err := vs.VehicleStoreHandler.Get(ctx)
	if err != nil {
		return nil, err
	}

	return vehicles, nil
}

func (vs *Service) Delete(ctx *gofr.Context, id int) error {
	if !validateID(id) {
		return errors.InvalidParam{Param: []string{"id"}}
	}

	err := vs.VehicleStoreHandler.Delete(ctx, id)

	return err
}

func (vs *Service) Update(ctx *gofr.Context, id int, vehicle *models.Vehicle) (*models.Vehicle, error) {
	if !validateID(id) {
		return nil, errors.InvalidParam{Param: []string{"id"}}
	}

	return vs.VehicleStoreHandler.Update(ctx, id, vehicle)
}
func (vs *Service) Create(ctx *gofr.Context, vehicle *models.Vehicle) (*models.Vehicle, error) {
	ok, err := validateInput(vehicle)
	if !ok {
		return nil, err
	}

	return vs.VehicleStoreHandler.Create(ctx, vehicle)
}
