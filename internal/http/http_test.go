package http

import (
	"bytes"
	"database/sql"
	"errors"
	"net/http/httptest"
	"reflect"
	"testing"

	gerror "developer.zopsmart.com/go/gofr/pkg/errors"
	"developer.zopsmart.com/go/gofr/pkg/gofr"
	"developer.zopsmart.com/go/gofr/pkg/gofr/request"
	"developer.zopsmart.com/go/gofr/pkg/gofr/responder"
	"github.com/golang/mock/gomock"
	"github.com/ishankochar09/go_pro/gofrTutorial/internal/models"
	"github.com/ishankochar09/go_pro/gofrTutorial/internal/service"
)

func TestGetVehicle(t *testing.T) {
	mockCtrl := gomock.NewController(t)
	defer mockCtrl.Finish()
	MockVehicleInterface := service.NewMockVehicleInterface(mockCtrl)
	handler := New(MockVehicleInterface)

	testCases := []struct {
		ID            string
		mock          *gomock.Call
		scode         int
		expectedError error
	}{
		{
			ID: "2",
			mock: MockVehicleInterface.EXPECT().
				GetIDVehicle(gomock.Any(), gomock.Any()).
				Return(models.Vehicle{}, nil),
			scode:         200,
			expectedError: nil,
		},
		{
			ID: "-7",
			mock: MockVehicleInterface.EXPECT().
				GetIDVehicle(gomock.Any(), gomock.Any()).
				Return(models.Vehicle{}, errors.New("error Invalid Id")),
			scode:         400,
			expectedError: errors.New("error Invalid Id"),
		},
		{
			ID:            "abcc",
			scode:         400,
			expectedError: gerror.MissingParam{Param: []string{"id"}},
		},
	}
	for _, tc := range testCases {
		link := "/vehicle" + "/" + tc.ID
		req := httptest.NewRequest("GET", link, nil)
		ctx := gofr.NewContext(responder.NewContextualResponder(httptest.NewRecorder(), req), request.NewHTTPRequest(req), nil)
		ctx.SetPathParams(map[string]string{"id": tc.ID})
		_, err := handler.GetByIDVehicle(ctx)

		if !reflect.DeepEqual(tc.expectedError, err) {
			t.Errorf("Expected error %v, got:%v", tc.expectedError, err)
		}
	}
}
func TestCreateVehicle(t *testing.T) {
	mockCtrl := gomock.NewController(t)
	defer mockCtrl.Finish()
	MockVehicleInterface := service.NewMockVehicleInterface(mockCtrl)
	handler := New(MockVehicleInterface)

	testcases := []struct {
		reqBody     []byte
		resBody     []byte
		mock        *gomock.Call
		scode       int
		ExpectedErr error
	}{
		{
			reqBody: []byte(`{"Id":5,"Model":"i8","Color":"black", "NumberPlate":"MH 03 AT 007","Name":"BMW", "Launched": true}`),
			resBody: []byte(
				`{"Code":200,
				"Status":"SUCCESS",
				"Data":
				{"Vehicle":
				{"id":5,
				"model":"i8",
				"color":"black",
				"numberPlate":"MH 03 AT 007",
				"updatedAt":"0001-01-01T00:00:00Z",
				"createdAt":"0001-01-01T00:00:00Z",
				"name":"BMW",
				"launched":true
				}
				}
				}`),
			mock: MockVehicleInterface.
				EXPECT().
				Create(gomock.Any(), gomock.Any()).
				Return(
					models.Vehicle{ID: 5, Model: "i8", Color: "black", NumberPlate: "MH 03 AT 007", Name: "BMW", Launched: true}, nil),
			scode:       200,
			ExpectedErr: nil,
		},
		{
			reqBody: []byte(`{"Id":5,"Model":"i8","Color":"black", "NumberPlate":"MH 03 AT 007","Name":"BMW", "Launched": true}`),
			resBody: []byte(
				`{"Code":200,
				 "Status":"SUCCESS",
				 "Data":{
					 "Vehicle":
					 {"id":5,
					 "model":"i8",
					 "color":"black",
					 "numberPlate":"MH 03 AT 007",
					 "updatedAt":"0001-01-01T00:00:00Z",
					 "createdAt":"0001-01-01T00:00:00Z",
					 "name":"BMW",
					 "launched":true
					 }
					 }
					 }`),
			mock: MockVehicleInterface.
				EXPECT().
				Create(gomock.Any(), gomock.Any()).
				Return(models.Vehicle{ID: 5, Model: "i8", Color: "black", NumberPlate: "MH 03 AT 007", Name: "BMW", Launched: true},
					gerror.InvalidParam{Param: []string{"body"}}),
			scode:       400,
			ExpectedErr: gerror.InvalidParam{Param: []string{"body"}},
		},
	}

	for _, tc := range testcases {
		r := httptest.NewRequest("POST", "/vehicle", bytes.NewReader(tc.reqBody))
		w := httptest.NewRecorder()

		req := request.NewHTTPRequest(r)
		res := responder.NewContextualResponder(w, r)
		ctx := gofr.NewContext(res, req, nil)
		_, err := handler.CreateVehicle(ctx)

		if !reflect.DeepEqual(err, tc.ExpectedErr) {
			t.Errorf(" expected %v, but got %v", tc.ExpectedErr, err)
		}
	}
}
func TestGetAllVehicles(t *testing.T) {
	mockCtrl := gomock.NewController(t)
	defer mockCtrl.Finish()
	MockVehicleInterface := service.NewMockVehicleInterface(mockCtrl)
	handler := New(MockVehicleInterface)

	testcases := []struct {
		mock        *gomock.Call
		scode       int
		expectedErr error
	}{

		{

			mock:        MockVehicleInterface.EXPECT().AllVehicles(gomock.Any()).Return([]models.Vehicle{}, nil),
			scode:       200,
			expectedErr: nil,
		},
		{
			mock:        MockVehicleInterface.EXPECT().AllVehicles(gomock.Any()).Return(nil, sql.ErrNoRows),
			scode:       400,
			expectedErr: sql.ErrNoRows,
		},
	}

	for _, tc := range testcases {
		r := httptest.NewRequest("GET", "/vehicle", nil)
		w := httptest.NewRecorder()
		req := request.NewHTTPRequest(r)
		res := responder.NewContextualResponder(w, r)
		ctx := gofr.NewContext(res, req, nil)
		_, err := handler.GetAllVehicles(ctx)

		if !reflect.DeepEqual(err, tc.expectedErr) {
			t.Errorf(" expected %v, but got %v", tc.expectedErr, err)
		}
	}
}
func TestDeleteVeh(t *testing.T) {
	mockCtrl := gomock.NewController(t)
	defer mockCtrl.Finish()
	MockVehicleInterface := service.NewMockVehicleInterface(mockCtrl)
	handler := New(MockVehicleInterface)

	testCases := []struct {
		ID          string
		mock        *gomock.Call
		scode       int
		expectedErr error
	}{
		{
			ID:          "-5",
			mock:        MockVehicleInterface.EXPECT().DeleteIDVehicle(gomock.Any(), gomock.Any()).Return(sql.ErrNoRows),
			scode:       400,
			expectedErr: sql.ErrNoRows,
		},
		{
			ID:          "1",
			mock:        MockVehicleInterface.EXPECT().DeleteIDVehicle(gomock.Any(), gomock.Any()).Return(nil),
			scode:       200,
			expectedErr: nil,
		},
		{
			ID:          "abc",
			scode:       400,
			expectedErr: gerror.MissingParam{Param: []string{"id"}},
		},
	}

	for _, tc := range testCases {
		link := "/vehicle" + tc.ID
		req := httptest.NewRequest("DELETE", link, nil)
		ctx := gofr.NewContext(responder.NewContextualResponder(httptest.NewRecorder(), req), request.NewHTTPRequest(req), nil)
		ctx.SetPathParams(map[string]string{"id": tc.ID})
		_, err := handler.DeleteVehicle(ctx)

		if !reflect.DeepEqual(err, tc.expectedErr) {
			t.Errorf(" expected %v, but got %v", tc.expectedErr, err)
		}
	}
}
func TestVehicleUpdate(t *testing.T) {
	mockCtrl := gomock.NewController(t)
	defer mockCtrl.Finish()
	MockVehicleInterface := service.NewMockVehicleInterface(mockCtrl)
	handler := New(MockVehicleInterface)
	testCases := []struct {
		ID          string
		scode       int
		mock        *gomock.Call
		expectedErr error
		veh         models.Vehicle
		reqBody     []byte
	}{
		{
			ID:          "1",
			scode:       200,
			mock:        MockVehicleInterface.EXPECT().UpdateIDVehicle(gomock.Any(), gomock.Any(), gomock.Any()).Return(models.Vehicle{}, nil),
			expectedErr: nil,
			veh:         models.Vehicle{ID: 5, Model: "i8", Color: "black", NumberPlate: "MH 03 AT 007", Name: "BMW", Launched: true},
			reqBody:     []byte(`{"Id":5,"Model":"i8","Color":"black", "NumberPlate":"MH 03 AT 007","Name":"BMW", "Launched": true}`),
		},
		{
			ID:    "-2",
			scode: 400,
			mock: MockVehicleInterface.EXPECT().
				UpdateIDVehicle(gomock.Any(), gomock.Any(), gomock.Any()).
				Return(models.Vehicle{}, errors.New("invalid id")),
			expectedErr: errors.New("invalid id"),
			veh:         models.Vehicle{ID: 5, Model: "i8", Color: "black", NumberPlate: "MH 03 AT 007", Name: "BMW", Launched: true},
			reqBody:     []byte(`{"Id":5,"Model":"i8","Color":"black", "NumberPlate":"MH 03 AT 007","Name":"BMW", "Launched": true}`),
		},
		{
			ID:          "abc",
			scode:       400,
			expectedErr: gerror.MissingParam{Param: []string{"id"}},
			veh:         models.Vehicle{ID: 5, Model: "i8", Color: "black", NumberPlate: "MH 03 AT 007", Name: "BMW", Launched: true},
			reqBody:     []byte(`{"Id":5,"Model":"i8","Color":"black", "NumberPlate":"MH 03 AT 007","Name":"BMW", "Launched": true}`),
		},
		{
			ID:    "2",
			scode: 400,
			mock: MockVehicleInterface.EXPECT().
				UpdateIDVehicle(gomock.Any(), gomock.Any(), gomock.Any()).
				Return(models.Vehicle{}, gerror.InvalidParam{Param: []string{"body"}}),
			expectedErr: gerror.InvalidParam{Param: []string{"body"}},
			veh:         models.Vehicle{ID: 5, Model: "i8", Color: "black", NumberPlate: "MH 03 AT 007", Name: "BMW", Launched: true},
			reqBody:     []byte(`{"Id":5,"Model":"","Color":"", "NumberPlate":"","Name":"BMW", "Launched": true}`),
		},
	}

	for _, tc := range testCases {
		link := "/vehicles" + tc.ID
		r := httptest.NewRequest("PUT", link, bytes.NewReader(tc.reqBody))
		w := httptest.NewRecorder()

		req := request.NewHTTPRequest(r)
		res := responder.NewContextualResponder(w, r)
		ctx := gofr.NewContext(res, req, nil)
		ctx.SetPathParams(map[string]string{"id": tc.ID})
		_, err := handler.UpdateVehicle(ctx)

		if !reflect.DeepEqual(err, tc.expectedErr) {
			t.Errorf(" expected %v, but got %v", tc.expectedErr, err)
		}
	}
}
