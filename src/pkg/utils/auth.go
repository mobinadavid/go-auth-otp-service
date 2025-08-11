package utils

import "github.com/gin-gonic/gin"

type Auth struct {
	OwnerId         string
	OwnerType       string
	TokenUUID       string
	IsAuthenticated bool
}

func GetAuthData(c *gin.Context) *Auth {
	auth := &Auth{IsAuthenticated: false}
	authenticatedUserID := c.GetString("authenticated-user-id")
	authenticatedUserType := c.GetString("authenticated-user-type")
	accessTokenUUID := c.GetString("access-token-uuid")

	if authenticatedUserID != "" && authenticatedUserType != "" && accessTokenUUID != "" {
		auth.IsAuthenticated = true
		auth.TokenUUID = accessTokenUUID
		auth.OwnerType = authenticatedUserType
		auth.OwnerId = authenticatedUserID
	}

	return auth
}
