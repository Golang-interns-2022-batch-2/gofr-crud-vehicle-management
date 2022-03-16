package vehicle

import (
	"developer.zopsmart.com/go/gofr/pkg/errors"
	"github.com/SukantArora/CRUD_Gofr/internal/models"
)

func validateInput(vehicle *models.Vehicle) (bool, error) {
	if vehicle.Name == "" {
		return false, errors.InvalidParam{Param: []string{"Name"}}
	}

	if vehicle.Color == "" {
		return false, errors.InvalidParam{Param: []string{"Color"}}
	}

	if vehicle.Model == "" {
		return false, errors.InvalidParam{Param: []string{"Model"}}
	}

	if vehicle.NumberPlate == "" {
		return false, errors.InvalidParam{Param: []string{"NumberPlate"}}
	}

	return true, nil
}

func validateID(id int) bool {
	return id > 0
}

// func validateName(name string) bool {
// 	return len(name) > 0
// }

// func validateColor(color string) bool {
// 	return len(color) > 0
// }

// func validateModel(model string) bool {
// 	return len(model) > 0
// }

// func validateNumberPlate(numberPlate string) bool {
// 	return len(numberPlate) > 0
// }
