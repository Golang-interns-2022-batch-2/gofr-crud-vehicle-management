package store

import (
	"context"
	"database/sql"
	"log"
	"reflect"
	"testing"
	"time"

	"developer.zopsmart.com/go/gofr/pkg/datastore"
	"developer.zopsmart.com/go/gofr/pkg/gofr"

	"github.com/DATA-DOG/go-sqlmock"
	"github.com/ishankochar09/go_pro/gofrTutorial/internal/models"

	gerror "developer.zopsmart.com/go/gofr/pkg/errors"
)

func newMock(*testing.T) (db *sql.DB, mock sqlmock.Sqlmock, store VehicleStore, ctx *gofr.Context) {
	db, mock, err := sqlmock.New(sqlmock.QueryMatcherOption(sqlmock.QueryMatcherEqual))

	if err != nil {
		log.Println(err)
	}

	store = New()

	ctx = gofr.NewContext(nil, nil, &gofr.Gofr{DataStore: datastore.DataStore{ORM: db}})
	ctx.Context = context.Background()

	return
}
func TestGetIdVehicle(t *testing.T) {
	db, mock, VehicleStore, ctx := newMock(t)
	defer db.Close()

	query := "Select Id, Model, Color, NumberPlate, UpdatedAt, CreatedAt, Name, Launched from Vehicle where id=? and DeletedAt is NULL;"

	testcases := []struct {
		ID          int
		Model       string
		Color       string
		NumberPlate string
		UpdatedAt   time.Time
		CreatedAt   time.Time
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
			UpdatedAt:   time.Now(),
			CreatedAt:   time.Now(),
			mockQuery:   mock.ExpectQuery(query).WithArgs(1).WillReturnError(gerror.Error("Internal Server Error")),
			expectError: gerror.Error("Internal Server Error"),
		},
		{
			ID:          1,
			Model:       "i8",
			Color:       "Black",
			NumberPlate: "MH 03 AT 007",
			Name:        "BMW",
			Launched:    true,
			UpdatedAt:   time.Now(),
			CreatedAt:   time.Now(),
			mockQuery: mock.ExpectQuery(query).
				WithArgs(1).
				WillReturnRows(mock.NewRows([]string{"Id", "Model", "Color", "NumberPlate", "UpdatedAt", "CreatedAt", "Name", "Launched"}).
					AddRow(1, "i8", "Black", "MH 03 AT 007", time.Date(2022, 1, 1, 1, 1, 1, 1, time.UTC),
						time.Date(2022, 1, 1, 1, 1, 1, 1, time.UTC), "BMW", true)),
			expectError: nil,
		},
	}
	for _, testCase := range testcases {
		testCase := testCase

		t.Run("", func(t *testing.T) {
			_, err := VehicleStore.GetVehicleByID(ctx, testCase.ID)
			if err != nil && !reflect.DeepEqual(testCase.expectError, err) {
				t.Errorf("expected error: %v, got: %v ", testCase.expectError, err)
			}
		})
	}
}
func TestInsertVehicle(t *testing.T) {
	db, mock, VehicleStore, ctx := newMock(t)
	defer db.Close()

	testcases := []struct {
		ID          int
		Model       string
		Color       string
		NumberPlate string
		UpdatedAt   time.Time
		CreatedAt   time.Time
		Name        string
		Launched    bool
		mockQuery   interface{}
		expectError error
	}{
		{
			Model: "i8", Color: "Black", NumberPlate: "MH 03 AT 007", Name: "BMW", Launched: true,
			mockQuery: mock.
				ExpectExec("insert into Vehicle (Model,Color,NumberPlate,Name,Launched) values(?,?,?,?,?);").
				WithArgs("i8", "Black", "MH 03 AT 007", "BMW", true).
				WillReturnResult(sqlmock.NewResult(1, 1)),
			expectError: nil,
		},
		{
			Model: "vs", Color: "grey", NumberPlate: "hr 12 9373", Name: "toyota", Launched: true,
			mockQuery: mock.
				ExpectExec("insert into Vehicle (Model,Color,NumberPlate,Name,Launched) values(?,?,?,?,?);").
				WithArgs("vs", "grey", "hr 12 9373", "toyota", true).
				WillReturnResult(sqlmock.NewResult(1, 1)),
			expectError: nil,
		},
		{
			Model: "vs", Color: "grey", NumberPlate: "hr 12 9373", Name: "toyota", Launched: true,
			mockQuery: mock.
				ExpectExec("insert into Vehicle (Model,Color,NumberPlate,Name,Launched) values(?,?,?,?,?);").
				WithArgs("vs", "grey", "hr 12 9373", "toyota", true).
				WillReturnResult(sqlmock.NewResult(1, 1)).
				WillReturnError(gerror.Error("Internal Server Error")),
			expectError: gerror.Error("Internal Server Error"),
		},
	}
	for _, testCase := range testcases {
		testCase := testCase
		vehicle := models.Vehicle{ID: testCase.ID, Model: testCase.Model, Color: testCase.Color,
			NumberPlate: testCase.NumberPlate, CreatedAt: testCase.CreatedAt,
			UpdatedAt: testCase.UpdatedAt, Name: testCase.Name, Launched: testCase.Launched}
		_, err := VehicleStore.InsertVehicle(ctx, &vehicle)

		if err != nil && !reflect.DeepEqual(testCase.expectError, err) {
			t.Errorf("expected error: %v, got: %v ", testCase.expectError, err)
		}
	}
}
func TestAllVehicle(t *testing.T) {
	db, mock, VehicleStore, ctx := newMock(t)
	defer db.Close()

	testcases := []struct {
		ID          int
		Model       string
		Color       string
		NumberPlate string
		UpdatedAt   time.Time
		CreatedAt   time.Time
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
			UpdatedAt:   time.Now(),
			CreatedAt:   time.Now(),
			Name:        "BMW",
			Launched:    true,
			mockQuery: mock.
				ExpectQuery("select Id, Model, Color, NumberPlate, UpdatedAt, CreatedAt, Name, Launched from Vehicle where DeletedAt is NULL").
				WillReturnRows(sqlmock.
					NewRows([]string{"Id", "Model", "Color", "NumberPlate", "UpdatedAt", "CreatedAt", "Name", "Launched"}).
					AddRow(1, "i8", "Black", "MH 03 AT 007", time.Now(), time.Now(), "BMW", true)).WillReturnError(nil),
			expectError: nil,
		},
	}
	for _, testCase := range testcases {
		testCase := testCase

		t.Run("", func(t *testing.T) {
			_, err := VehicleStore.GetAll(ctx)
			if err != nil && !reflect.DeepEqual(testCase.expectError, err) {
				t.Errorf("expected error: %v, got: %v ", testCase.expectError, err)
			}
		})
	}
}

func TestUpdateVeh(t *testing.T) {
	db, mock, VehicleStore, ctx := newMock(t)
	defer db.Close()

	testcases := []struct {
		ID          int
		Model       string
		Color       string
		NumberPlate string
		UpdatedAt   time.Time
		CreatedAt   time.Time
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
			UpdatedAt:   time.Now(),
			CreatedAt:   time.Now(),
			Name:        "BMW",
			Launched:    true,
			mockQuery: mock.ExpectExec(
				"update Vehicle set model=?,color=?,numberPlate=?,name=? where Id=? and DeletedAt is null").
				WithArgs(sqlmock.AnyArg(), sqlmock.AnyArg(), sqlmock.AnyArg(), sqlmock.AnyArg(), sqlmock.AnyArg()).
				WillReturnError(gerror.Error("Internal Server Error")),
			expectError: gerror.Error("Internal Server Error"),
		},
		{
			ID:          1,
			Model:       "i8",
			Color:       "Black",
			NumberPlate: "MH 03 AT 007",
			Name:        "BMW",
			UpdatedAt:   time.Now(),
			CreatedAt:   time.Now(),
			Launched:    true,
			expectError: nil,
			mockQuery: []interface{}{mock.ExpectExec(
				"update Vehicle set model=?,color=?,numberPlate=?,name=? where Id=? and DeletedAt is null").
				WithArgs(sqlmock.AnyArg(), sqlmock.AnyArg(), sqlmock.AnyArg(), sqlmock.AnyArg(), 1).
				WillReturnResult(sqlmock.NewResult(1, 1)),
				mock.
					ExpectQuery(`Select Id, Model, Color, NumberPlate, UpdatedAt, 
					CreatedAt, Name, Launched from Vehicle where id=? and DeletedAt is NULL;`).
					WithArgs(1).
					WillReturnRows(mock.
						NewRows([]string{"Id", "Model", "Color", "NumberPlate", "UpdatedAt", "CreatedAt", "Name", "Launched"}).
						AddRow(1, "i8", "Black", "MH 03 AT 007", time.Date(2022, 1, 1, 1, 1, 1, 1, time.UTC),
							time.Date(2022, 1, 1, 1, 1, 1, 1, time.UTC), "BMW", true)),
			},
		},
		{
			ID:          1,
			Model:       "",
			Color:       "",
			NumberPlate: "",
			UpdatedAt:   time.Now(),
			CreatedAt:   time.Now(),
			Name:        "",
			Launched:    true,
			expectError: gerror.Error("Empty Query"),
			// mockQuery: mock.ExpectExec(
			// 	"update Vehicle set where Id=? and DeletedAt is null").
			// 	WithArgs(sqlmock.AnyArg(), sqlmock.AnyArg(), sqlmock.AnyArg(), sqlmock.AnyArg(), 1).
			// 	WillReturnError(errors.New("nothing to update")),
		},
	}
	for _, testCase := range testcases {
		testCase := testCase

		t.Run("", func(t *testing.T) {
			vehicle :=
				models.Vehicle{ID: testCase.ID,
					Model:       testCase.Model,
					Color:       testCase.Color,
					NumberPlate: testCase.NumberPlate,
					CreatedAt:   testCase.CreatedAt,
					UpdatedAt:   testCase.UpdatedAt,
					Name:        testCase.Name,
					Launched:    testCase.Launched}
			_, err := VehicleStore.UpdateVehByID(ctx, testCase.ID, &vehicle)
			if err != nil && !reflect.DeepEqual(testCase.expectError, err) {
				t.Errorf("expected error: %v, got: %v", testCase.expectError, err)
			}
		})
	}
}

func TestDelete(t *testing.T) {
	db, mock, VehicleStore, ctx := newMock(t)
	defer db.Close()

	testcases := []struct {
		ID          int
		Model       string
		Color       string
		NumberPlate string
		UpdatedAt   time.Time
		CreatedAt   time.Time
		Name        string
		Launched    bool
		mockQuery   interface{}
		expectError error
	}{
		{
			ID:    2,
			Model: "vs", Color: "grey", NumberPlate: "hr 12 9373", Name: "toyota", Launched: true,
			mockQuery: mock.ExpectPrepare("UPDATE Vehicle SET DeletedAt= ? WHERE Id=?").
				ExpectExec().WithArgs(sqlmock.AnyArg(), 2).
				WillReturnResult(sqlmock.NewResult(1, 1)),
			expectError: nil,
		},
		{
			ID:    1,
			Model: "vs", Color: "grey", NumberPlate: "hr 12 9373", Name: "toyota", Launched: true,
			mockQuery: mock.ExpectPrepare(
				"UPDATE Vehicle SET DeletedAt= ? WHERE Id=?").
				ExpectExec().WithArgs(sqlmock.AnyArg(), 1).
				WillReturnResult(sqlmock.NewResult(0, 1)),
			expectError: nil,
		},

		{
			ID:    1,
			Model: "vs", Color: "grey", NumberPlate: "hr 12 9373", Name: "toyota", Launched: true,
			mockQuery: mock.ExpectPrepare(
				"UPDATE Vehicle SET DeletedAt= ? WHERE Id=?").
				ExpectExec().
				WithArgs(sqlmock.AnyArg(), 1).
				WillReturnError(gerror.Error("Internal Server Error")),
			expectError: gerror.Error("Internal Server Error"),
		},

		{
			ID:    1,
			Model: "vs", Color: "grey", NumberPlate: "hr 12 9373", Name: "toyota", Launched: true,
			mockQuery: mock.ExpectPrepare("UPDATE Vehicle SET DeletedAt= ? WHERE Id=?").
				WillReturnError(gerror.Error("Internal Server Error")),
			expectError: gerror.Error("Internal Server Error"),
		},
	}
	for _, testCase := range testcases {
		testCase := testCase

		t.Run("", func(t *testing.T) {
			err := VehicleStore.DeleteVehicle(ctx, testCase.ID)
			if err != nil && !reflect.DeepEqual(testCase.expectError, err) {
				t.Errorf("expected error: %v, got: %v ", testCase.expectError, err)
			}
		})
	}
}
