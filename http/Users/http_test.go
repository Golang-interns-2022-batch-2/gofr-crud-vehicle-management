package vehicle

import (
	"bytes"
	"net/http"
	"net/http/httptest"
	"reflect"
	"testing"
	"time"

	"github.com/Gofr-VMS/model"
	"github.com/Gofr-VMS/service"

	"developer.zopsmart.com/go/gofr/pkg/errors"
	"developer.zopsmart.com/go/gofr/pkg/gofr"
	"developer.zopsmart.com/go/gofr/pkg/gofr/request"
	"developer.zopsmart.com/go/gofr/pkg/gofr/responder"
	"developer.zopsmart.com/go/gofr/pkg/gofr/types"
	"github.com/golang/mock/gomock"
)

func setMock(t *testing.T) (*gofr.Gofr, Handler, *service.MockService) {
	app := gofr.New()
	ctrl := gomock.NewController(t)
	mock := service.NewMockService(ctrl)
	h := New(mock)

	return app, h, mock
}
func setMockHTTP(app *gofr.Gofr, method string, body []byte) *gofr.Context {
	w := httptest.NewRecorder()
	r := httptest.NewRequest(method, "http://dummy", bytes.NewReader(body))

	req := request.NewHTTPRequest(r)
	res := responder.NewContextualResponder(w, r)

	ctx := gofr.NewContext(res, req, app)

	return ctx
}
func TestGet(t *testing.T) {
	app, h, mock := setMock(t)

	TC := []struct {
		desc   string
		input  string
		output types.Response
		mock   []*gomock.Call
		experr error
	}{
		{
			desc:  "success",
			input: "1",
			output:  types.Response{
				Data: model.Vehicles{
					ID:        1,
					Model:     "i8",
					Color:     "Black",
					NumPlate:  "MH 03 AT 007",
					CreatedAt: time.Date(2022, 1, 1, 1, 1, 1, 1, time.UTC),
					UpdatedAt: time.Date(2022, 1, 1, 1, 1, 1, 1, time.UTC),
					Name:      "BMW",
					Launched:  true,
				},
				
			},
			mock: []*gomock.Call{
				mock.EXPECT().GetByIdService(gomock.Any(), 1).Return(&model.Vehicles{
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
			experr: nil,
		},
		{
			desc:   "fail",
			input:  "-1",
			output: types.Response{Data:nil},
			experr: errors.InvalidParam{
				Param: []string{"id"},
			},
			mock: []*gomock.Call{
				mock.EXPECT().GetByIdService(gomock.Any(), -1).Return(nil, errors.InvalidParam{
					Param: []string{"id"},
				}),
			},
		},
	}

	for _, test := range TC {
		tc := test
		t.Run(tc.desc, func(t *testing.T) {
			ctx := setMockHTTP(app, http.MethodGet, nil)
			ctx.SetPathParams(map[string]string{
				"id": tc.input,
			})

			resp, err := h.Get(ctx)
			if !reflect.DeepEqual(err, tc.experr) {
				t.Errorf("%s : expected %v, but got %v", tc.desc, tc.experr, err)
			}

			if tc.experr == nil && !reflect.DeepEqual(resp, tc.output) {
				t.Errorf("%s : expected %v, but got %v", tc.desc, tc.output, resp)
			}
		})
	}
}
func TestRemove(t *testing.T) {
	app, h, mock := setMock(t)
	TC := []struct {
		desc   string
		id     string
		output errors.Response
		experr error
		mock   []*gomock.Call
	}{
		{
			desc: "success",
			id:   "1",
			output: errors.Response{
				StatusCode:   http.StatusOK,
				Code: "success",
				Reason:   "Vehicle deleted successfully.",
			},
			experr: nil,
			mock: []*gomock.Call{
				mock.EXPECT().RemoveService(gomock.Any(), 1).Return(nil),
			},
		},
		{
			desc:   "missing params",
			id:     "",
			experr: errors.MissingParam{Param: []string{"id"}},
			mock:   nil,
		},
	}

	for _, tc := range TC {
		ctx := setMockHTTP(app, http.MethodDelete, nil)
		ctx.SetPathParams(map[string]string{
			"id": tc.id,
		})

		resp, err := h.Remove(ctx)
		if !reflect.DeepEqual(err, tc.experr) {
			t.Errorf("%s : expected %v, but got %v", tc.desc, tc.experr, err)
		}

		if tc.experr == nil && !reflect.DeepEqual(resp, tc.output) {
			t.Errorf("%s : expected %v, but got %v", tc.desc, tc.output, resp)
		}
	}
}
func TestInsert(t *testing.T) {
	app, h, mock := setMock(t)
	TC := []struct {
		desc   string
		input  []byte
		experr error
		mock   []*gomock.Call
	}{
		{
			desc: "success",
			input: []byte(`{
				"id": 12,
				"model": "i9",
				"color": "Black",
				"numplate": "MH 04 AT 890",
				"name": "BMW",
				"launched": true
			}`),
			experr: nil,
			mock: []*gomock.Call{
				mock.EXPECT().InsertService(gomock.Any(), &model.Vehicles{
					ID:       12,
					Model:    "i9",
					Color:    "Black",
					NumPlate: "MH 04 AT 890",
					Name:     "BMW",
					Launched: true,
				}).Return(&model.Vehicles{
					ID:       12,
					Model:    "i9",
					Color:    "Black",
					NumPlate: "MH 04 AT 890",
					Name:     "BMW",
					Launched: true,
				}, nil),
			},
		},
		{
			desc:   "Failure case no body",
			input:  []byte(``),
			experr: errors.InvalidParam{Param: []string{"body"}},
			mock:   nil,
		},
	}

	for _, test := range TC {
		tc := test
		ctx := setMockHTTP(app, http.MethodPost, tc.input)

		_, err := h.Create(ctx)
		if !reflect.DeepEqual(err, tc.experr) {
			t.Errorf("%s : expected %v, but got %v", tc.desc, tc.experr, err)
		}
	}
}
func TestUpadte(t *testing.T) {
	app, h, mock := setMock(t)
	TC := []struct {
		desc   string
		id     string
		input  []byte
		experr error
		mock   []*gomock.Call
	}{
		{
			desc: "string",
			input: []byte(`{
				"id": 1,
				"model": "i9",
				"color": "Black",
				"numberPlate": "MH 04 AT 890",
				"name": "BMW",
				"launched": true
		    }`),
			id:     "1",
			experr: nil,
			mock: []*gomock.Call{
				mock.EXPECT().UpdateService(gomock.Any(), &model.Vehicles{
					ID:       12,
					Model:    "i9",
					Color:    "Black",
					NumPlate: "MH 04 AT 890",
					Name:     "BMW",
					Launched: true,
				}).Return(&model.Vehicles{
					ID:       1,
					Model:    "i9",
					Color:    "Black",
					NumPlate: "MH 04 AT 890",
					Name:     "BMW",
					Launched: true,
				}, nil),
			},
		},
	}

	for _, test := range TC {
		tc := test
		ctx := setMockHTTP(app, http.MethodPost, tc.input)

		_, err := h.Update(ctx)
		if !reflect.DeepEqual(err, tc.experr) {
			t.Errorf("%s : expected %v, but got %v", tc.desc, tc.experr, err)
		}
	}
}
func TestGetAll(t *testing.T) {
	app, h, mock := setMock(t)

	TC := []struct {
		desc     string
		mockCall []*gomock.Call
		output   types.Response
		experr   error
	}{
		{
			desc:   "success",
			experr: nil,
			output: types.Response{
				Data: model.Vehicles{
					ID:        1,
					Model:     "i8",
					Color:     "Black",
					NumPlate:  "MH 03 AT 007",
					CreatedAt: time.Date(2022, 1, 1, 1, 1, 1, 1, time.UTC),
					UpdatedAt: time.Date(2022, 1, 1, 1, 1, 1, 1, time.UTC),
					Name:      "BMW",
					Launched:  true,
				},
			
			},
			mockCall: []*gomock.Call{
				mock.EXPECT().GetAllService(gomock.Any()).Return([]*model.Vehicles{
					{
						ID:        1,
						Model:     "i8",
						Color:     "Black",
						NumPlate:  "MH 03 AT 007",
						CreatedAt: time.Date(2022, 1, 1, 1, 1, 1, 1, time.UTC),
						UpdatedAt: time.Date(2022, 1, 1, 1, 1, 1, 1, time.UTC),
						Name:      "BMW",
						Launched:  true,
					},
				}, nil),
			},
		},
	}
	for _, test := range TC {
		tc := test
		ctx := setMockHTTP(app, http.MethodGet, nil)
		resp, err := h.GetAll(ctx)

		if !reflect.DeepEqual(err, tc.experr) {
			t.Errorf("%s : expected %v, but got %v", tc.desc, tc.experr, err)
		}

		if tc.experr == nil && !reflect.DeepEqual(resp, tc.output) {
			t.Errorf("%s : expected %v, but got %v", tc.desc, tc.output, resp)
		}
	}
}
