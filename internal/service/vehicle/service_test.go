package vehicle

import (
	"database/sql"
	"reflect"
	"testing"

	"developer.zopsmart.com/go/gofr/pkg/errors"
	"gopkg.in/guregu/null.v4"

	"developer.zopsmart.com/go/gofr/pkg/gofr"
	"github.com/SN786/gofr_vms/internal/model"
	"github.com/SN786/gofr_vms/internal/store"
	"github.com/golang/mock/gomock"
)

func TestGet(t *testing.T) {
	ctrl := gomock.NewController(t)

	defer ctrl.Finish()

	StoreSvc := store.NewMockVehicleManager(ctrl)
	h := New(StoreSvc)

	tcs := []struct {
		ID int
		// resp        interface{}
		expectError error
		mock        *gomock.Call
	}{

		{
			ID:          -1,
			expectError: errors.InvalidParam{Param: []string{"id"}},
		},
		{
			ID:          2,
			expectError: nil,
			mock:        StoreSvc.EXPECT().GetDetailsByID(gomock.Any(), 2).Return(&model.Vehicle{}, nil),
		},
		{
			ID:          3,
			expectError: errors.Error("no data found"),
			mock:        StoreSvc.EXPECT().GetDetailsByID(gomock.Any(), 3).Return(nil, errors.Error("no data found")),
		},
	}

	for _, tc := range tcs {
		ctx := gofr.NewContext(nil, nil, gofr.New())
		_, err := h.Get(ctx, tc.ID)

		if !reflect.DeepEqual(err, tc.expectError) {
			t.Errorf("expected error:%v, got:%v", tc.expectError, err)
		}
	}
}

func TestPost(t *testing.T) {
	ctrl := gomock.NewController(t)

	defer ctrl.Finish()

	StoreSvc := store.NewMockVehicleManager(ctrl)

	tcs := []struct {
		vehicle     model.Vehicle
		expectError error
		mock        []*gomock.Call
	}{
		{
			vehicle: model.Vehicle{
				ID:          6,
				Model:       "Iter",
				Color:       "Black",
				NumberPlate: "MH 03 AT 007",
				Name:        "BMW",
				Launched:    null.BoolFrom(true),
			},
			expectError: nil,
			mock: []*gomock.Call{
				StoreSvc.EXPECT().InsertVehicle(
					gomock.Any(),
					&model.Vehicle{
						ID:          6,
						Model:       "Iter",
						Color:       "Black",
						NumberPlate: "MH 03 AT 007",
						Name:        "BMW",
						Launched:    null.BoolFrom(true),
					}).Return(&model.Vehicle{
					ID:          6,
					Model:       "Iter",
					Color:       "Black",
					NumberPlate: "MH 03 AT 007",
					Name:        "BMW",
					Launched:    null.BoolFrom(true),
				}, nil),

				StoreSvc.EXPECT().GetDetailsByID(gomock.Any(), 6).Return(&model.Vehicle{
					ID:          6,
					Model:       "Iter",
					Color:       "Black",
					NumberPlate: "MH 03 AT 007",
					Name:        "BMW",
					Launched:    null.BoolFrom(true),
				}, nil),
			},
		},

		{
			vehicle: model.Vehicle{
				ID:          6,
				Model:       "Iter",
				Color:       "Black",
				NumberPlate: "MH 03 AT 007",
				Name:        "BMW",
				Launched:    null.BoolFrom(true),
			},
			expectError: sql.ErrNoRows,
			mock: []*gomock.Call{
				StoreSvc.EXPECT().
					InsertVehicle(gomock.Any(),
						&model.Vehicle{
							ID:          6,
							Model:       "Iter",
							Color:       "Black",
							NumberPlate: "MH 03 AT 007",
							Name:        "BMW",
							Launched:    null.BoolFrom(true),
						}).
					Return(
						&model.Vehicle{
							ID:          6,
							Model:       "Iter",
							Color:       "Black",
							NumberPlate: "MH 03 AT 007",
							Name:        "BMW",
							Launched:    null.BoolFrom(true),
						}, nil),
				StoreSvc.EXPECT().GetDetailsByID(gomock.Any(), 6).Return(nil, sql.ErrNoRows),
			},
		},

		{
			vehicle: model.Vehicle{
				ID:          6,
				Model:       "Iter",
				Color:       "Black",
				NumberPlate: "MH 03 AT 007",
				Name:        "BMW",
				Launched:    null.BoolFrom(true),
			},
			expectError: errors.Error("no data found"),
			mock: []*gomock.Call{
				StoreSvc.EXPECT().InsertVehicle(gomock.Any(),
					&model.Vehicle{
						ID:          6,
						Model:       "Iter",
						Color:       "Black",
						NumberPlate: "MH 03 AT 007",
						Name:        "BMW",
						Launched:    null.BoolFrom(true),
					}).
					Return(nil, errors.Error("no data found")),
			},
		},
	}
	h := New(StoreSvc)

	for _, tc := range tcs {
		ctx := gofr.NewContext(nil, nil, gofr.New())
		_, err := h.Post(ctx, &tc.vehicle)

		if !reflect.DeepEqual(err, tc.expectError) {
			t.Errorf("expected error:%v, got:%v", tc.expectError, err)
		}
	}
}

func TestDelete(t *testing.T) {
	ctrl := gomock.NewController(t)

	defer ctrl.Finish()

	StoreSvc := store.NewMockVehicleManager(ctrl)

	tcs := []struct {
		ID          int
		expectError error
		mock        *gomock.Call
	}{

		{
			ID:          -1,
			expectError: errors.InvalidParam{Param: []string{"id"}},
		},
		{
			ID:          2,
			expectError: nil,
			mock:        StoreSvc.EXPECT().DeleteVehicleByID(gomock.Any(), 2).Return(nil),
		},
		{
			ID:          2,
			expectError: errors.Error("no data found"),
			mock:        StoreSvc.EXPECT().DeleteVehicleByID(gomock.Any(), 2).Return(errors.Error("no data found")),
		},
	}
	h := New(StoreSvc)

	for _, tc := range tcs {
		ctx := gofr.NewContext(nil, nil, gofr.New())
		err := h.Delete(ctx, tc.ID)

		if !reflect.DeepEqual(err, tc.expectError) {
			t.Errorf("expected error:%v, got:%v", tc.expectError, err)
		}
	}
}

func TestUpdate(t *testing.T) {
	ctrl := gomock.NewController(t)

	defer ctrl.Finish()

	StoreSvc := store.NewMockVehicleManager(ctrl)

	tcs := []struct {
		vehicle     model.Vehicle
		expectError error
		mock        []*gomock.Call
	}{
		{
			vehicle: model.Vehicle{
				ID:          -8,
				Model:       "hd",
				Color:       "Black",
				NumberPlate: "MH 03 AT 007",
				Name:        "",
				Launched:    null.BoolFrom(true),
			},
			expectError: errors.InvalidParam{Param: []string{"id"}},
		},

		{
			vehicle: model.Vehicle{
				ID:          7,
				Model:       "hd",
				Color:       "Black",
				NumberPlate: "MH 03 AT 007",
				Name:        "yamaha",
				Launched:    null.BoolFrom(true),
			},
			expectError: errors.Error("no data found"),
			mock: []*gomock.Call{
				StoreSvc.EXPECT().UpdateVehicleByID(
					gomock.Any(),
					&model.Vehicle{
						ID:          7,
						Model:       "hd",
						Color:       "Black",
						NumberPlate: "MH 03 AT 007",
						Name:        "yamaha",
						Launched:    null.BoolFrom(true)},
				).Return(errors.Error("no data found")),
			},
		},

		{
			vehicle: model.Vehicle{
				ID:          7,
				Model:       "hd",
				Color:       "Black",
				NumberPlate: "MH 03 AT 007",
				Name:        "yamaha",
				Launched:    null.BoolFrom(true),
			},
			expectError: nil,
			mock: []*gomock.Call{
				StoreSvc.EXPECT().UpdateVehicleByID(
					gomock.Any(),
					&model.Vehicle{
						ID:          7,
						Model:       "hd",
						Color:       "Black",
						NumberPlate: "MH 03 AT 007",
						Name:        "yamaha",
						Launched:    null.BoolFrom(true)},
				).Return(nil),

				StoreSvc.EXPECT().GetDetailsByID(gomock.Any(), 7).Return(&model.Vehicle{
					ID:          7,
					Model:       "hd",
					Color:       "Black",
					NumberPlate: "MH 03 AT 007",
					Name:        "yamaha",
					Launched:    null.BoolFrom(true)}, nil),
			},
		},
	}

	h := New(StoreSvc)

	for _, tc := range tcs {
		ctx := gofr.NewContext(nil, nil, gofr.New())
		_, err := h.Update(ctx, int(tc.vehicle.ID), &tc.vehicle)

		if !reflect.DeepEqual(err, tc.expectError) {
			t.Errorf("expected error:%v, got:%v", tc.expectError, err)
		}
	}
}

func TestGetAll(t *testing.T) {
	ctrl := gomock.NewController(t)

	defer ctrl.Finish()

	StoreSvc := store.NewMockVehicleManager(ctrl)

	tcs := []struct {
		expectError error
		mock        *gomock.Call
	}{
		{
			expectError: sql.ErrNoRows,
			mock:        StoreSvc.EXPECT().GetAll(gomock.Any()).Return(nil, sql.ErrNoRows),
		},
		{
			expectError: nil,
			mock:        StoreSvc.EXPECT().GetAll(gomock.Any()).Return([]*model.Vehicle{}, nil),
		},
	}
	h := New(StoreSvc)

	for _, tc := range tcs {
		ctx := gofr.NewContext(nil, nil, gofr.New())
		_, err := h.GetAll(ctx)

		if !reflect.DeepEqual(err, tc.expectError) {
			t.Errorf("expected error:%v, got:%v", tc.expectError, err)
		}
	}
}
