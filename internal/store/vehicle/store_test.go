package vehicle

import (
	"context"
	"fmt"
	"log"
	"reflect"
	"testing"
	"time"

	"developer.zopsmart.com/go/gofr/pkg/datastore"
	"developer.zopsmart.com/go/gofr/pkg/errors"
	"developer.zopsmart.com/go/gofr/pkg/gofr"
	"github.com/DATA-DOG/go-sqlmock"
	"github.com/SukantArora/CRUD_Gofr/internal/models"
	"gopkg.in/guregu/null.v4"
)

func NewMock() (*gofr.Context, sqlmock.Sqlmock) {
	db, mock, err := sqlmock.New(sqlmock.QueryMatcherOption(sqlmock.QueryMatcherEqual))
	if err != nil {
		log.Fatalf("an error '%s' was not expected when opening a stub database connection", err)
	}

	ctx := gofr.NewContext(nil, nil, &gofr.Gofr{DataStore: datastore.DataStore{ORM: db}})
	ctx.Context = context.Background()

	return ctx, mock
}

func TestGetByID(t *testing.T) {
	ctx, mock := NewMock()
	handler := New()
	query := "select id,model,color,numberPlate,updatedAt,createdAt,name,launched from Vehicle where id = ? and deletedAt is NULL"
	tcs := []struct {
		Desc        string
		ID          int
		Model       string
		Color       string
		NumberPlate string
		Name        string
		Launched    bool
		mockQ       interface{}
		expectedErr error
	}{
		{
			Desc:        "SUCCESS",
			ID:          1,
			Model:       "Q5",
			Color:       "Black",
			NumberPlate: "UK 07 1896",
			Name:        "Audi",
			Launched:    true,
			mockQ: mock.ExpectQuery(query).
				WithArgs(1).
				WillReturnRows(mock.NewRows([]string{"id", "model", "color", "numberPlate", "updatedAt", "createdAt", "name", "launched"}).
					AddRow(1, "Q5", "Black", "UK 07 1896", time.Date(2022, 1, 1, 1, 1, 1, 1, time.UTC),
						time.Date(2022, 1, 1, 1, 1, 1, 1, time.UTC), "Audi", true)),
			expectedErr: nil,
		},
		{
			Desc:        "FAILURE",
			ID:          1,
			Model:       "Q5",
			Color:       "Black",
			NumberPlate: "UK 07 1896",
			Name:        "Audi",
			Launched:    true,
			mockQ: mock.ExpectQuery(query).
				WithArgs(1).
				WillReturnError(errors.Error("internal server error")),
			expectedErr: errors.Error("internal server error"),
		},
	}

	for _, testCase := range tcs {
		testCase := testCase
		t.Run(testCase.Desc, func(t *testing.T) {
			respObj, err := handler.GetByID(ctx, testCase.ID)
			if respObj == nil && err != nil && err.Error() != testCase.expectedErr.Error() {
				t.Errorf("Got : %v, Want : %v", err.Error(), testCase.expectedErr.Error())
			}
		})
	}
}

func TestCreateVehicle(t *testing.T) {
	ctx, mock := NewMock()

	query := "INSERT INTO Vehicle (model,color,numberPlate,name,launched) values (?,?,?,?,?);"

	tcs := []struct {
		Desc        string
		ID          int
		Model       string
		Color       string
		NumberPlate string
		Name        string
		Launched    null.Bool
		mockQ       interface{}
		expectedErr error
	}{
		{
			Desc:        "SUCCESS",
			ID:          1,
			Model:       "Q5",
			Color:       "Black",
			NumberPlate: "UK 07 1896",
			Name:        "Audi",
			Launched:    null.BoolFrom(true),
			mockQ: []interface{}{
				mock.ExpectExec(query).
					WithArgs("Q5", "Black", "UK 07 1896", "Audi", true).
					WillReturnResult(sqlmock.NewResult(1, 1)),
				mock.ExpectQuery(squery).
					WithArgs(1).
					WillReturnRows(mock.NewRows([]string{"id", "model", "color", "numberPlate", "updatedAt", "createdAt", "name", "launched"}).
						AddRow(1, "Q5", "Black", "UK 07 1896", time.Date(2022, 1, 1, 1, 1, 1, 1, time.UTC),
							time.Date(2022, 1, 1, 1, 1, 1, 1, time.UTC), "Audi", true)),
			},
			expectedErr: nil,
		},

		// Failure : prepare command fail
		{
			Desc:        "FAILURE",
			ID:          1,
			Model:       "Q5",
			Color:       "Black",
			NumberPlate: "UK 07 1896",
			Name:        "Audi",
			Launched:    null.BoolFrom(true),
			mockQ:       mock.ExpectExec(query).WillReturnError(errors.Error("internal server error")),
			expectedErr: errors.Error("internal server error"),
		},

		{
			Desc:        "GetByID failure",
			ID:          1,
			Model:       "Q5",
			Color:       "Black",
			NumberPlate: "UK 07 1896",
			Name:        "Audi",
			Launched:    null.BoolFrom(true),
			mockQ: []interface{}{
				mock.ExpectExec(query).WithArgs("Q5", "Black", "UK 07 1896", "Audi", true).WillReturnResult(sqlmock.NewResult(1, 1)),
				mock.ExpectQuery(squery).WithArgs(1).WillReturnError(errors.Error("internal server error")),
			},
			expectedErr: errors.Error("internal server error"),
		},
	}
	handler := New()

	for _, testCase := range tcs {
		testCase := testCase
		t.Run(testCase.Desc, func(t *testing.T) {
			testObj := &models.Vehicle{
				ID:          testCase.ID,
				Model:       testCase.Model,
				Color:       testCase.Color,
				NumberPlate: testCase.NumberPlate,
				Name:        testCase.Name,
				Launched:    null.BoolFrom(true),
			}
			respObj, err := handler.Create(ctx, testObj)
			fmt.Println(respObj, err)
			if respObj == nil && err != nil && err.Error() != testCase.expectedErr.Error() {
				t.Errorf("Got : %v, Want : %v", err.Error(), testCase.expectedErr.Error())
			}
		})
	}
}

func TestUpdateByID(t *testing.T) {
	ctx, mock := NewMock()
	handler := New()
	query := "update Vehicle set model=?,color=?,numberPlate=?,name=?,launched=? where id=? and deletedAt is null"
	squery := "select id,model,color,numberPlate,updatedAt,createdAt,name,launched from Vehicle where id = ? and deletedAt is NULL"

	tcs := []struct {
		Desc        string
		ID          int
		Model       string
		Color       string
		NumberPlate string
		Name        string
		Launched    null.Bool
		mockQ       interface{}
		expectedErr error
	}{
		{
			Desc:        "GET BY ID FAILURE",
			ID:          1,
			Model:       "Q5",
			Color:       "Black",
			NumberPlate: "UK 07 1896",
			Name:        "Audi",
			Launched:    null.BoolFrom(true),
			mockQ: []interface{}{
				mock.
					ExpectExec("update Vehicle set model=?,color=?,numberPlate=?,name=?,launched=? where id=? and deletedAt is null").
					WithArgs("Q5", "Black", "UK 07 1896", "Audi", true, 1).WillReturnResult(sqlmock.NewResult(1, 1)),
				mock.ExpectQuery(squery).WithArgs(1).WillReturnError(errors.Error("internal server error")),
			},
			expectedErr: errors.Error("internal server error"),
		},
		{
			Desc:        "SUCCESS",
			ID:          1,
			Model:       "Q5",
			Color:       "Black",
			NumberPlate: "UK 07 1896",
			Name:        "Audi",
			Launched:    null.BoolFrom(true),
			mockQ: []interface{}{mock.ExpectExec(query).
				WithArgs("Q5", "Black", "UK 07 1896", "Audi", true, 1).
				WillReturnResult(sqlmock.NewResult(1, 1)),
				mock.ExpectQuery(squery).
					WithArgs(1).
					WillReturnRows(mock.NewRows([]string{"id", "model", "color", "numberPlate", "updatedAt", "createdAt", "name", "launched"}).
						AddRow(1, "Q5", "Black", "UK 07 1896", time.Date(2022, 1, 1, 1, 1, 1, 1, time.UTC),
							time.Date(2022, 1, 1, 1, 1, 1, 1, time.UTC), "Audi", true)),
			},
			expectedErr: nil,
		},

		{
			Desc:        "FAILURE",
			ID:          -1,
			Model:       "Q5",
			Color:       "Black",
			NumberPlate: "UK 07 1896",
			Name:        "Audi",
			Launched:    null.BoolFrom(true),
			mockQ: mock.ExpectExec(query).
				WithArgs("Q5", "Black", "UK 07 1896", "Audi", true, -1).
				WillReturnError(errors.Error("internal server error")),
			expectedErr: errors.Error("internal server error"),
		},
		{
			Desc:        "Nothing to update",
			ID:          1,
			Model:       "",
			Color:       "",
			NumberPlate: "",
			Name:        "",
			Launched:    null.NewBool(false, false),
			expectedErr: errors.Error("nothing to update"),
		},
	}

	for _, testCase := range tcs {
		testCase := testCase
		t.Run(testCase.Desc, func(t *testing.T) {
			vehicle := &models.Vehicle{Model: testCase.Model, Color: testCase.Color, NumberPlate: testCase.NumberPlate,
				Name: testCase.Name, Launched: testCase.Launched}
			res, err := handler.Update(ctx, testCase.ID, vehicle)
			if res == nil && !reflect.DeepEqual(err, testCase.expectedErr) {
				t.Errorf("expected error:%v, got:%v", testCase.expectedErr, err)
			}
		})
	}
}

func TestDeleteByID(t *testing.T) {
	ctx, mock := NewMock()
	handler := New()
	query := "update Vehicle set deletedAt = ? where id=?"

	tcs := []struct {
		ID          int
		Desc        string
		expectedErr error
		mockQ       interface{}
	}{
		{
			ID:          1,
			Desc:        "SUCCESS",
			expectedErr: nil,
			mockQ: mock.ExpectPrepare(query).ExpectExec().
				WithArgs(sqlmock.AnyArg(), 1).
				WillReturnResult(sqlmock.NewResult(0, 1)),
		},

		{
			ID:          1,
			Desc:        "FAILURE",
			expectedErr: errors.Error("internal server error"),
			mockQ: mock.ExpectPrepare(query).ExpectExec().
				WithArgs(sqlmock.AnyArg(), 1).
				WillReturnError(errors.Error("internal server error")),
		},

		{
			ID:          1,
			Desc:        "FAILURE",
			expectedErr: errors.Error("internal server error"),
			mockQ: mock.ExpectPrepare(query).
				WillReturnError(errors.Error("internal server error")),
		},
		{
			ID:          1,
			Desc:        "FAILURE",
			expectedErr: errors.EntityNotFound{Entity: "Movie", ID: "1"},
			mockQ: mock.ExpectPrepare(query).ExpectExec().
				WithArgs(sqlmock.AnyArg(), 1).
				WillReturnResult(sqlmock.NewResult(0, 0)),
		},
	}

	for _, testCase := range tcs {
		testCase := testCase
		t.Run(testCase.Desc, func(t *testing.T) {
			err := handler.Delete(ctx, testCase.ID)
			if !reflect.DeepEqual(err, testCase.expectedErr) {
				t.Errorf("expected error:%v, got:%v", testCase.expectedErr, err)
			}
		})
	}
}
func TestGetAll(t *testing.T) {
	ctx, mock := NewMock()
	handler := New()
	query := "Select id,model,color,numberPlate,updatedAt,createdAt,name,launched from Vehicle where deletedAt is NULL;"
	tcs := []struct {
		Desc        string
		expectedErr error
		mockQ       interface{}
	}{
		{
			Desc:        "SUCCESS",
			expectedErr: nil,
			mockQ: mock.ExpectQuery(query).
				WillReturnRows(sqlmock.
					NewRows([]string{"id", "model", "color", "numberPlate", "updatedAt", "createdAt", "name", "launched"}).
					AddRow(1, "Q5", "Black", "UK 07 1896", time.Date(2022, 1, 1, 1, 1, 1, 1, time.UTC),
						time.Date(2022, 1, 1, 1, 1, 1, 1, time.UTC), "Audi", true)).
				WillReturnError(nil),
		},
		{
			Desc:        "FAILURE",
			expectedErr: errors.Error("internal server error"),
			mockQ:       mock.ExpectQuery(query).WillReturnError(errors.Error("internal server error")),
		},
		{
			Desc:        "SCAN FAILURE",
			expectedErr: errors.Error("internal server error"),
			mockQ: mock.ExpectQuery(query).
				WillReturnRows(sqlmock.
					NewRows([]string{"id", "model", "color", "numberPlate", "updatedAt", "createdAt", "name"}).
					AddRow(1, "Q5", "Black", "UK 07 1896", time.Date(2022, 1, 1, 1, 1, 1, 1, time.UTC),
						time.Date(2022, 1, 1, 1, 1, 1, 1, time.UTC), "Audi")).
				WillReturnError(errors.Error("internal server error")),
		},
	}

	for _, testCase := range tcs {
		testCase := testCase
		t.Run(testCase.Desc, func(t *testing.T) {
			_, err := handler.Get(ctx)
			if !reflect.DeepEqual(err, testCase.expectedErr) {
				t.Errorf("expected error:%v, got:%v", testCase.expectedErr, err)
			}
		})
	}
}
