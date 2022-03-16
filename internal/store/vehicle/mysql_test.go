package vehicle

import (
	"context"
	"database/sql"
	"fmt"
	"time"

	"reflect"
	"testing"

	"developer.zopsmart.com/go/gofr/pkg/datastore"
	"developer.zopsmart.com/go/gofr/pkg/errors"
	"developer.zopsmart.com/go/gofr/pkg/gofr"
	"github.com/DATA-DOG/go-sqlmock"
	"github.com/SN786/gofr_vms/internal/model"
	"gopkg.in/guregu/null.v4"
)

func TestGetDetailsByID(t *testing.T) {
	db, mock, err := sqlmock.New(sqlmock.QueryMatcherOption(sqlmock.QueryMatcherEqual))

	var handler = New()

	ctx := gofr.NewContext(nil, nil, &gofr.Gofr{DataStore: datastore.DataStore{ORM: db}})

	ctx.Context = context.Background()

	if err != nil {
		t.Errorf("An error %s occurred when opening a database connection", err)
		return
	}

	tcs := []struct {
		ID          int
		Model       string
		Color       string
		NumberPlate string
		UpdatedAt   string
		CreatedAt   string
		Name        string
		Launched    bool
		mockQuery   interface{}
		expectError error
	}{
		{
			ID:          1,
			Model:       "i8",
			Color:       "Black",
			NumberPlate: "MH 03 AT 007",
			Name:        "BMW",
			Launched:    true,
			UpdatedAt:   "2021-12-17 13:39:41",
			CreatedAt:   "2021-12-17 13:39:41",
			mockQuery: mock.
				ExpectQuery("select id,model,color,numberPlate,updatedAt,createdAt,name,launched from vehicles where id=? and deletedAt is NULL").
				WithArgs(1).WillReturnRows(mock.NewRows([]string{"id", "model", "color", "numberPlate", "updatedAt", "createdAt", "name", "launched"}).
				AddRow(1, "i8", "Black", "MH 03 AT 007", "2021-12-17 13:39:41", "2021-12-17 13:39:41", "BMW", true)),
			expectError: nil,
		},

		{
			ID:          -1,
			Model:       "i10",
			Color:       "White",
			NumberPlate: "MH 03 AT 123",
			Name:        "Suzuki",
			Launched:    false,
			mockQuery: mock.
				ExpectQuery("select id,model,color,numberPlate,updatedAt,createdAt,name,launched from vehicles where id=? and deletedAt is NULL").
				WithArgs(-1).WillReturnError(errors.Error(fmt.Sprintf("couldn't get the vechicle: %v", -1))),
			expectError: errors.Error(fmt.Sprintf("couldn't get the vechicle: %v", -1)),
		},

		{
			ID:          23,
			Model:       "i10",
			Color:       "White",
			NumberPlate: "MH 03 AT 123",
			Name:        "Suzuki",
			Launched:    false,
			mockQuery: mock.
				ExpectQuery("select id,model,color,numberPlate,updatedAt,createdAt,name,launched from vehicles where id=? and deletedAt is NULL").
				WithArgs(23).WillReturnError(sql.ErrNoRows),
			expectError: sql.ErrNoRows,
		},
	}

	for _, testCase := range tcs {
		testCase := testCase

		t.Run("", func(t *testing.T) {
			res, err := handler.GetDetailsByID(ctx, testCase.ID)

			if res != nil && err != nil && err.Error() != testCase.expectError.Error() {
				t.Errorf("expected error:%v, got:%v", testCase.expectError, err)
			}
		})
	}
}

func TestInsertVehicle(t *testing.T) {
	db, mock, err := sqlmock.New(sqlmock.QueryMatcherOption(sqlmock.QueryMatcherEqual))

	var handler = New()

	ctx := gofr.NewContext(nil, nil, &gofr.Gofr{DataStore: datastore.DataStore{ORM: db}})
	ctx.Context = context.Background()

	if err != nil {
		t.Errorf("An error %s occurred when opening a database connection", err)
		return
	}

	tcs := []struct {
		ID          int64
		Model       string
		Color       string
		NumberPlate string
		UpdatedAt   string
		CreatedAt   string
		Name        string
		Launched    null.Bool
		mockQuery   interface{}
		expectError error
	}{
		{
			Model:       "i8",
			Color:       "Black",
			NumberPlate: "MH 03 AT 007",
			Name:        "BMW",
			Launched:    null.BoolFrom(true),
			mockQuery: mock.ExpectExec("insert into vehicles (model,color,numberPlate,name,launched) values(?,?,?,?,?)").
				WithArgs("i8", "Black", "MH 03 AT 007", "BMW", true).WillReturnResult(sqlmock.NewResult(1, 1)),
			expectError: nil,
		},

		{
			Model:       "i10",
			Color:       "White",
			NumberPlate: "MH 03 AT 123",
			Name:        "Suzuki",
			Launched:    null.BoolFrom(false),
			mockQuery: mock.ExpectExec("insert into vehicles (model,color,numberPlate,name,launched) values(?,?,?,?,?)").
				WithArgs("i10", "White", "MH 03 AT 123", "Suzuki", false).WillReturnError(errors.Error("couldn't insert the vechicle")),
			expectError: errors.Error("couldn't insert the vechicle"),
		},
	}

	for _, testCase := range tcs {
		testCase := testCase

		t.Run("", func(t *testing.T) {
			vehicle := model.Vehicle{
				ID:          testCase.ID,
				Model:       testCase.Model,
				Color:       testCase.Color,
				NumberPlate: testCase.NumberPlate,
				CreatedAt:   testCase.CreatedAt,
				UpdatedAt:   testCase.UpdatedAt,
				Name:        testCase.Name,
				Launched:    testCase.Launched,
			}
			res, err := handler.InsertVehicle(ctx, &vehicle)
			if err != nil && err.Error() != testCase.expectError.Error() && res == nil {
				t.Errorf("expected error:%v, got:%v", testCase.expectError, err)
			}
		})
	}
}
func TestDeleteVehicleByID(t *testing.T) {
	db, mock, err := sqlmock.New(sqlmock.QueryMatcherOption(sqlmock.QueryMatcherEqual))

	var handler = New()

	ctx := gofr.NewContext(nil, nil, &gofr.Gofr{DataStore: datastore.DataStore{ORM: db}})
	ctx.Context = context.Background()

	if err != nil {
		t.Errorf("An error %s occurred when opening a database connection", err)
		return
	}

	tcs := []struct {
		ID          int
		deleateTime string
		mockQuery   interface{}
		expectError error
	}{
		{
			ID: 1,
			mockQuery: mock.ExpectExec("update vehicles set deletedAt=? where id=? and deletedAt is null").WithArgs(sqlmock.AnyArg(), 1).
				WillReturnError(errors.Error(fmt.Sprintf("couldn't delete the vechicle: %v", 1))),
			expectError: errors.Error(fmt.Sprintf("couldn't delete the vechicle: %v", 1)),
		},
		{
			ID: 1,
			mockQuery: mock.ExpectExec("update vehicles set deletedAt=? where id=? and deletedAt is null").WithArgs(sqlmock.AnyArg(), 1).
				WillReturnError(sql.ErrNoRows),
			expectError: sql.ErrNoRows,
		},
		{
			ID: 2,
			mockQuery: mock.ExpectExec("update vehicles set deletedAt=? where id=? and deletedAt is null").WithArgs(sqlmock.AnyArg(), 2).
				WillReturnResult(sqlmock.NewResult(1, 1)),
			expectError: nil,
		},
		{
			ID: 2,
			mockQuery: mock.ExpectExec("update vehicles set deletedAt=? where id=? and deletedAt is null").WithArgs(sqlmock.AnyArg(), 2).
				WillReturnResult(sqlmock.NewResult(1, 0)),
			expectError: errors.EntityNotFound{
				Entity: "vehicle",
				ID:     "2",
			},
		},
	}

	for _, testCase := range tcs {
		testCase := testCase

		t.Run("", func(t *testing.T) {
			err := handler.DeleteVehicleByID(ctx, testCase.ID)
			if !reflect.DeepEqual(err, testCase.expectError) {
				t.Errorf("expected error:%v, got:%v", testCase.expectError, err)
			}
		})
	}
}

func TestUpdateVehicleById(t *testing.T) {
	db, mock, err := sqlmock.New(sqlmock.QueryMatcherOption(sqlmock.QueryMatcherEqual))

	var handler = New()

	ctx := gofr.NewContext(nil, nil, &gofr.Gofr{DataStore: datastore.DataStore{ORM: db}})
	ctx.Context = context.Background()

	if err != nil {
		t.Errorf("An error %s occurred when opening a database connection", err)
		return
	}

	tcs := []struct {
		ID          int64
		Model       string
		Color       string
		NumberPlate string
		UpdatedAt   string
		CreatedAt   string
		Name        string
		Launched    null.Bool
		mockQuery   interface{}
		expectError error
	}{
		{
			ID:        -1,
			Model:     "i8",
			Name:      "BMW",
			Launched:  null.BoolFrom(true),
			UpdatedAt: time.Now().String()[:19],
			mockQuery: mock.ExpectExec("update vehicles set model=?,name=?,launched=?,updatedAt=? where id=? and deletedAt is NULL").
				WithArgs("i8", "BMW", true, sqlmock.AnyArg(), -1).WillReturnError(errors.Error(fmt.Sprintf("couldn't update the vechicle: %v", 1))),
			expectError: errors.Error(fmt.Sprintf("couldn't update the vechicle: %v", -1)),
		},
		{
			ID:        -1,
			Model:     "i8",
			Name:      "BMW",
			Launched:  null.BoolFrom(true),
			UpdatedAt: time.Now().String()[:19],
			mockQuery: mock.ExpectExec("update vehicles set model=?,name=?,launched=?,updatedAt=? where id=? and deletedAt is NULL").
				WithArgs("i8", "BMW", true, sqlmock.AnyArg(), -1).WillReturnError(sql.ErrNoRows),
			expectError: sql.ErrNoRows,
		},
		{
			ID:        1,
			Model:     "i8",
			Name:      "BMW",
			Launched:  null.BoolFrom(true),
			UpdatedAt: time.Now().String()[:19],
			mockQuery: mock.ExpectExec("update vehicles set model=?,name=?,launched=?,updatedAt=? where id=? and deletedAt is NULL").
				WithArgs("i8", "BMW", true, sqlmock.AnyArg(), 1).WillReturnResult(sqlmock.NewResult(0, 1)),
			expectError: nil,
		},
	}

	for _, testCase := range tcs {
		testCase := testCase

		t.Run("", func(t *testing.T) {
			vehicle := model.Vehicle{
				ID:       testCase.ID,
				Model:    testCase.Model,
				Name:     testCase.Name,
				Launched: testCase.Launched,
			}
			err := handler.UpdateVehicleByID(ctx, &vehicle)
			if !reflect.DeepEqual(err, testCase.expectError) {
				t.Errorf("expected error:%v, got:%v", testCase.expectError, err)
			}
		})
	}
}

func TestGetAll(t *testing.T) {
	db, mock, err := sqlmock.New(sqlmock.QueryMatcherOption(sqlmock.QueryMatcherEqual))

	var handler = New()

	ctx := gofr.NewContext(nil, nil, &gofr.Gofr{DataStore: datastore.DataStore{ORM: db}})
	ctx.Context = context.Background()

	if err != nil {
		t.Errorf("An error %s occurred when opening a database connection", err)
		return
	}

	tcs := []struct {
		ID          int64
		Model       string
		Color       string
		NumberPlate string
		UpdatedAt   string
		CreatedAt   string
		Name        string
		Launched    bool
		mockQuery   interface{}
		expectError error
	}{

		{
			ID:          1,
			Model:       "i8",
			Color:       "Black",
			NumberPlate: "MH 03 AT 007",
			Name:        "BMW",
			Launched:    true,
			UpdatedAt:   "2021-12-17 13:39:41",
			CreatedAt:   "2021-12-17 13:39:41",
			mockQuery: mock.ExpectQuery("select id,model,color,numberPlate,updatedAt,createdAt,name,launched from vehicles where deletedAt is NULL").
				WillReturnRows(sqlmock.NewRows([]string{"id", "model", "color", "numberPlate", "updatedAt", "createdAt", "name", "launched"}).
					AddRow(1, "i8", "Black", "MH 03 AT 007", "2021-12-17 13:39:41", "2021-12-17 13:39:41", "BMW", true)).WillReturnError(nil),
			expectError: nil,
		},

		{
			ID:          1,
			Model:       "i8",
			Color:       "Black",
			NumberPlate: "MH 03 AT 007",
			Name:        "BMW",
			Launched:    true,
			UpdatedAt:   "2021-12-17 13:39:41",
			CreatedAt:   "2021-12-17 13:39:41",
			mockQuery: mock.ExpectQuery("select id,model,color,numberPlate,updatedAt,createdAt,name,launched from vehicles where deletedAt is NULL").
				WillReturnRows(sqlmock.NewRows([]string{"id", "model", "color", "numberPlate", "updatedAt", "createdAt", "name", "launched"}).
					AddRow("str", "i8", "Black", "MH 03 AT 007", "2021-12-17 13:39:41", "2021-12-17 13:39:41", "BMW", true)).RowsWillBeClosed().
				WillReturnError(errors.Error("couldn't get list of vechicles:")),
			expectError: errors.Error("couldn't get list of vechicles:"),
		},
		{
			ID:          1,
			Model:       "i8",
			Color:       "Black",
			NumberPlate: "MH 03 AT 007",
			Name:        "BMW",
			Launched:    true,
			UpdatedAt:   "2021-12-17 13:39:41",
			CreatedAt:   "2021-12-17 13:39:41",
			mockQuery: mock.ExpectQuery("select id,model,color,numberPlate,updatedAt,createdAt,name,launched from vehicles where deletedAt is NULL").
				WillReturnRows(sqlmock.NewRows([]string{"id", "model", "color", "numberPlate", "updatedAt", "createdAt", "name", "launched"}).
					AddRow("str", "i8", "Black", "MH 03 AT 007", "2021-12-17 13:39:41", "2021-12-17 13:39:41", "BMW", true)).RowsWillBeClosed().
				WillReturnError(sql.ErrNoRows),
			expectError: sql.ErrNoRows,
		},
	}
	for _, testCase := range tcs {
		testCase := testCase

		t.Run("", func(t *testing.T) {
			_, err := handler.GetAll(ctx)
			if !reflect.DeepEqual(err, testCase.expectError) {
				t.Errorf("expected error:%v, got:%v", testCase.expectError, err)
			}
		})
	}
}
