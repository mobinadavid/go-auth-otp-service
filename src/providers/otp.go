package providers

import "go-auth-otp-service/src/services"

func ProvideOTPService() *services.OTPService {
	return &services.OTPService{}
}
