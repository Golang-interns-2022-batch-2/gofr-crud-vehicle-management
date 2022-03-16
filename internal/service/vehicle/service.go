package vehicle

import (
	"developer.zopsmart.com/go/gofr/pkg/errors"
	"developer.zopsmart.com/go/gofr/pkg/gofr"
	"github.com/SN786/gofr_vms/internal/model"
	"github.com/SN786/gofr_vms/internal/store"
)

type StoreHandler struct {
	datastore store.VehicleManager
}

func New(v store.VehicleManager) *StoreHandler {
	return &StoreHandler{datastore: v}
}
func (v StoreHandler) Get(c *gofr.Context, id int) (*model.Vehicle, error) {
	if id > 0 {
		vehicleData, err := v.datastore.GetDetailsByID(c, id)

		if err != nil {
			return &model.Vehicle{}, err
		}

		return vehicleData, nil
	}

	return nil, errors.InvalidParam{Param: []string{"id"}}
}

func (v StoreHandler) Post(c *gofr.Context, vehicle *model.Vehicle) (*model.Vehicle, error) {
	vehicleData, err := v.datastore.InsertVehicle(c, vehicle)
	if err != nil {
		return &model.Vehicle{}, err
	}

	vehicleData, err = v.datastore.GetDetailsByID(c, int(vehicle.ID))
	if err != nil {
		return nil, err
	}

	return vehicleData, nil
}
func (v StoreHandler) Delete(c *gofr.Context, id int) error {
	if id > 0 {
		err := v.datastore.DeleteVehicleByID(c, id)
		if err != nil {
			return err
		}

		return nil
	}

	return errors.InvalidParam{Param: []string{"id"}}
}

func (v StoreHandler) Update(c *gofr.Context, id int, vehicle *model.Vehicle) (*model.Vehicle, error) {
	if id > 0 {
		vehicle.ID = int64(id)
		err := v.datastore.UpdateVehicleByID(c, vehicle)

		if err != nil {
			return &model.Vehicle{}, err
		}

		return v.Get(c, id)
	}

	return nil, errors.InvalidParam{Param: []string{"id"}}
}

func (v StoreHandler) GetAll(c *gofr.Context) ([]*model.Vehicle, error) {
	vehicleData, err := v.datastore.GetAll(c)

	if err != nil {
		return nil, err
	}

	return vehicleData, nil
}
