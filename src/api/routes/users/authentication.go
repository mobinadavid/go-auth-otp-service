package users

import (
	"github.com/gin-gonic/gin"
	"go-auth-otp-service/src/providers"
	"go-auth-otp-service/src/services"
)

func AuthenticationRouter(router *gin.RouterGroup) {
	registerController := providers.GetAuthenticationContainer()

	rateLimiterSendOtp := providers.ProvideRateLimiterMiddleware(
		providers.ProvideRateLimiterService(),
	).SetLimiter(services.CriticalLimiter()).SetKey(services.OtpKeyGetter)

	rateLimiterVerifyOtp := providers.ProvideRateLimiterMiddleware(
		providers.ProvideRateLimiterService(),
	).SetLimiter(services.CriticalLimiter()).SetKey(services.GenericCriticalKeyGetter("verify-otp"))

	// define route
	authentication := router.Group("authentication")

	// register
	register := authentication.Group("register")
	{
		register.POST("send-otp", rateLimiterSendOtp.Middleware, registerController.UserRegisterController.SendOtp)
		register.POST("verify-otp", rateLimiterVerifyOtp.Middleware, registerController.UserRegisterController.VerifyOtp)
	}

}
