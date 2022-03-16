package http

import (
	"bytes"
	"database/sql"
	"net/http/httptest"
	"testing"

	"developer.zopsmart.com/go/gofr/pkg/errors"
	"developer.zopsmart.com/go/gofr/pkg/gofr"
	"developer.zopsmart.com/go/gofr/pkg/gofr/request"
	"developer.zopsmart.com/go/gofr/pkg/gofr/responder"
	"github.com/SukantArora/CRUD_Gofr/internal/models"
	"github.com/SukantArora/CRUD_Gofr/internal/service"
	"github.com/golang/mock/gomock"
	"github.com/stretchr/testify/assert"
)

var baseURL = "/vehicles"

func TestGetByID(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()
	vehServ := service.NewMockVehicle(ctrl)
	h := New(vehServ)

	testCases := []struct {
		ID          string
		mock        *gomock.Call
		expectedErr error
	}{
		{
			ID:          "1",
			mock:        vehServ.EXPECT().GetByID(gomock.Any(), gomock.Any()).Return(&models.Vehicle{}, nil),
			expectedErr: nil,
		},
		{
			ID:          "-5",
			mock:        vehServ.EXPECT().GetByID(gomock.Any(), gomock.Any()).Return(&models.Vehicle{}, errors.InvalidParam{Param: []string{"id"}}),
			expectedErr: errors.InvalidParam{Param: []string{"id"}},
		},
		{
			ID:          "-5",
			mock:        vehServ.EXPECT().GetByID(gomock.Any(), gomock.Any()).Return(&models.Vehicle{}, sql.ErrNoRows),
			expectedErr: sql.ErrNoRows,
		},
		{
			ID:          "aaa",
			expectedErr: errors.InvalidParam{Param: []string{"id"}},
		},
	}

	for _, testCase := range testCases {
		link := baseURL
		req := httptest.NewRequest("GET", link, nil)
		ctx := gofr.NewContext(responder.NewContextualResponder(httptest.NewRecorder(), req), request.NewHTTPRequest(req), nil)
		ctx.SetPathParams(map[string]string{"id": testCase.ID})
		_, err := h.GetByID(ctx)

		assert.Equal(t, testCase.expectedErr, err)
	}
}

func TestCreate(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()
	vehServ := service.NewMockVehicle(ctrl)
	h := New(vehServ)

	tcs := []struct {
		body   []byte
		mock   *gomock.Call
		stCode int
		err    error
	}{

		{
			body: []byte(`{
				"model": "i9",
				"color": "Black",
				"numberPlate": "MH 04 AT 890",
				"name": "BMW",
				"launched": true
			}`),
			stCode: 200,
			mock:   vehServ.EXPECT().Create(gomock.Any(), gomock.Any()).Return(&models.Vehicle{}, nil),
			err:    nil,
		},

		{
			// Failure Case Empty Model
			body: []byte(`{
				"model": "",
				"color": "Black",
				"numberPlate": "MH 04 AT 890",
				"name": "BMW",
				"launched": true
			}`),
			stCode: 400,
			mock:   vehServ.EXPECT().Create(gomock.Any(), gomock.Any()).Return(&models.Vehicle{}, errors.InvalidParam{Param: []string{"Body"}}),
			err:    errors.InvalidParam{Param: []string{"Body"}},
		},

		{
			// Failure Case Empty Name
			body: []byte(`{
				"model": "i9",
				"color": "Black",
				"numberPlate": "MH 04 AT 890",
				"name": "",
				"launched": true
			}`),
			// service layer
			stCode: 400,
			mock:   vehServ.EXPECT().Create(gomock.Any(), gomock.Any()).Return(&models.Vehicle{}, errors.InvalidParam{Param: []string{"Body"}}),
			err:    errors.InvalidParam{Param: []string{"Body"}},
		},

		{
			body: []byte(`{
				"model": "i9",
				"color": "Black",
				"numberPlate": "",
				"name": "BMW",
				"launched": true
			}`),
			stCode: 400,
			mock:   vehServ.EXPECT().Create(gomock.Any(), gomock.Any()).Return(&models.Vehicle{}, errors.InvalidParam{Param: []string{"Body"}}),
			err:    errors.InvalidParam{Param: []string{"Body"}},
		},

		{
			// Failure Case parsing body
			body: []byte(`
				model: "i9",
				"color": "Black",
				"numberPlate": "MH 04 AT 890",
				"name": "",
				"launched": true
			}`),
			// service layer
			stCode: 400,
			//mock:   vehServ.EXPECT().Create(gomock.Any(), gomock.Any()).Return(&models.Vehicle{}, errors.New("parsing body error")),
			err: errors.InvalidParam{Param: []string{"Body"}},
		},
	}

	for _, tc := range tcs {
		req := httptest.NewRequest("POST", "/employee", bytes.NewReader(tc.body))
		ctx := gofr.NewContext(responder.NewContextualResponder(httptest.NewRecorder(), req), request.NewHTTPRequest(req), nil)
		_, err := h.Create(ctx)
		assert.Equal(t, tc.err, err)
	}
}

func TestGet(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()
	vehServ := service.NewMockVehicle(ctrl)
	h := New(vehServ)

	testCases := []struct {
		mock        *gomock.Call
		statusCode  int
		expectedErr error
	}{
		{

			mock:        vehServ.EXPECT().Get(gomock.Any()).Return([]*models.Vehicle{}, nil),
			statusCode:  200,
			expectedErr: nil,
		},
		{
			mock:        vehServ.EXPECT().Get(gomock.Any()).Return(nil, sql.ErrNoRows),
			statusCode:  404,
			expectedErr: sql.ErrNoRows,
		},
		{
			mock:        vehServ.EXPECT().Get(gomock.Any()).Return(nil, errors.Error("internal server error")),
			statusCode:  500,
			expectedErr: errors.Error("internal server error"),
		},
	}

	for _, testCase := range testCases {
		req := httptest.NewRequest("POST", baseURL, nil)
		ctx := gofr.NewContext(responder.NewContextualResponder(httptest.NewRecorder(), req), request.NewHTTPRequest(req), nil)
		_, err := h.Get(ctx)
		assert.Equal(t, testCase.expectedErr, err)
	}
}

func TestDelete(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()
	vehServ := service.NewMockVehicle(ctrl)
	h := New(vehServ)

	testCases := []struct {
		ID          string
		mock        *gomock.Call
		statusCode  int
		expectedErr error
	}{
		{
			ID:          "200",
			mock:        vehServ.EXPECT().Delete(gomock.Any(), gomock.Any()).Return(sql.ErrNoRows),
			statusCode:  404,
			expectedErr: sql.ErrNoRows,
		},
		{
			ID:          "-1",
			mock:        vehServ.EXPECT().Delete(gomock.Any(), -1).Return(errors.InvalidParam{Param: []string{"id"}}),
			statusCode:  400,
			expectedErr: errors.InvalidParam{Param: []string{"id"}},
		},
		{
			ID:          "200",
			mock:        vehServ.EXPECT().Delete(gomock.Any(), gomock.Any()).Return(errors.Error("internal server error")),
			statusCode:  500,
			expectedErr: errors.Error("internal server error"),
		},
		{
			ID:          "1",
			mock:        vehServ.EXPECT().Delete(gomock.Any(), gomock.Any()).Return(nil),
			statusCode:  200,
			expectedErr: nil,
		},
		{
			ID:          "aaa",
			statusCode:  400,
			expectedErr: errors.InvalidParam{Param: []string{"id"}},
		},
	}

	for _, testCase := range testCases {
		link := baseURL
		req := httptest.NewRequest("DELETE", link, nil)
		ctx := gofr.NewContext(responder.NewContextualResponder(httptest.NewRecorder(), req), request.NewHTTPRequest(req), nil)
		ctx.SetPathParams(map[string]string{"id": testCase.ID})
		_, err := h.Delete(ctx)
		assert.Equal(t, testCase.expectedErr, err)
	}
}

func TestUpdate(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()
	vehServ := service.NewMockVehicle(ctrl)
	h := New(vehServ)

	testCases := []struct {
		body          []byte
		ID            string
		mock          *gomock.Call
		statusCode    int
		expectedError error
	}{
		// Success Case
		{
			body: []byte(`{
				"id": 12,
				"model": "i9",
				"color": "Black",
				"numberPlate": "MH 04 AT 890",
				"name": "BMW",
				"launched": true
		    }`),
			ID:            "1",
			mock:          vehServ.EXPECT().Update(gomock.Any(), gomock.Any(), gomock.Any()).Return(&models.Vehicle{}, nil),
			statusCode:    200,
			expectedError: nil,
		},
		{
			body: []byte(`{
				"id": 12,
				"model": "i9",
				"color": "Black",
				"numberPlate": "MH 04 AT 890",
				"name": "BMW",
				"launched": true
		    }`),
			ID:            "1",
			mock:          vehServ.EXPECT().Update(gomock.Any(), gomock.Any(), gomock.Any()).Return(nil, sql.ErrNoRows),
			statusCode:    404,
			expectedError: sql.ErrNoRows,
		},
		{
			body: []byte(`{
				"id": 12,
				"model": "i9",
				"color": "Black",
				"numberPlate": "MH 04 AT 890",
				"name": "BMW",
				"launched": true
		    }`),
			ID:            "1",
			mock:          vehServ.EXPECT().Update(gomock.Any(), gomock.Any(), gomock.Any()).Return(nil, sql.ErrConnDone),
			statusCode:    500,
			expectedError: sql.ErrConnDone,
		},
		{
			// Failure Case parsing body
			ID: "1",
			body: []byte(`{
				id: 12,
				model: "i9",
				"color": "Black",
				"numberPlate": "MH 04 AT 890",
				"name": "",
				"launched": true
			}`),
			// service layer
			statusCode: 400,
			//mock:          vehServ.EXPECT().Update(gomock.Any(), gomock.Any(), gomock.Any()).Return(nil, errors.New("parsing body error")),
			expectedError: errors.InvalidParam{Param: []string{"Body"}},
		},
		{
			ID:            "aaa",
			statusCode:    400,
			expectedError: errors.InvalidParam{Param: []string{"id"}},
		},
	}

	for _, testCase := range testCases {
		link := "/vehicles"
		req := httptest.NewRequest("PUT", link, bytes.NewReader(testCase.body))
		ctx := gofr.NewContext(responder.NewContextualResponder(httptest.NewRecorder(), req), request.NewHTTPRequest(req), nil)
		ctx.SetPathParams(map[string]string{"id": testCase.ID})

		_, err := h.Update(ctx)

		assert.Equal(t, testCase.expectedError, err)
	}
}
