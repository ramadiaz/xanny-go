package mapper

import (
	"xanny-go-template/models/database"
	"xanny-go-template/models/dto"

	"github.com/go-viper/mapstructure/v2"
)

func MapExampleInputToModel(input dto.ExampleInput) database.Example {
	var example database.Example

	mapstructure.Decode(input, &example)
	return example
}
