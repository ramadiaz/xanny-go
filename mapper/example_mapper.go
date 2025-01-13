package mapper

import (
	"layered-template/models/database"
	"layered-template/models/dto"

	"github.com/go-viper/mapstructure/v2"
)

func MapExampleInputToModel(input dto.ExampleInput) database.Example {
	var example database.Example

	mapstructure.Decode(input, &example)
	return example
}
