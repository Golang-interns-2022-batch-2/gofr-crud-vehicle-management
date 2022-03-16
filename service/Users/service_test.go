package vehicle

import (
	"database/sql"
	"reflect"
	"testing"
	"time"

	"github.com/Gofr-VMS/model"
	"github.com/Gofr-VMS/service"
	"github.com/Gofr-VMS/store"

	"developer.zopsmart.com/go/gofr/pkg/errors"
	"developer.zopsmart.com/go/gofr/pkg/gofr"
	"github.com/golang/mock/gomock"
)

func setMock(t *testing.T) (*gofr.Gofr, service.Service, *store.MockStore) {
	app := gofr.New()
	ctrl := gomock.NewController(t)
	mock := store.NewMockStore(ctrl)
	s := New(mock)

	return app, s, mock
}

func TestGtByIDService(t *testing.T) {
	app, s, mock := setMock(t)

	TC := []struct {
		desc   string
		input  int
		output *model.Vehicles
		experr error
		mock   []*gomock.Call
	}{
		{
			desc:  "success",
			input: 1,
			output: &model.Vehicles{
				ID:        1,
				Model:     "i8",
				Color:     "Black",
				NumPlate:  "MH 03 AT 007",
				CreatedAt: time.Date(2022, 1, 1, 1, 1, 1, 1, time.UTC),
				UpdatedAt: time.Date(2022, 1, 1, 1, 1, 1, 1, time.UTC),
				Name:      "BMW",
				Launched:  true,
			},
			experr: nil,
			mock: []*gomock.Call{
				mock.EXPECT().GetByID(gomock.Any(), 1).Return(&model.Vehicles{
					ID:        1,
					Model:     "i8",
					Color:     "Black",
					NumPlate:  "MH 03 AT 007",
					CreatedAt: time.Date(2022, 1, 1, 1, 1, 1, 1, time.UTC),
					UpdatedAt: time.Date(2022, 1, 1, 1, 1, 1, 1, time.UTC),
					Name:      "BMW",
					Launched:  true,
				}, nil),
			},
		},
		{
			desc:   "fail",
			input:  -1,
			experr: errors.InvalidParam{Param: []string{"id"}},
		},
	}
	for _, tc := range TC {
		ctx := gofr.NewContext(nil, nil, app)
		r, err := s.GetByIdService(ctx, tc.input)

		if !reflect.DeepEqual(err, tc.experr) {
			t.Errorf("%s : expected %v, got %v ", tc.desc, tc.experr, err)
		}
		if tc.experr == nil && !reflect.DeepEqual(r, tc.output) {
			t.Errorf("%s : expected %v  , got: %v ", tc.desc, tc.output, r)
		}
	}
}
func TestGetAllService(t *testing.T) {
	app, s, mock := setMock(t)
	TC := []struct {
		desc   string
		output []*model.Vehicles
		experr error
		mock   []*gomock.Call
	}{
		{
			desc: "success",
			output: []*model.Vehicles{
				{
					ID:        1,
					Model:     "Q5",
					Color:     "Black",
					NumPlate:  "UK 07 1896",
					CreatedAt: time.Date(2022, 1, 1, 1, 1, 1, 1, time.UTC),
					UpdatedAt: time.Date(2022, 1, 1, 1, 1, 1, 1, time.UTC),
					Name:      "Audi",
					Launched:  true,
				},
			},
			experr: nil,
			mock: []*gomock.Call{
				mock.EXPECT().GetAll(gomock.Any()).Return([]*model.Vehicles{
					{
						ID:        1,
						Model:     "Q5",
						Color:     "Black",
						NumPlate:  "UK 07 1896",
						CreatedAt: time.Date(2022, 1, 1, 1, 1, 1, 1, time.UTC),
						UpdatedAt: time.Date(2022, 1, 1, 1, 1, 1, 1, time.UTC),
						Name:      "Audi",
						Launched:  true,
					},
				}, nil),
			},
		},
		{
			desc:   "fail",
			output: nil,
			experr: sql.ErrNoRows,
			mock: []*gomock.Call{
				mock.EXPECT().GetAll(gomock.Any()).Return(nil, sql.ErrNoRows),
			},
		},
	}

	for _, tc := range TC {
		ctx := gofr.NewContext(nil, nil, app)
		r, err := s.GetAllService(ctx)

		if !reflect.DeepEqual(err, tc.experr) {
			t.Errorf("%s : expected %v, got %v ", tc.desc, tc.experr, err)
		}
		if tc.experr == nil && !reflect.DeepEqual(r, tc.output) {
			t.Errorf("%s : expected %v  , got: %v ", tc.desc, tc.output, r)
		}
	}
}
func TestRemoveService(t *testing.T) {
	app, s, mock := setMock(t)

	TC := []struct {
		desc   string
		input  int
		experr error
		mock   []*gomock.Call
	}{
		{
			desc:   "success",
			input:  1,
			experr: nil,
			mock: []*gomock.Call{
				mock.EXPECT().Delete(gomock.Any(), 1).Return(nil),
			},
		},
		{
			desc:   "fail",
			input:  -1,
			experr: errors.InvalidParam{Param: []string{"id"}},
			mock:   nil,
		},
	}
	for _, tc := range TC {
		ctx := gofr.NewContext(nil, nil, app)
		err := s.RemoveService(ctx, tc.input)

		if !reflect.DeepEqual(err, tc.experr) {
			t.Errorf("%s : expected %v, got %v ", tc.desc, tc.experr, err)
		}
	}
}
func TestUpdateService(t *testing.T) {
	app, s, mock := setMock(t)
	TC := []struct {
		desc   string
		input  *model.Vehicles
		experr error
		mock   []*gomock.Call
	}{
		{
			desc: "success",
			input: &model.Vehicles{
				ID:       1,
				Model:    "Q5",
				Color:    "Black",
				NumPlate: "UK 07 1896",
				Name:     "Audi",
				Launched: true,
			},
			experr: nil,
			mock: []*gomock.Call{
				mock.EXPECT().Update(gomock.Any(), &model.Vehicles{
					ID:       1,
					Model:    "Q5",
					Color:    "Black",
					NumPlate: "UK 07 1896",
					Name:     "Audi",
					Launched: true,
				}).Return(&model.Vehicles{
					ID:       1,
					Model:    "Q5",
					Color:    "Black",
					NumPlate: "UK 07 1896",
					Name:     "Audi",
					Launched: true,
				}, nil),
			},
		},
		{

			desc: "fail",
			input: &model.Vehicles{
				ID: -1,
			},
			experr: errors.InvalidParam{Param: []string{"id"}},
			mock:   nil,
		},
	}

	for _, tc := range TC {
		ctx := gofr.NewContext(nil, nil, app)
		r, err := s.UpdateService(ctx, tc.input)

		if !reflect.DeepEqual(err, tc.experr) {
			t.Errorf("%s : expected %v, got %v ", tc.desc, tc.experr, err)
		}

		if tc.experr == nil && !reflect.DeepEqual(r, tc.input) {
			t.Errorf("%s : expected %v, but got %v", tc.desc, tc.input, r)
		}
	}
}
func TestInsertService(t *testing.T) {
	app, s, mock := setMock(t)
	TC := []struct {
		desc   string
		input  *model.Vehicles
		experr error
		mock   []*gomock.Call
	}{
		{
			desc: "success",
			input: &model.Vehicles{
				Model:    "Q5",
				Color:    "Black",
				NumPlate: "UK 07 1896",
				Name:     "Audi",
				Launched: true,
			},
			experr: nil,
			mock: []*gomock.Call{
				mock.EXPECT().Insert(gomock.Any(), &model.Vehicles{

					Model:    "Q5",
					Color:    "Black",
					NumPlate: "UK 07 1896",
					Name:     "Audi",
					Launched: true,
				}).Return(&model.Vehicles{

					Model:    "Q5",
					Color:    "Black",
					NumPlate: "UK 07 1896",
					Name:     "Audi",
					Launched: true,
				}, nil),
			},
		},
		{
			desc: "fail",
			input: &model.Vehicles{
				ID: -1,
			},
			experr: errors.Error("invalid  model field"),
			mock:   nil,
		},
	}

	for _, tc := range TC {
		ctx := gofr.NewContext(nil, nil, app)
		r, err := s.InsertService(ctx, tc.input)

		if !reflect.DeepEqual(err, tc.experr) {
			t.Errorf("%s : expected %v, got %v ", tc.desc, tc.experr, err)
		}

		if tc.experr == nil && !reflect.DeepEqual(r, tc.input) {
			t.Errorf("%s : expected %v, but got %v", tc.desc, tc.input, r)
		}
	}
}
