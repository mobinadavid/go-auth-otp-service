package routes

import (
	"github.com/gin-gonic/gin"
	"go-auth-otp-service/src/providers"
)

func UserRouter(router *gin.RouterGroup) {
	userContainer := providers.GetUserContainer()
	authenticationContainer := providers.GetAuthenticationContainer()

	// define route
	users := router.Group("users")

	// user
	{
		users.GET("", authenticationContainer.AuthenticationMiddleware.Middleware("user"), userContainer.UserController.GetList)
		users.GET(":uuid", authenticationContainer.AuthenticationMiddleware.Middleware("user"), userContainer.UserController.GetByUuid)
	}

}
