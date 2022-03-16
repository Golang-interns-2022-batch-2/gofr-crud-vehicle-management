package vehicle

import (
	"context"
	"database/sql"
	"log"
	"reflect"
	"testing"
	"time"

	"github.com/DATA-DOG/go-sqlmock"
	"github.com/Gofr-VMS/model"
	"github.com/Gofr-VMS/store"

	"developer.zopsmart.com/go/gofr/pkg/datastore"

	"developer.zopsmart.com/go/gofr/pkg/errors"
	"developer.zopsmart.com/go/gofr/pkg/gofr"
)

func setMock() (*gofr.Context, store.Store, sqlmock.Sqlmock, *sql.DB) {
	db, mock, err := sqlmock.New(sqlmock.QueryMatcherOption(sqlmock.QueryMatcherEqual))
	if err != nil {
		log.Println(err)
	}

	ctx := gofr.NewContext(nil, nil, &gofr.Gofr{DataStore: datastore.DataStore{ORM: db}})
	ctx.Context = context.Background()
	ctx.DB().DB = db
	s := New()

	return ctx, s, mock, db
}
func TestGetByID(t *testing.T) {
	ctx, s, mock, db := setMock()

	defer db.Close()

	query := "SELECT id,model,color,numPlate,createdAt,updatedAt,name,launched FROM vehicle WHERE ID=? AND deletedAt is NULL"
	dateTime := time.Date(2022, 1, 1, 1, 1, 1, 1, time.UTC)

	TC := []struct {
		desc   string
		id     int
		output model.Vehicles
		experr error
		mock   interface{}
	}{

		{
			desc: "success",
			id:   1,
			output: model.Vehicles{
				ID:        1,
				Model:     "i7",
				Color:     "Black",
				NumPlate:  "MH 03 AT 007",
				CreatedAt: dateTime,
				UpdatedAt: dateTime,
				Name:      "BMW",
				Launched:  true,
			},
			experr: nil,
			mock: mock.ExpectQuery(query).WithArgs(1).
				WillReturnRows(sqlmock.NewRows([]string{"ID", "model", "color", "numPlate",
					"updatedAt", "createdAt", "name", "launched"}).
					AddRow(1, "i8", "Black", "MH 03 AT 007", dateTime, dateTime, "BMW", true)),
		},
		{
			desc: "fail",
			id:   100,
			experr: errors.EntityNotFound{
				Entity: "vehicle",
				ID:     "100",
			},
			mock: mock.ExpectQuery(query).WithArgs(100).WillReturnError(sql.ErrNoRows),
		},
	}
	for _, tc := range TC {
		t.Run(tc.desc, func(t *testing.T) {
			res, err := s.GetByID(ctx, tc.id)

			if err != nil {
				return
			}
			if tc.experr == nil && !reflect.DeepEqual(res, &tc.output) {
				t.Errorf("%s : expected %v, but got %v", tc.desc, tc.output, res)
			}
		})
	}
}
func TestInsert(t *testing.T) {
	ctx, s, mock, db := setMock()

	defer db.Close()

	query := "INSERT INTO vehicle(model,color,numPlate,name,launched) VALUES(?,?,?,?,?)"

	TC := []struct {
		desc   string
		input  model.Vehicles
		output model.Vehicles
		experr error
		mock   interface{}
	}{
		{
			desc:   "Succ",
			input:  model.Vehicles{ID: 1, Model: "i8", Color: "Black", NumPlate: "MH 03 AT 007", Name: "BMW", Launched: true},
			output: model.Vehicles{ID: 1, Model: "i8", Color: "Black", NumPlate: "MH 03 AT 007", Name: "BMW", Launched: true},
			experr: nil,
			mock: []interface{}{
				*mock.ExpectExec(query).WithArgs(1, "i8", "Black", "MH 03 AT 007", "BMW", 1).WillReturnResult(sqlmock.NewResult(1, 1)),
			},
		},
		{
			desc:   "fail",
			input:  model.Vehicles{ID: 1, Model: "i8", Color: "Black", NumPlate: "MH 03 AT 007", Name: "BMW", Launched: false},
			output: model.Vehicles{},
			experr: errors.EntityAlreadyExists{},
			mock: []interface{}{
				*mock.ExpectExec(query).WithArgs(1, "i8", "Black", "MH 03 AT 007", "BMW", 0).WillReturnError(errors.EntityAlreadyExists{}),
			},
		},
	}
	for _, tc := range TC {
		t.Run(tc.desc, func(t *testing.T) {
			res, err := s.Insert(ctx, &tc.input)

			if err != nil {
				return
			}

			if tc.experr == nil && !reflect.DeepEqual(res, &tc.output) {
				t.Errorf("%s : expected %v, but got %v", tc.desc, tc.output, res)
			}
		})
	}
}

func TestGetAll(t *testing.T) {
	ctx, s, mock, db := setMock()
	defer db.Close()

	dateTime := time.Date(2022, 1, 1, 1, 1, 1, 1, time.UTC)
	query := "SELECT id,model,color,numPlate,createdAt,updatedAt,name,launched FROM vehicle WHERE  deletedAt is NULL"

	TC := []struct {
		desc      string
		experr    error
		expoutput []*model.Vehicles
		mock      interface{}
	}{
		{
			desc:   "Suceess",
			experr: nil,
			expoutput: []*model.Vehicles{
				{
					ID:        1,
					Model:     "i8",
					Color:     "Black",
					NumPlate:  "MH 03 AT 007",
					CreatedAt: dateTime,
					UpdatedAt: dateTime,
					Name:      "BMW",
					Launched:  true,
				},
				{
					ID:        2,
					Model:     "truck",
					Color:     "Brown",
					NumPlate:  "MH 03 A1 657",
					CreatedAt: dateTime,
					UpdatedAt: dateTime,
					Name:      "BMW",
					Launched:  true,
				},
			},

			mock: mock.ExpectQuery(query).WillReturnRows(sqlmock.NewRows(
				[]string{"ID", "model", "color", "numPlate", "updatedAt", "createdAt", "name", "launched"}).
				AddRow(1, "i8", "Black", "MH 03 AT 007", dateTime, dateTime, "BMW", true).
				AddRow(2, "truck", "Browm", "MH 03 A1 657", dateTime, dateTime, "BMW", true)),
		},

		{
			desc:      "fail",
			experr:    errors.EntityNotFound{Entity: "vehicle", ID: "all"},
			expoutput: nil,
			mock:      mock.ExpectQuery(query).WillReturnError(sql.ErrNoRows),
		},
	}
	for _, tc := range TC {
		t.Run(tc.desc, func(t *testing.T) {
			_, err := s.GetAll(ctx)

			if err != nil {
				return
			}

			if tc.experr == nil && !reflect.DeepEqual(err, tc.experr) {
				t.Errorf("%s : expected %v, but got %v", tc.desc, tc.experr, err)
			}
		})
	}
}
func TestDelete(t *testing.T) {
	ctx, s, mock, db := setMock()
	defer db.Close()

	query := "update vehicle set deletedAt=Now() WHERE  id=? "

	TC := []struct {
		desc   string
		input  int
		experr error
		mock   interface{}
	}{
		{
			desc:   "success",
			input:  1,
			experr: nil,
			mock:   mock.ExpectExec(query).WillReturnResult(sqlmock.NewResult(0, 1)),
		},
		{
			desc:   "fail",
			input:  10000,
			experr: errors.EntityNotFound{Entity: "vehicle"},
			mock:   mock.ExpectExec(query).WillReturnError(sql.ErrConnDone),
		},
	}
	for _, tc := range TC {
		t.Run(tc.desc, func(t *testing.T) {
			err := s.Delete(ctx, tc.input)

			if !reflect.DeepEqual(err, tc.experr) {
				t.Errorf("%s : expected %v, but got %v", tc.desc, tc.experr, err)
			}
		})
	}
}
func TestUpdate(t *testing.T) {
	ctx, s, mock, db := setMock()
	defer db.Close()

	query := "update vehicle set model=?,name=?,launched=? where ID=?"
	TC := []struct {
		desc   string
		input  *model.Vehicles
		experr error
		mock   interface{}
	}{
		{
			desc: "success",
			input: &model.Vehicles{
				ID:       1,
				Model:    "Q5",
				Color:    "Black",
				NumPlate: "UK 07 1896",
				Name:     "Audi",
				Launched: true,
			},
			experr: nil,
			mock:   mock.ExpectExec(query).WithArgs(1, "Q5", "Black", "UK 07 1896", "Audi", true).WillReturnResult(sqlmock.NewResult(1, 1)),
		},
		{
			desc: "fail",
			input: &model.Vehicles{
				ID:       2,
				Model:    "truck",
				Name:     "BMW",
				Launched: true,
			},
			experr: errors.Error("error updating record"),
			mock:   mock.ExpectExec(query).WithArgs("truck", "BMW", true).WillReturnError(sql.ErrConnDone),
		},
	}
	for _, tc := range TC {
		t.Run(tc.desc, func(t *testing.T) {
			_, err := s.Update(ctx, tc.input)

			if !reflect.DeepEqual(err, tc.experr) {
				t.Errorf("%s : expected %v, but got %v", tc.desc, tc.experr, err)
			}
		})
	}
}
