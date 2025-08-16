package providers

import (
	authentication_controller "go-auth-otp-service/src/api/http/controllers/authentication"
	"go-auth-otp-service/src/services"
	"go-auth-otp-service/src/services/authentication"
)

func ProvideUserRegisterController(registerService *authentication.RegisterService) *authentication_controller.RegisterController {
	return &authentication_controller.RegisterController{
		RegisterService: registerService,
	}
}

func ProvideRegisterService(userService *services.UserService, otpService *services.OTPService, jwtService *authentication.JwtService, accessTokenService *authentication.AccessTokenService) *authentication.RegisterService {
	return &authentication.RegisterService{
		UserService:        userService,
		OTPService:         otpService,
		AccessTokenService: accessTokenService,
		JwtService:         jwtService,
	}
}

func ProvideJwtService() *authentication.JwtService {
	return &authentication.JwtService{}
}
