package authentication

import (
	"github.com/gin-gonic/gin"
	authRequests "go-auth-otp-service/src/api/http/requests/authentication"
	response "go-auth-otp-service/src/api/http/responses"
	"go-auth-otp-service/src/pkg/validator"
	"go-auth-otp-service/src/services/authentication"

	"net/http"
)

type RegisterController struct {
	RegisterService authentication.IRegisterService
}

func (controller *RegisterController) Register(c *gin.Context) {
	// Bind check payload.
	var req authRequests.RegisterRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		response.Api(c).SetLog().Send()
		return
	}

	// validate the payload.
	if err := validator.Validate(&req, c.GetString("locale")); err != nil {
		response.Api(c).SetErrors(err).SetLog().Send()
		return
	}

	err := controller.RegisterService.SaveStateAndSendOTP(&req)
	if err != nil {
		response.Api(c).SetMessage(err.Error()).SetLog().Send()
		return
	}

	// Return response.
	response.Api(c).SetMessage("request-successful").
		SetStatusCode(http.StatusOK).SetMessage("register-request-successful").SetLog().Send()
}
