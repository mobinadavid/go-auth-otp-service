package routes

import (
	"github.com/gin-gonic/gin"
	"go-auth-otp-service/src/api/http/middlewares"
	"go-auth-otp-service/src/models"
	"go-auth-otp-service/src/providers"
)

func UserRouter(router *gin.RouterGroup) {
	userContainer := providers.GetUserContainer()
	authenticationContainer := providers.GetAuthenticationContainer()

	// define route
	users := router.Group("users")

	// user
	{
		users.GET("", authenticationContainer.AuthenticationMiddleware.Middleware("user"),
			middlewares.QueryParametersBuilderMiddleware(models.UserModel{}),
			userContainer.UserController.GetList)

		users.GET(":uuid", authenticationContainer.AuthenticationMiddleware.Middleware("user"),
			userContainer.UserController.GetByUuid)
	}

}
