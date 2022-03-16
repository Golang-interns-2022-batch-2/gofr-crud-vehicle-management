package vehicle

import (
	"fmt"
	"strconv"
	"time"

	"developer.zopsmart.com/go/gofr/pkg/errors"
	"developer.zopsmart.com/go/gofr/pkg/gofr"
	"github.com/SN786/gofr_vms/internal/model"
	"gopkg.in/guregu/null.v4"
)

type Storer struct {
}

func New() *Storer {
	return &Storer{}
}
func (s Storer) GetDetailsByID(ctx *gofr.Context, id int) (*model.Vehicle, error) {
	var vehResponse model.Vehicle

	query := "select id,model,color,numberPlate,updatedAt,createdAt,name,launched from vehicles where id=? and deletedAt is NULL"
	err := ctx.DB().QueryRowContext(ctx, query, id).
		Scan(&vehResponse.ID, &vehResponse.Model, &vehResponse.Color, &vehResponse.NumberPlate,
			&vehResponse.UpdatedAt, &vehResponse.CreatedAt, &vehResponse.Name, &vehResponse.Launched)

	if err != nil {
		return nil, errors.Error(fmt.Sprintf("couldn't get the vechicle: %v", id))
	}

	return &vehResponse, nil
}
func (s Storer) InsertVehicle(ctx *gofr.Context, vehicle *model.Vehicle) (*model.Vehicle, error) {
	res, err := ctx.DB().ExecContext(ctx, "insert into vehicles (model,color,numberPlate,name,launched) values(?,?,?,?,?)",
		vehicle.Model, vehicle.Color, vehicle.NumberPlate, vehicle.Name, vehicle.Launched)

	if err != nil {
		return nil, errors.Error("couldn't insert the vechicle")
	}

	id, _ := res.LastInsertId()
	vehicle.ID = id

	return vehicle, nil
}

func (s Storer) DeleteVehicleByID(ctx *gofr.Context, id int) error {
	var curentDateTime = time.Now()

	res, err := ctx.DB().ExecContext(ctx, "update vehicles set deletedAt=? where id=? and deletedAt is null", curentDateTime, id)

	if err != nil {
		return errors.Error(fmt.Sprintf("couldn't delete the vechicle: %v", id))
	}

	rowAffected, _ := res.RowsAffected()

	if rowAffected == 0 {
		ID := strconv.Itoa(id)

		return errors.EntityNotFound{
			Entity: "vehicle",
			ID:     ID,
		}
	}

	return nil
}

func (s Storer) UpdateVehicleByID(ctx *gofr.Context, vehicle *model.Vehicle) error {
	var curentDateTime = time.Now()

	query, values := qureyGenerate(vehicle)

	values = append(values, curentDateTime, vehicle.ID)
	_, err := ctx.DB().ExecContext(ctx, query, values...)

	if err != nil {
		fmt.Println(err.Error())
		return errors.Error(fmt.Sprintf("couldn't update the vechicle: %v", vehicle.ID))
	}

	return nil
}

func (s Storer) GetAll(ctx *gofr.Context) ([]*model.Vehicle, error) {
	res, err := ctx.DB().
		QueryContext(ctx, "select id,model,color,numberPlate,updatedAt,createdAt,name,launched from vehicles where deletedAt is NULL")

	if err != nil {
		return nil, errors.Error("couldn't get list of vechicles:")
	}

	var (
		vehicles    []*model.Vehicle
		ID          int64
		Model       string
		Color       string
		NumberPlate string
		UpdatedAt   string
		CreatedAt   string
		Name        string
		Launched    null.Bool
	)

	for res.Next() {
		err := res.Scan(&ID, &Model, &Color, &NumberPlate, &UpdatedAt, &CreatedAt, &Name, &Launched)

		if err != nil {
			return nil, errors.Error("scanning error")
		}

		vehicle := &model.Vehicle{
			ID:          ID,
			Model:       Model,
			Color:       Color,
			NumberPlate: NumberPlate,
			UpdatedAt:   UpdatedAt,
			CreatedAt:   CreatedAt,
			Name:        Name,
			Launched:    Launched,
		}

		vehicles = append(vehicles, vehicle)
	}

	return vehicles, nil
}

func qureyGenerate(vehicles *model.Vehicle) (query string, values []interface{}) {
	if vehicles.Model != "" {
		query += "model=?,"

		values = append(values, vehicles.Model)
	}

	if vehicles.Name != "" {
		query += "name=?,"

		values = append(values, vehicles.Name)
	}

	if vehicles.Launched.Valid {
		query += "launched=?,"

		values = append(values, vehicles.Launched)
	}

	query = "update vehicles set " + query + "updatedAt=? where id=? and deletedAt is NULL"

	return
}
