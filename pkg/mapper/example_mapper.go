package mapper

import (
	"xanny-go-template/api/example/dto"
	"xanny-go-template/models"

	"github.com/go-viper/mapstructure/v2"
)

func MapExampleInputToModel(input dto.ExampleInput) models.Example {
	var example models.Example

	mapstructure.Decode(input, &example)
	return example
}
