package vehicle

import (
	"fmt"
	"strconv"
	"time"

	"developer.zopsmart.com/go/gofr/pkg/errors"
	"developer.zopsmart.com/go/gofr/pkg/gofr"
	"github.com/SukantArora/CRUD_Gofr/internal/models"
	"github.com/SukantArora/CRUD_Gofr/internal/store"
)

type dbStore struct {
}

func New() store.Vehicle {
	return &dbStore{}
}

func generateQuery(vehicle *models.Vehicle) (sqlQuery string, parameters []interface{}) {
	query := ""

	var values []interface{}

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

	if vehicle.Launched.Valid {
		query += "launched=?,"

		values = append(values, vehicle.Launched)
	}

	if len(query) > 0 {
		query = query[:len(query)-1]
	}

	return query, values
}

const squery = "select id,model,color,numberPlate,updatedAt,createdAt,name,launched from Vehicle where id = ? and deletedAt is NULL"

func (s *dbStore) GetByID(ctx *gofr.Context, id int) (*models.Vehicle, error) {
	vehicle := &models.Vehicle{}
	err := ctx.DB().QueryRowContext(ctx, squery, id).
		Scan(&vehicle.ID, &vehicle.Model, &vehicle.Color, &vehicle.NumberPlate,
			&vehicle.UpdatedAt, &vehicle.CreatedAt, &vehicle.Name, &vehicle.Launched)

	if err != nil {
		return nil, errors.Error("internal server error")
	}

	return vehicle, nil
}

func (s *dbStore) Create(ctx *gofr.Context, vehicle *models.Vehicle) (*models.Vehicle, error) {
	query := "INSERT INTO Vehicle (model,color,numberPlate,name,launched) values (?,?,?,?,?);"
	res, err := ctx.DB().Exec(query, vehicle.Model, vehicle.Color, vehicle.NumberPlate, vehicle.Name, vehicle.Launched)

	if err != nil {
		return nil, errors.Error("internal server error")
	}

	id, _ := res.LastInsertId()

	return s.GetByID(ctx, int(id))
}

func (s *dbStore) Get(ctx *gofr.Context) ([]*models.Vehicle, error) {
	rows, err := ctx.DB().Query("Select id,model,color,numberPlate,updatedAt,createdAt,name,launched from Vehicle where deletedAt is NULL;")

	var vehicles []*models.Vehicle

	if err != nil {
		return nil, errors.Error("internal server error")
	}
	defer rows.Close()

	for rows.Next() {
		vehicle := &models.Vehicle{}
		err2 := rows.Scan(&vehicle.ID, &vehicle.Model, &vehicle.Color, &vehicle.NumberPlate,
			&vehicle.UpdatedAt, &vehicle.CreatedAt, &vehicle.Name, &vehicle.Launched)

		if err2 == nil {
			vehicles = append(vehicles, vehicle)
		} else {
			return []*models.Vehicle{}, errors.Error("internal server error")
		}
	}

	return vehicles, nil
}

// Update
func (s *dbStore) Update(ctx *gofr.Context, id int, vehicle *models.Vehicle) (*models.Vehicle, error) {
	query := "update Vehicle set "
	q, values := generateQuery(vehicle)

	if q == "" {
		fmt.Println("Hello")
		return nil, errors.Error("nothing to update")
	}

	query += q + " where id=? and deletedAt is null"

	values = append(values, id)
	_, err := ctx.DB().Exec(query, values...)

	if err != nil {
		return nil, errors.Error("internal server error")
	}

	respvehicle, err := s.GetByID(ctx, id)

	return respvehicle, err
}

// Delete
func (s *dbStore) Delete(ctx *gofr.Context, id int) error {
	stmt, err := ctx.DB().Prepare("update Vehicle set deletedAt = ? where id=?")
	if err != nil {
		return errors.Error("internal server error")
	}
	defer stmt.Close()
	res, err := stmt.Exec(time.Now(), id)

	if err != nil {
		return errors.Error("internal server error")
	}

	rowsAffect, err := res.RowsAffected()
	if rowsAffect == 0 {
		return errors.EntityNotFound{Entity: "Movie", ID: strconv.Itoa(id)}
	}

	return err
}
