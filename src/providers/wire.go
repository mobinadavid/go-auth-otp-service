//go:build wireinject
// +build wireinject

package providers

import (
	"github.com/google/wire"
	"go-auth-otp-service/src/api/http/controllers"
	authentication2 "go-auth-otp-service/src/api/http/controllers/authentication"
	"go-auth-otp-service/src/api/http/middlewares"
	"go-auth-otp-service/src/database"
)

type (
	AuthenticationContainer struct {
		UserRegisterController   *authentication2.RegisterController
		AuthenticationMiddleware *middlewares.AuthenticationMiddleware
	}
	UserContainer struct {
		UserController *controllers.UserController
	}
)

func GetAuthenticationContainer() *AuthenticationContainer {
	wire.Build(
		// Repositories
		database.GetInstance,
		ProvideUserRepository,
		ProvideAccessTokenRepository,
		// Services
		ProvideRegisterService,
		ProvideUserService,
		ProvideOTPService,
		ProvideJwtService,
		ProvideAccessTokenService,
		// Controllers
		ProvideUserRegisterController,

		// Middlewares
		ProvideAuthenticationMiddleware,

		wire.Struct(new(AuthenticationContainer), "*"),
	)
	return nil
}

func GetUserContainer() *UserContainer {
	wire.Build(
		// Repositories
		database.GetInstance,
		ProvideUserRepository,
		// Services
		ProvideUserService,
		// Controllers
		ProvideUserController,
		wire.Struct(new(UserContainer), "*"),
	)
	return nil
}
