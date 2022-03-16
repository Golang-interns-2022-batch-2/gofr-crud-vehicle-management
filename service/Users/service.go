package vehicle

import (
	"github.com/Gofr-VMS/model"
	"github.com/Gofr-VMS/service"
	"github.com/Gofr-VMS/store"

	"developer.zopsmart.com/go/gofr/pkg/errors"
	"developer.zopsmart.com/go/gofr/pkg/gofr"
)

type services struct {
	st store.Store
}

func New(serv store.Store) service.Service {
	return &services{
		st: serv,
	}
}
func (s *services) GetByIdService(ctx *gofr.Context, id int) (*model.Vehicles, error) {
	if !IsIDValid(id) {
		return nil, errors.InvalidParam{Param: []string{"id"}}
	}

	return s.st.GetByID(ctx, id)
}

func (s *services) GetAllService(ctx *gofr.Context) ([]*model.Vehicles, error) {
	v, err := s.st.GetAll(ctx)
	if err != nil {
		return nil, err
	}

	return v, nil
}
func (s *services) InsertService(ctx *gofr.Context, vh *model.Vehicles) (*model.Vehicles, error) {
	if !validateModel(vh) {
		return nil, errors.Error("invalid  model field")
	}

	if !validateColor(vh) {
		return nil, errors.Error("invalid  color field")
	}

	if !validateNumPlate(vh) {
		return nil, errors.Error("invalid  numberPlate field")
	}

	if !validateName(vh) {
		return nil, errors.Error("invalid  name field")
	}

	return s.st.Insert(ctx, vh)
}
func (s *services) UpdateService(ctx *gofr.Context, vh *model.Vehicles) (*model.Vehicles, error) {
	if !IsIDValid(vh.ID) {
		return nil, errors.InvalidParam{Param: []string{"id"}}
	}

	v, err := s.st.Update(ctx, vh)

	if err != nil {
		return nil, err
	}

	return v, nil
}
func (s *services) RemoveService(ctx *gofr.Context, id int) error {
	if !IsIDValid(id) {
		return errors.InvalidParam{Param: []string{"id"}}
	}

	err := s.st.Delete(ctx, id)

	if err != nil {
		return errors.Error("error inavlid id ")
	}

	return nil
}
