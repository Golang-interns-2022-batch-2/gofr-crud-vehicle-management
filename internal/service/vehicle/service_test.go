package vehicle

import (
	"database/sql"
	"reflect"
	"testing"

	"developer.zopsmart.com/go/gofr/pkg/errors"
	"developer.zopsmart.com/go/gofr/pkg/gofr"
	"github.com/SukantArora/CRUD_Gofr/internal/models"
	"github.com/SukantArora/CRUD_Gofr/internal/store"
	"github.com/golang/mock/gomock"
	"gopkg.in/guregu/null.v4"
)

func TestGetByID(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()
	vehStore := store.NewMockVehicle(ctrl)
	handler := New(vehStore)

	testCases := []struct {
		ID            int
		expectedError error
		mock          *gomock.Call
	}{
		{
			ID:            -1,
			expectedError: errors.InvalidParam{Param: []string{"id"}},
		},
		{
			ID:            1,
			expectedError: nil,
			mock:          vehStore.EXPECT().GetByID(gomock.Any(), gomock.Any()).Return(nil, nil),
		},
	}

	for _, testCase := range testCases {
		ctx := gofr.NewContext(nil, nil, gofr.New())
		_, err := handler.GetByID(ctx, testCase.ID)

		if !reflect.DeepEqual(err, testCase.expectedError) {
			t.Errorf("expected error:%v, got:%v", testCase.expectedError, err)
		}
	}
}

func TestInsert(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	vehstore := store.NewMockVehicle(ctrl)
	handler := New(vehstore)

	testCases := []struct {
		vehicle       *models.Vehicle
		expectedError error
		mock          *gomock.Call
	}{

		{
			vehicle: &models.Vehicle{
				Model:       "Q5",
				Color:       "Black",
				NumberPlate: "UK 07 1896",
				Name:        "Audi",
				Launched:    null.BoolFrom(true),
			},
			expectedError: nil,
			mock:          vehstore.EXPECT().Create(gomock.Any(), gomock.Any()).Return(nil, nil),
		},

		{
			vehicle: &models.Vehicle{
				Model:       "Q5",
				Color:       "",
				NumberPlate: "UK 07 1896",
				Name:        "Audi",
				Launched:    null.BoolFrom(true),
			},
			expectedError: errors.InvalidParam{Param: []string{"Color"}},
		},
		{
			vehicle: &models.Vehicle{
				Model:       "",
				Color:       "Black",
				NumberPlate: "UK 07 1896",
				Name:        "Audi",
				Launched:    null.BoolFrom(true),
			},
			expectedError: errors.InvalidParam{Param: []string{"Model"}},
		},
		{
			vehicle: &models.Vehicle{
				Model:       "Q5",
				Color:       "Black",
				NumberPlate: "UK 07 1896",
				Name:        "",
				Launched:    null.BoolFrom(true),
			},
			expectedError: errors.InvalidParam{Param: []string{"Name"}},
		},
		{
			vehicle: &models.Vehicle{
				Model:       "Q5",
				Color:       "Black",
				NumberPlate: "",
				Name:        "Audi",
				Launched:    null.BoolFrom(true),
			},
			expectedError: errors.InvalidParam{Param: []string{"NumberPlate"}},
		},
	}

	for _, testCase := range testCases {
		ctx := gofr.NewContext(nil, nil, gofr.New())
		_, err := handler.Create(ctx, testCase.vehicle)

		if !reflect.DeepEqual(err, testCase.expectedError) {
			t.Errorf("expected error:%v, got:%v", testCase.expectedError, err)
		}
	}
}

func TestUpdate(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()
	vehStore := store.NewMockVehicle(ctrl)
	handler := New(vehStore)

	testCases := []struct {
		vehicle       *models.Vehicle
		expectedError error
		mock          *gomock.Call
	}{

		{
			vehicle: &models.Vehicle{
				ID:          1,
				Model:       "Q5",
				Color:       "Black",
				NumberPlate: "UK 07 1896",
				Name:        "Audi",
				Launched:    null.BoolFrom(true),
			},
			expectedError: nil,
			mock:          vehStore.EXPECT().Update(gomock.Any(), gomock.Any(), gomock.Any()).Return(nil, nil),
		},
		{
			vehicle: &models.Vehicle{
				ID:          -1,
				Model:       "Q5",
				Color:       "Black",
				NumberPlate: "UK 07 1896",
				Name:        "Audi",
				Launched:    null.BoolFrom(true),
			},
			expectedError: errors.InvalidParam{Param: []string{"id"}},
		},
	}

	for _, testCase := range testCases {
		ctx := gofr.NewContext(nil, nil, gofr.New())
		_, err := handler.Update(ctx, testCase.vehicle.ID, testCase.vehicle)

		if !reflect.DeepEqual(err, testCase.expectedError) {
			t.Errorf("expected error:%v, got:%v", testCase.expectedError, err)
		}
	}
}

func TestGetAll(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()
	vehStore := store.NewMockVehicle(ctrl)
	handler := New(vehStore)

	testCases := []struct {
		expectedError error
		mock          *gomock.Call
	}{

		{
			expectedError: nil,
			mock:          vehStore.EXPECT().Get(gomock.Any()).Return([]*models.Vehicle{}, nil),
		},
		{
			expectedError: sql.ErrNoRows,
			mock:          vehStore.EXPECT().Get(gomock.Any()).Return(nil, sql.ErrNoRows),
		},
	}

	for _, testCase := range testCases {
		ctx := gofr.NewContext(nil, nil, gofr.New())
		_, err := handler.Get(ctx)

		if !reflect.DeepEqual(err, testCase.expectedError) {
			t.Errorf("expected error:%v, got:%v", testCase.expectedError, err)
		}
	}
}

func TestDelete(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()
	vehStore := store.NewMockVehicle(ctrl)
	handler := New(vehStore)
	testCases := []struct {
		ID            int
		expectedError error
		mock          *gomock.Call
	}{
		{
			ID:            0,
			expectedError: errors.InvalidParam{Param: []string{"id"}},
		},
		{
			ID:            1,
			expectedError: nil,
			mock:          vehStore.EXPECT().Delete(gomock.Any(), gomock.Any()).Return(nil),
		},
	}

	for _, testCase := range testCases {
		ctx := gofr.NewContext(nil, nil, gofr.New())
		err := handler.Delete(ctx, testCase.ID)

		if !reflect.DeepEqual(err, testCase.expectedError) {
			t.Errorf("expected error:%v, got:%v", testCase.expectedError, err)
		}
	}
}
