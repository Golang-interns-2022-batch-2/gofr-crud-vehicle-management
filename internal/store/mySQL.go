package store

import (
	"errors"
	"time"

	gerror "developer.zopsmart.com/go/gofr/pkg/errors"
	"developer.zopsmart.com/go/gofr/pkg/gofr"
	"github.com/ishankochar09/go_pro/gofrTutorial/internal/models"
)

type dbStore struct {
}

func New() VehicleStore {
	return &dbStore{}
}

func (s *dbStore) GetVehicleByID(ctx *gofr.Context, id int) (models.Vehicle, error) {
	var veh models.Vehicle

	query := "Select Id, Model, Color, NumberPlate, UpdatedAt, CreatedAt, Name, Launched from Vehicle where id=? and DeletedAt is NULL;"
	row := ctx.DB().QueryRowContext(ctx, query, id)
	err := row.Scan(&veh.ID, &veh.Model, &veh.Color, &veh.NumberPlate, &veh.UpdatedAt, &veh.CreatedAt, &veh.Name, &veh.Launched)

	if err != nil {
		return models.Vehicle{}, gerror.Error("Internal Server Error")
	}

	return veh, nil
}
func (s *dbStore) InsertVehicle(ctx *gofr.Context, veh *models.Vehicle) (models.Vehicle, error) {
	res, err := ctx.DB().ExecContext(ctx,
		"insert into Vehicle (Model,Color,NumberPlate,Name,Launched) values(?,?,?,?,?);",
		veh.Model,
		veh.Color,
		veh.NumberPlate,
		veh.Name,
		veh.Launched)

	if err != nil {
		return models.Vehicle{}, gerror.Error("Internal Server Error")
	}

	id, _ := res.LastInsertId()
	resp, _ := s.GetVehicleByID(ctx, int(id))

	return resp, nil
}
func (s *dbStore) GetAll(ctx *gofr.Context) ([]models.Vehicle, error) {
	rows, err := ctx.DB().QueryContext(ctx, `select Id, Model, Color, NumberPlate, UpdatedAt, CreatedAt, 
	Name, Launched from Vehicle where DeletedAt is NULL`)

	var vehicles []models.Vehicle

	var vehi models.Vehicle

	if err == nil {
		defer rows.Close()

		for rows.Next() {
			e := rows.Scan(&vehi.ID, &vehi.Model, &vehi.Color, &vehi.NumberPlate, &vehi.UpdatedAt, &vehi.CreatedAt, &vehi.Name, &vehi.Launched)
			if e == nil {
				vehi = models.Vehicle{
					ID:          vehi.ID,
					Model:       vehi.Model,
					Color:       vehi.Color,
					NumberPlate: vehi.NumberPlate,
					UpdatedAt:   vehi.UpdatedAt,
					CreatedAt:   vehi.CreatedAt,
					Name:        vehi.Name, Launched: vehi.Launched}
				vehicles = append(vehicles, vehi)
			}
		}

		return vehicles, nil
	}

	return []models.Vehicle{}, errors.New("error exec select query")
}
func (s *dbStore) DeleteVehicle(ctx *gofr.Context, id int) error {
	res, err := ctx.DB().PrepareContext(ctx, "UPDATE Vehicle SET DeletedAt= ? WHERE Id=?")

	if err != nil {
		return gerror.Error("Internal Server Error")
	}
	defer res.Close()
	_, err = res.Exec(time.Now(), id)

	if err != nil {
		return gerror.Error("Internal Server Error")
	}

	return err
}
func (s *dbStore) UpdateVehByID(ctx *gofr.Context, id int, veh *models.Vehicle) (models.Vehicle, error) {
	query := "update Vehicle set "
	q, values := generateQuery(veh)

	if q == "" {
		return models.Vehicle{}, gerror.Error("Empty Query")
	}

	query += q

	query += " where Id=? and DeletedAt is null"

	values = append(values, id)
	_, err := ctx.DB().ExecContext(ctx, query, values...)

	if err != nil {
		return models.Vehicle{}, err
	}

	respvehicle, err := s.GetVehicleByID(ctx, id)

	return respvehicle, err
}
func generateQuery(vehicle *models.Vehicle) (query string, values []interface{}) {
	query = ""

	values = []interface{}{}

	if vehicle.Model != "" {
		query += "model=?,"

		values = append(values, vehicle.Model)
	}

	if vehicle.Color != "" {
		query += "color=?,"

		values = append(values, vehicle.Color)
	}

	if vehicle.NumberPlate != "" {
		query += "numberPlate=?,"

		values = append(values, vehicle.NumberPlate)
	}

	if vehicle.Name != "" {
		query += "name=?,"

		values = append(values, vehicle.Name)
	}

	if len(query) > 0 {
		query = query[:len(query)-1]
	}

	return query, values
}
