package vehicle

import (
	"database/sql"

	filter "github.com/Gofr-VMS/Filters"
	"github.com/Gofr-VMS/model"
	"github.com/Gofr-VMS/store"

	"developer.zopsmart.com/go/gofr/pkg/errors"
	"developer.zopsmart.com/go/gofr/pkg/gofr"
)

type vehicle struct{}

func New() store.Store {
	return &vehicle{}
}
func (v *vehicle) GetByID(ctx *gofr.Context, id int) (*model.Vehicles, error) {
	vh := &model.Vehicles{}

	query := "SELECT id,model,color,numPlate,createdAt,updatedAt,name,launched FROM vehicle WHERE id=? AND deletedAt is NULL"
	row := ctx.DB().QueryRowContext(ctx, query, id)
	err := row.Scan(&vh.ID, &vh.Model, &vh.Color, &vh.NumPlate, &vh.CreatedAt, &vh.UpdatedAt, &vh.Name, &vh.Launched)

	if err != nil {
		return nil, sql.ErrNoRows
	}

	return vh, nil
}
func (v *vehicle) GetAll(ctx *gofr.Context) ([]*model.Vehicles, error) {
	var vhcl []*model.Vehicles

	query := "SELECT id,model,color,numPlate,createdAt,updatedAt,name,launched FROM vehicle WHERE  deletedAt is NULL"

	results, err := ctx.DB().QueryContext(ctx, query)

	if err != nil {
		return nil, errors.EntityNotFound{Entity: "vehicle", ID: "all"}
	}

	for results.Next() {
		var vh = &model.Vehicles{}
		err = results.Scan(&vh.ID, &vh.Model, &vh.Color, &vh.NumPlate, &vh.CreatedAt, &vh.UpdatedAt, &vh.Name, &vh.Launched)

		if err != nil {
			return nil, sql.ErrNoRows
		}

		vhcl = append(vhcl, vh)
	}

	return vhcl, nil
}
func (v *vehicle) Insert(ctx *gofr.Context, vh *model.Vehicles) (*model.Vehicles, error) {
	query := "INSERT INTO vehicle(model,color,numPlate,name,launched) VALUES(?,?,?,?,?)"

	r, e := ctx.DB().ExecContext(ctx, query, vh.Model, vh.Color, vh.NumPlate, vh.Name, vh.Launched)
	if e != nil {
		return nil, errors.EntityAlreadyExists{}
	}

	lastID, _ := r.LastInsertId()
	vh.ID = int(lastID)

	return v.GetByID(ctx, int(lastID))
}
func (v *vehicle) Update(ctx *gofr.Context, vh *model.Vehicles) (*model.Vehicles, error) {
	query := "update vehicle set "
	where, val := filter.UpdateFilter(vh)

	query += where
	query += " where id=? and deletedAt is NULL"

	val = append(val, vh.ID)
	_, e := ctx.DB().ExecContext(ctx, query, val...)

	if e != nil {
		return nil, errors.Error("error updating record")
	}

	return v.GetByID(ctx, vh.ID)
}
func (v *vehicle) Delete(ctx *gofr.Context, id int) error {
	query := "update vehicle set deletedAt=Now() WHERE  id=? "
	result, err := ctx.DB().ExecContext(ctx, query, id)

	if err != nil {
		return errors.EntityNotFound{Entity: "vehicle"}
	}

	res, _ := result.RowsAffected()

	if res == 0 {
		return errors.Error("Delete failed")
	}

	return nil
}
