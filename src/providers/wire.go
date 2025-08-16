//go:build wireinject
// +build wireinject

package providers

import (
	"github.com/google/wire"
	authentication2 "go-auth-otp-service/src/api/http/controllers/authentication"
	"go-auth-otp-service/src/database"
)

type (
	AuthenticationContainer struct {
		UserRegisterController *authentication2.RegisterController
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

		wire.Struct(new(AuthenticationContainer), "*"),
	)
	return nil
}
