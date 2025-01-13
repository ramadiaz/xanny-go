package mapper

import (
	"xanny-go-template/models"
	"xanny-go-template/dto"

	"github.com/go-viper/mapstructure/v2"
)

func MapExampleInputToModel(input dto.ExampleInput) models.Example {
	var example models.Example

	mapstructure.Decode(input, &example)
	return example
}
