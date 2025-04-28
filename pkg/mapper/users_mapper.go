package mapper

import (
	"xanny-go-template/api/users/dto"
	"xanny-go-template/models"

	"github.com/go-viper/mapstructure/v2"
)

func MapUserInputToModel(input dto.Users) models.Users {
	var user models.Users

	mapstructure.Decode(input, &user)
	return user
}
