package routes

import (
	"github.com/gin-gonic/gin"
	"go-auth-otp-service/src/api/http/middlewares"
	"go-auth-otp-service/src/models"
	"go-auth-otp-service/src/providers"
)

func RegisterAccessTokensRouter(router *gin.RouterGroup) {
	authenticationContainer := providers.GetAuthenticationContainer()
	accessTokens := router.Group("access-tokens")
	{
		accessTokens.GET("", authenticationContainer.AuthenticationMiddleware.Middleware("user"),
			middlewares.QueryParametersBuilderMiddleware(models.AccessTokenModel{}),
			authenticationContainer.AccessTokenController.GetList)

		accessTokens.GET(":uuid",
			authenticationContainer.AuthenticationMiddleware.Middleware("user"),
			authenticationContainer.AccessTokenController.GetByUuid,
		)

		accessTokens.POST("refresh", authenticationContainer.AccessTokenController.RefreshAccessToken)

		revoke := accessTokens.Group("revoke").
			Use(authenticationContainer.AuthenticationMiddleware.Middleware("user"))
		{
			revoke.DELETE("", authenticationContainer.AccessTokenController.RevokeTokens)
			revoke.DELETE(":uuid", authenticationContainer.AccessTokenController.RevokeTokenByUUID)
			revoke.DELETE("current-token", authenticationContainer.AccessTokenController.RevokeCurrentToken)
		}
	}

	activeAccessTokens := router.Group("active-access-tokens")
	{
		activeAccessTokens.GET("", authenticationContainer.AuthenticationMiddleware.Middleware("user"),
			middlewares.QueryParametersBuilderMiddleware(models.AccessTokenModel{}),
			authenticationContainer.AccessTokenController.GetActiveTokens)
	}
}
