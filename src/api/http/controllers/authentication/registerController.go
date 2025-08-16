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

	key, err := controller.RegisterService.SaveStateAndSendOTP(&req)
	if err != nil {
		response.Api(c).SetMessage(err.Error()).SetLog().Send()
		return
	}

	// Return response.
	response.Api(c).SetMessage("request-successful").
		SetStatusCode(http.StatusOK).SetMessage("register-request-successful").SetData(
		map[string]interface{}{
			"key": key,
		}).Send()
}

func (controller *RegisterController) VerifyRegister(c *gin.Context) {
	// Bind check payload.
	var req authRequests.VerifyRegisterOTP
	if err := c.ShouldBindJSON(&req); err != nil {
		response.Api(c).SetLog().Send()
		return
	}

	// validate the payload.
	if err := validator.Validate(&req, c.GetString("locale")); err != nil {
		response.Api(c).SetErrors(err).SetLog().Send()
		return
	}
	// register user
	err := controller.RegisterService.VerifyRegisterOTPViaRedisKey(&req)
	if err != nil {
		resp := response.Api(c).SetMessage(err.Error())
		resp.SetLog().Send()
		return
	}
	// Return response.
	response.Api(c).SetMessage("verify-register-request-successful").
		SetStatusCode(http.StatusOK).
		SetLog().
		Send()
}
