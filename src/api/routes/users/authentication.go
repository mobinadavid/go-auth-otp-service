package users

import (
	"github.com/gin-gonic/gin"
	"go-auth-otp-service/src/providers"
)

func AuthenticationRouter(router *gin.RouterGroup) {
	registerController := providers.GetAuthenticationContainer()

	// define route
	authentication := router.Group("authentication")

	// register
	register := authentication.Group("register")
	{
		register.POST("send-otp", registerController.UserRegisterController.SendOtp)
		register.POST("verify-otp", registerController.UserRegisterController.VerifyOtp)
	}

}
