// go:build wireinject
// go:build wireinject
//go:build wireinject
// +build wireinject

package injectors

import (
	authControllers "xanny-go-template/internal/auth/controllers"
	authServices "xanny-go-template/internal/auth/services"

	"github.com/go-playground/validator/v10"
	"github.com/google/wire"
)

var authFeatureSet = wire.NewSet(
	authServices.NewComponentServices,
	authControllers.NewCompController,
)

func InitializeAuthController(validate *validator.Validate) authControllers.CompControllers {
	wire.Build(authFeatureSet)
	return nil
}
