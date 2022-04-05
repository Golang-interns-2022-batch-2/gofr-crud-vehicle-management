package vehicle

import (
	"database/sql"
	"net/http"

	"strconv"
	"time"

	"developer.zopsmart.com/go/gofr/pkg/errors"
	"developer.zopsmart.com/go/gofr/pkg/gofr"
	"github.com/SN786/gofr_vms/internal/model"
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

	if err == sql.ErrNoRows {
		idStr := strconv.Itoa(id)

		return nil, errors.EntityNotFound{
			Entity: "Vehicle",
			ID:     idStr,
		}
	}

	if err != nil {
		return nil, &errors.Response{
			StatusCode: http.StatusInternalServerError,
			Code:       http.StatusText(http.StatusInternalServerError),
			Reason:     "couldn't fetch the vehicle details",
		}
	}

	return &vehResponse, nil
}
func (s Storer) InsertVehicle(ctx *gofr.Context, vehicle *model.Vehicle) (*model.Vehicle, error) {
	res, err := ctx.DB().ExecContext(ctx, "insert into vehicles (model,color,numberPlate,name,launched) values(?,?,?,?,?)",
		vehicle.Model, vehicle.Color, vehicle.NumberPlate, vehicle.Name, vehicle.Launched)

	if err != nil {
		return nil, &errors.Response{
			StatusCode: http.StatusInternalServerError,
			Code:       http.StatusText(http.StatusInternalServerError),
			Reason:     "couldn't insert the vehicle",
		}
	}

	id, _ := res.LastInsertId()
	vehicle.ID = id

	return vehicle, nil
}

func (s Storer) DeleteVehicleByID(ctx *gofr.Context, id int) error {
	var curentDateTime = time.Now()

	res, err := ctx.DB().ExecContext(ctx, "update vehicles set deletedAt=? where id=? and deletedAt is null", curentDateTime, id)

	if err == sql.ErrNoRows {
		idStr := strconv.Itoa(id)

		return errors.EntityNotFound{
			Entity: "Vehicle",
			ID:     idStr,
		}
	}

	if err != nil {
		return &errors.Response{
			StatusCode: http.StatusInternalServerError,
			Code:       http.StatusText(http.StatusInternalServerError),
			Reason:     "couldn't delete the vehicle",
		}
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

	if err == sql.ErrNoRows {
		idStr := strconv.Itoa(int(vehicle.ID))

		return errors.EntityNotFound{
			Entity: "Vehicle",
			ID:     idStr,
		}
	}

	if err != nil {
		return &errors.Response{
			StatusCode: http.StatusInternalServerError,
			Code:       http.StatusText(http.StatusInternalServerError),
			Reason:     "couldn't update the vehicle",
		}
	}

	return nil
}

func (s Storer) GetAll(ctx *gofr.Context) ([]*model.Vehicle, error) {
	res, err := ctx.DB().
		QueryContext(ctx, "select id,model,color,numberPlate,updatedAt,createdAt,name,launched from vehicles where deletedAt is NULL")

	if err == sql.ErrNoRows {
		return nil, errors.EntityNotFound{
			Entity: "Vehicle",
			ID:     "all",
		}
	}

	if err != nil {
		return nil, &errors.Response{
			StatusCode: http.StatusInternalServerError,
			Code:       http.StatusText(http.StatusInternalServerError),
			Reason:     "couldn't get list of vehicle",
		}
	}

	var vehicles []*model.Vehicle

	for res.Next() {
		var vehicle model.Vehicle

		err := res.Scan(&vehicle.ID, &vehicle.Model, &vehicle.Color, &vehicle.NumberPlate,
			&vehicle.UpdatedAt, &vehicle.CreatedAt, &vehicle.Name, &vehicle.Launched)

		if err != nil {
			return nil, &errors.Response{
				StatusCode: http.StatusInternalServerError,
				Code:       http.StatusText(http.StatusInternalServerError),
				Reason:     "couldn't ger the vehicles",
			}
		}

		vehicles = append(vehicles, &vehicle)
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
