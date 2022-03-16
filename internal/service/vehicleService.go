package service

import (
	"developer.zopsmart.com/go/gofr/pkg/gofr"
	"github.com/ishankochar09/go_pro/gofrTutorial/internal/models"
	"github.com/ishankochar09/go_pro/gofrTutorial/internal/store"

	gerror "developer.zopsmart.com/go/gofr/pkg/errors"
)

type VehicleService struct {
	vehStore store.VehicleStore
}

func New(dstore store.VehicleStore) *VehicleService {
	return &VehicleService{vehStore: dstore}
}

func (es *VehicleService) GetIDVehicle(ctx *gofr.Context, id int) (models.Vehicle, error) {
	if id < 0 {
		return models.Vehicle{}, gerror.Error("Validation Error")
	}

	res, err := es.vehStore.GetVehicleByID(ctx, id)

	if err != nil {
		return models.Vehicle{}, err
	}

	return res, err
}
func (es *VehicleService) Create(ctx *gofr.Context, veh *models.Vehicle) (models.Vehicle, error) {
	var vehi models.Vehicle

	var err error

	if veh.ID < 0 {
		return vehi, gerror.Error("invalid Id")
	}

	if veh.NumberPlate == "" {
		return vehi, gerror.Error("invalid numberPlate")
	}

	if veh.Model == "" {
		return vehi, gerror.Error("Invalid Model")
	}

	if veh.Name == "" {
		return vehi, gerror.Error("Invalid name")
	}

	res, err := es.vehStore.InsertVehicle(ctx, veh)

	return res, err
}
func (es *VehicleService) AllVehicles(ctx *gofr.Context) ([]models.Vehicle, error) {
	res, err := es.vehStore.GetAll(ctx)
	return res, err
}
func (es *VehicleService) DeleteIDVehicle(ctx *gofr.Context, id int) error {
	if id < 0 {
		return gerror.Error("Validation Error")
	}

	res, err := es.vehStore.GetVehicleByID(ctx, id)

	if err != nil {
		return gerror.Error("Internal Server Error")
	}

	err = es.vehStore.DeleteVehicle(ctx, res.ID)

	if err != nil {
		return err
	}

	return nil
}
func (es *VehicleService) UpdateIDVehicle(ctx *gofr.Context, id int, veh *models.Vehicle) (models.Vehicle, error) {
	var vehi models.Vehicle
	if veh.ID < 0 {
		return vehi, gerror.Error("Validation Error")
	}

	res, err1 := es.vehStore.UpdateVehByID(ctx, id, veh)

	return res, err1
}
