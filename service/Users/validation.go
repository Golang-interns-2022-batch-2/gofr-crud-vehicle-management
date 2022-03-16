package vehicle

import "github.com/Gofr-VMS/model"

func IsIDValid(id int) bool {
	return id > 0
}
func validateModel(vehicle *model.Vehicles) bool {
	return vehicle.Model != ""
}
func validateColor(vehicle *model.Vehicles) bool {
	return vehicle.Color != ""
}
func validateNumPlate(vehicle *model.Vehicles) bool {
	return vehicle.NumPlate != ""
}
func validateName(vehicle *model.Vehicles) bool {
	return vehicle.Name != ""
}
