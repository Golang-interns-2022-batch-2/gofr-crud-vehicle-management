package filter

import "github.com/Gofr-VMS/model"

func UpdateFilter(v *model.Vehicles) (where string, query []interface{}) {
	if v.Model != "" {
		where += " model=?,"

		query = append(query, v.Model)
	}

	if v.Color != "" {
		where += " color=?,"

		query = append(query, v.Color)
	}

	if v.NumPlate != "" {
		where += " numPlate=?,"

		query = append(query, v.NumPlate)
	}

	if v.Name != "" {
		where += " name=?,"

		query = append(query, v.Name)
	}

	if len(where) > 0 {
		where = where[:len(where)-1]
	}

	return
}
