package vehicle

import (
	"bytes"
	"database/sql"
	"net/http"
	"net/http/httptest"
	"reflect"
	"testing"

	"developer.zopsmart.com/go/gofr/pkg/errors"
	"developer.zopsmart.com/go/gofr/pkg/gofr"
	"developer.zopsmart.com/go/gofr/pkg/gofr/request"
	"developer.zopsmart.com/go/gofr/pkg/gofr/responder"
	"github.com/SN786/gofr_vms/internal/model"
	"github.com/SN786/gofr_vms/internal/service"
	"github.com/bmizerany/assert"
	"github.com/golang/mock/gomock"
	"gopkg.in/guregu/null.v4"
)

func TestGet(t *testing.T) {
	ctrl := gomock.NewController(t)

	defer ctrl.Finish()

	VehSvc := service.NewMockVehicleManager(ctrl)

	h := New(VehSvc)

	app := gofr.New()

	tests := []struct {
		id   string
		desc string
		resp interface{}
		mock *gomock.Call
		err  error
	}{
		{
			id:   "1",
			desc: "get by id succuss",
			resp: model.Vehicle{
				ID:          1,
				Model:       "some name",
				Color:       "Red",
				NumberPlate: "xy-123",
				UpdatedAt:   gomock.Any().String(),
				CreatedAt:   gomock.Any().String(),
				Name:        "Yamaha",
				Launched:    null.BoolFrom(true),
			},
			mock: VehSvc.EXPECT().Get(gomock.Any(), 1).Return(&model.Vehicle{}, nil),
			err:  nil,
		},
		{
			id:   "1",
			desc: "get by id fail",
			resp: model.Vehicle{
				ID:          1,
				Model:       "some name",
				Color:       "Red",
				NumberPlate: "xy-123",
				UpdatedAt:   gomock.Any().String(),
				CreatedAt:   gomock.Any().String(),
				Name:        "Yamaha",
				Launched:    null.BoolFrom(true),
			},
			mock: VehSvc.EXPECT().Get(gomock.Any(), 1).Return(nil, sql.ErrNoRows),
			err:  sql.ErrNoRows,
		},
	}

	for _, tc := range tests {
		w := httptest.NewRecorder()
		r := httptest.NewRequest(http.MethodGet, "http://localhost:8090", nil)
		req := request.NewHTTPRequest(r)
		res := responder.NewContextualResponder(w, r)
		ctx := gofr.NewContext(res, req, app)

		ctx.SetPathParams(map[string]string{
			"id": tc.id,
		})

		_, err := h.Get(ctx)
		assert.Equal(t, tc.err, err)
	}
}

func TestCreate(t *testing.T) {
	ctrl := gomock.NewController(t)

	defer ctrl.Finish()

	VehSvc := service.NewMockVehicleManager(ctrl)
	h := Handler{VehSvc}
	app := gofr.New()

	tcs := []struct {
		desc          string
		body          []byte
		mock          *gomock.Call
		statusCode    int
		expectedError error
	}{

		{
			desc: "Success Case",
			body: []byte(`{
				"model": 		"why12345",
				"color": 		"Red",
				"numberPlate": 	"MH 03 AT 007",
				"name": 		"HondaCity",
				"launched": 	false
		    }`),
			mock: VehSvc.EXPECT().Post(
				gomock.Any(),
				&model.Vehicle{
					Model:       "why12345",
					Color:       "Red",
					NumberPlate: "MH 03 AT 007",
					Name:        "HondaCity",
					Launched:    null.BoolFrom(false),
				}).Return(
				&model.Vehicle{
					Model:       "why12345",
					Color:       "Red",
					NumberPlate: "MH 03 AT 007",
					Name:        "HondaCity",
					Launched:    null.BoolFrom(false),
				}, nil),

			statusCode:    200,
			expectedError: nil,
		},

		{
			desc: "Faliure Case: Sql No Row",
			body: []byte(`{
				"model": 		"why12345",
				"color": 		"Red",
				"numberPlate": 	"MH 03 AT 007",
				"name": 		"HondaCity",
				"launched": 	false
		    }`),
			mock: VehSvc.EXPECT().Post(
				gomock.Any(),
				&model.Vehicle{
					Model:       "why12345",
					Color:       "Red",
					NumberPlate: "MH 03 AT 007",
					Name:        "HondaCity",
					Launched:    null.BoolFrom(false),
				}).Return(nil, sql.ErrNoRows),
			statusCode:    404,
			expectedError: sql.ErrNoRows,
		},

		{
			desc: "Faliure Case: Unmarshall",
			body: []byte(`{
				"model": 		1,
				"color": 		"Red",
				"numberPlate": 	"MH 03 AT 007",
				"name": 		"susjh",
				"launched": 	false
		    }`),
			statusCode:    400,
			expectedError: errors.InvalidParam{Param: []string{"body"}},
		},
	}

	for _, tc := range tcs {
		r := httptest.NewRequest(http.MethodPost, "http://localhost:8090", bytes.NewReader(tc.body))
		w := httptest.NewRecorder()
		req := request.NewHTTPRequest(r)
		res := responder.NewContextualResponder(w, r)
		ctx := gofr.NewContext(res, req, app)

		_, err := h.Post(ctx)

		if !reflect.DeepEqual(err, tc.expectedError) {
			t.Errorf("expected error:%v, got:%v", tc.expectedError, err)
		}
	}
}
func TestDelete(t *testing.T) {
	ctrl := gomock.NewController(t)

	defer ctrl.Finish()

	app := gofr.New()
	VehSvc := service.NewMockVehicleManager(ctrl)
	h := Handler{VehSvc}

	tcs := []struct {
		ID            string
		mock          *gomock.Call
		expectedError error
		statusCode    int
	}{
		{
			ID:            "0",
			mock:          VehSvc.EXPECT().Delete(gomock.Any(), 0).Return(nil),
			expectedError: nil,
			statusCode:    200,
		},
		{
			ID:            "1",
			mock:          VehSvc.EXPECT().Delete(gomock.Any(), 1).Return(errors.Error("no data found")),
			expectedError: errors.Error("no data found"),
			statusCode:    404,
		},
	}

	for _, tc := range tcs {
		w := httptest.NewRecorder()
		r := httptest.NewRequest(http.MethodDelete, "http://localhost:8090", nil)
		req := request.NewHTTPRequest(r)
		res := responder.NewContextualResponder(w, r)
		ctx := gofr.NewContext(res, req, app)

		ctx.SetPathParams(map[string]string{
			"id": tc.ID,
		})

		_, err := h.Delete(ctx)

		if !reflect.DeepEqual(err, tc.expectedError) {
			t.Errorf("expected error:%v, got:%v", tc.expectedError, err)
		}
	}
}

func TestUpdate(t *testing.T) {
	ctrl := gomock.NewController(t)

	defer ctrl.Finish()

	VehSvc := service.NewMockVehicleManager(ctrl)
	app := gofr.New()

	tcs := []struct {
		body          []byte
		ID            string
		mock          *gomock.Call
		statusCode    int
		expectedError error
	}{
		{
			body: []byte(`{
				"model": 		"why12345",
				"name": 		"HondaCity",
				"launched": 	false
		    }`),
			ID: "1",
			mock: VehSvc.EXPECT().Update(
				gomock.Any(),
				1,
				&model.Vehicle{
					Model:    "why12345",
					Name:     "HondaCity",
					Launched: null.BoolFrom(false),
				}).Return(&model.Vehicle{
				Model:    "why12345",
				Name:     "HondaCity",
				Launched: null.BoolFrom(false)}, nil),
			statusCode:    200,
			expectedError: nil,
		},
		{
			body: []byte(`{
				"model": 		"why12345",
				"name": 		"HondaCity",
				"launched": 	false
		    }`),
			ID: "-1",
			mock: VehSvc.EXPECT().Update(
				gomock.Any(),
				-1,
				&model.Vehicle{
					Model:    "why12345",
					Name:     "HondaCity",
					Launched: null.BoolFrom(false),
				}).Return(nil, errors.Error("parsing body error")),
			statusCode:    500,
			expectedError: errors.Error("parsing body error"),
		},
		{
			body: []byte(`{
				"model": 		"why12345",
				"name": 		"HondaCity",
				"launched": 	false
		    }`),
			ID: "1",
			mock: VehSvc.EXPECT().Update(
				gomock.Any(),
				1,
				&model.Vehicle{
					Model:    "why12345",
					Name:     "HondaCity",
					Launched: null.BoolFrom(false),
				}).Return(nil, sql.ErrNoRows),
			statusCode:    404,
			expectedError: sql.ErrNoRows,
		},
		{
			body: []byte(`{
				"model": 		1,
				"name": 		"HondaCity",
				"launched": 	false
		    }`),
			ID:            "1",
			statusCode:    400,
			expectedError: errors.InvalidParam{Param: []string{"body"}},
		},
	}

	h := New(VehSvc)

	for _, tc := range tcs {
		r := httptest.NewRequest(http.MethodPost, "http://localhost:8090", bytes.NewReader(tc.body))
		w := httptest.NewRecorder()
		req := request.NewHTTPRequest(r)
		res := responder.NewContextualResponder(w, r)
		ctx := gofr.NewContext(res, req, app)

		ctx.SetPathParams(map[string]string{
			"id": tc.ID,
		})

		_, err := h.Update(ctx)

		if !reflect.DeepEqual(err, tc.expectedError) {
			t.Errorf("expected error:%v, got:%v", tc.expectedError, err)
		}
	}
}

func TestGetAll(t *testing.T) {
	ctrl := gomock.NewController(t)

	defer ctrl.Finish()

	VehSvc := service.NewMockVehicleManager(ctrl)
	h := Handler{VehSvc}
	app := gofr.New()

	tcs := []struct {
		mock          *gomock.Call
		expectedError error
	}{
		{

			mock:          VehSvc.EXPECT().GetAll(gomock.Any()).Return([]*model.Vehicle{}, nil),
			expectedError: nil,
		},
		{
			mock:          VehSvc.EXPECT().GetAll(gomock.Any()).Return(nil, sql.ErrNoRows),
			expectedError: sql.ErrNoRows,
		},
	}

	for _, tc := range tcs {
		r := httptest.NewRequest(http.MethodPost, "http://localhost:8090", nil)
		w := httptest.NewRecorder()
		req := request.NewHTTPRequest(r)
		res := responder.NewContextualResponder(w, r)
		ctx := gofr.NewContext(res, req, app)

		_, err := h.GetAll(ctx)

		if !reflect.DeepEqual(err, tc.expectedError) {
			t.Errorf("expected error:%v, got:%v", tc.expectedError, err)
		}
	}
}
