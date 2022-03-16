package service

import (
	"database/sql"
	"fmt"
	"reflect"
	"testing"

	"developer.zopsmart.com/go/gofr/pkg/gofr"
	gomock "github.com/golang/mock/gomock"
	"github.com/ishankochar09/go_pro/gofrTutorial/internal/models"
	"github.com/ishankochar09/go_pro/gofrTutorial/internal/store"

	gerror "developer.zopsmart.com/go/gofr/pkg/errors"
)

func TestGetVehicle(t *testing.T) {
	mockCtrl := gomock.NewController(t)
	defer mockCtrl.Finish()
	MockVehicleStore := store.NewMockVehicleStore(mockCtrl)
	VehicleService := New(MockVehicleStore)

	tcs := []struct {
		id          int
		mock        *gomock.Call
		expectedErr error
	}{
		{
			id: 1,
			mock: MockVehicleStore.
				EXPECT().
				GetVehicleByID(gomock.Any(), gomock.Any()).
				Return(models.Vehicle{}, nil),
			expectedErr: nil,
		},
		{
			id:          -2,
			expectedErr: gerror.Error("Validation Error"),
		},
	}

	for _, testCase := range tcs {
		testCase := testCase
		ctx := gofr.NewContext(nil, nil, gofr.New())

		t.Run("testing create service", func(t *testing.T) {
			_, err := VehicleService.GetIDVehicle(ctx, testCase.id)

			if !reflect.DeepEqual(err, testCase.expectedErr) {
				t.Errorf(" Expected Error: %v , Found : %v", testCase.expectedErr, err)
			}
		})
	}
}
func TestVehicleService(t *testing.T) {
	mockCtrl := gomock.NewController(t)
	defer mockCtrl.Finish()
	MockVehicleStore := store.NewMockVehicleStore(mockCtrl)
	VehicleService := New(MockVehicleStore)

	tcs := []struct {
		desc string
		id   int
		veh  models.Vehicle
		out  models.Vehicle
		mock *gomock.Call
		err  error
	}{
		{
			desc: "success case",
			id:   1,
			veh:  models.Vehicle{NumberPlate: "MH 03 AT 007", Name: "BMW", Model: "i8", Color: "Black"},
			out:  models.Vehicle{NumberPlate: "MH 03 AT 007", Name: "BMW", Model: "i8", Color: "Black"},
			mock: MockVehicleStore.
				EXPECT().
				InsertVehicle(gomock.Any(), gomock.Any()).
				Return(models.Vehicle{NumberPlate: "MH 03 AT 007", Name: "BMW", Model: "i8", Color: "Black"}, nil),
			err: nil,
		},
		{
			desc: "error in name",
			id:   1,
			veh:  models.Vehicle{ID: 1, Model: "BMW", Color: "Black", Name: "", NumberPlate: "MH 03 AT 007"},
			out:  models.Vehicle{},
			err:  gerror.Error("Invalid name"),
		},
		{
			desc: "error in ID",
			id:   -2,
			veh:  models.Vehicle{ID: -2, Model: "BMW", Color: "Black", Name: "BMW", NumberPlate: "MH 03 AT 007"},
			out:  models.Vehicle{},
			err:  gerror.Error("invalid Id"),
		},
		{
			desc: "error in model",
			id:   3,
			veh:  models.Vehicle{ID: 3, Model: "", Color: "Black", Name: "BMW", NumberPlate: "MH 03 AT 007"},
			out:  models.Vehicle{},
			err:  gerror.Error("Invalid Model"),
		},
		{
			desc: "error in numberPlate",
			id:   4,
			veh:  models.Vehicle{ID: 4, Model: "i8", Color: "Black", Name: "BMW", NumberPlate: ""},
			out:  models.Vehicle{},
			err:  gerror.Error("invalid numberPlate"),
		},
	}
	for i, testCase := range tcs {
		testCase := testCase
		ctx := gofr.NewContext(nil, nil, gofr.New())

		t.Run("testing create service", func(t *testing.T) {
			out, er := VehicleService.Create(ctx,
				&models.Vehicle{
					ID:          testCase.veh.ID,
					Model:       testCase.veh.Model,
					NumberPlate: testCase.veh.NumberPlate,
					Color:       testCase.veh.Color, Name: testCase.veh.Name})

			if er != nil && (testCase.err.Error() != er.Error()) {
				t.Errorf("desc ---> %v Expected Error: %v , Found : %v", testCase.desc, testCase.err, er)
			}

			if !reflect.DeepEqual(out, testCase.out) {
				t.Errorf(" Expected Error: %v , Found : %v", testCase.out, out)
			}
		})
		fmt.Println("PASSED : ", i)
	}
}
func TestDeleteVehicle(t *testing.T) {
	mockCtrl := gomock.NewController(t)
	defer mockCtrl.Finish()
	MockVehicleStore := store.NewMockVehicleStore(mockCtrl)
	VehicleService := New(MockVehicleStore)
	tcs := []struct {
		id   int
		mock []*gomock.Call
		err  error
	}{
		{
			id: 1,
			mock: []*gomock.Call{
				MockVehicleStore.EXPECT().GetVehicleByID(gomock.Any(), gomock.Any()).Return(models.Vehicle{}, nil),
				MockVehicleStore.EXPECT().DeleteVehicle(gomock.Any(), gomock.Any()).Return(nil)},
			err: nil,
		},
		{
			id:  -2,
			err: gerror.Error("Validation Error"),
		},
	}

	for _, testCase := range tcs {
		testCase := testCase
		ctx := gofr.NewContext(nil, nil, gofr.New())

		t.Run("testing create service", func(t *testing.T) {
			er := VehicleService.DeleteIDVehicle(ctx, testCase.id)

			if !reflect.DeepEqual(er, testCase.err) {
				t.Errorf(" Expected Error: %v , Found : %v", testCase.err, er)
			}
		})
	}
}
func TestUpdateVehicle(t *testing.T) {
	mockCtrl := gomock.NewController(t)
	defer mockCtrl.Finish()
	MockVehicleStore := store.NewMockVehicleStore(mockCtrl)
	VehicleService := New(MockVehicleStore)
	tcs := []struct {
		id      int
		vehicle models.Vehicle
		err     error
		mock    *gomock.Call
	}{
		{
			vehicle: models.Vehicle{ID: -8, Model: "hd", Color: "Black", NumberPlate: "MH 03 AT 007", Name: "", Launched: true},
			err:     gerror.Error("Validation Error"),
		},
		{
			vehicle: models.Vehicle{
				ID:          7,
				Model:       "i9",
				Color:       "Black",
				NumberPlate: "hr 12 9889",
				Name:        "Toyota",
				Launched:    true,
			},
			err: nil,
			mock: MockVehicleStore.
				EXPECT().
				UpdateVehByID(gomock.Any(), gomock.Any(), gomock.Any()).Return(models.Vehicle{}, nil),
		},
	}

	for _, testCase := range tcs {
		testCase := testCase

		t.Run("testing create service", func(t *testing.T) {
			ctx := gofr.NewContext(nil, nil, gofr.New())
			_, er := VehicleService.UpdateIDVehicle(ctx, testCase.id, &testCase.vehicle)
			if !reflect.DeepEqual(er, testCase.err) {
				t.Errorf(" Expected Error: %v , Found : %v", testCase.err, er)
			}
		})
	}
}
func TestGetAll(t *testing.T) {
	mockCtrl := gomock.NewController(t)
	defer mockCtrl.Finish()
	MockVehicleStore := store.NewMockVehicleStore(mockCtrl)

	tcs := []struct {
		expectError error
		mock        *gomock.Call
	}{
		{
			expectError: sql.ErrNoRows,
			mock:        MockVehicleStore.EXPECT().GetAll(gomock.Any()).Return(nil, sql.ErrNoRows),
		},
		{
			expectError: nil,
			mock:        MockVehicleStore.EXPECT().GetAll(gomock.Any()).Return([]models.Vehicle{}, nil),
		},
	}
	h := New(MockVehicleStore)

	for _, tc := range tcs {
		ctx := gofr.NewContext(nil, nil, gofr.New())
		_, err := h.AllVehicles(ctx)

		if !reflect.DeepEqual(err, tc.expectError) {
			t.Errorf("expected error:%v, got:%v", tc.expectError, err)
		}
	}
}
