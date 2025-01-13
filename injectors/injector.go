// go:build wireinject
// go:build wireinject
//go:build wireinject
// +build wireinject

package injectors

import (
	exampleControllers "layered-template/example/controllers"
	exampleRepositories "layered-template/example/repositories"
	exampleServices "layered-template/example/services"

	"github.com/go-playground/validator/v10"
	"github.com/google/wire"
	"gorm.io/gorm"
)

var exampleFeatureSet = wire.NewSet(
	exampleRepositories.NewComponentRepository,
	exampleServices.NewComponentServices,
	exampleControllers.NewCompController,
)

func InitializeExampleController(db *gorm.DB, validate *validator.Validate) exampleControllers.CompControllers {
	wire.Build(exampleFeatureSet)
	return nil
}