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
	if IsIDValid(id) {
		veh, err := v.datastore.GetDetailsByID(c, id)

		if err != nil {
			return &model.Vehicle{}, err
		}

		return veh, nil
	}

	return &model.Vehicle{}, errors.Error("validation error")
}

func (v StoreHandler) Post(c *gofr.Context, vehicle *model.Vehicle) (*model.Vehicle, error) {
	veh, err := v.datastore.InsertVehicle(c, vehicle)
	if err != nil {
		return &model.Vehicle{}, err
	}

	vehicleData, err := v.datastore.GetDetailsByID(c, int(veh.ID))
	if err != nil {
		return nil, err
	}

	return vehicleData, nil
}
func (v StoreHandler) Delete(c *gofr.Context, id int) error {
	if IsIDValid(id) {
		err := v.datastore.DeleteVehicleByID(c, id)
		if err != nil {
			return err
		}

		return nil
	}

	return errors.Error("validation error")
}

func (v StoreHandler) Update(c *gofr.Context, id int, vehicle *model.Vehicle) (*model.Vehicle, error) {
	if IsIDValid(id) {
		vehicle.ID = int64(id)
		err := v.datastore.UpdateVehicleByID(c, vehicle)

		if err != nil {
			return &model.Vehicle{}, err
		}

		return v.Get(c, id)
	}

	return nil, errors.Error("validation error")
}

func (v StoreHandler) GetAll(c *gofr.Context) ([]*model.Vehicle, error) {
	veh, err := v.datastore.GetAll(c)

	if err != nil {
		return nil, err
	}

	return veh, nil
}
