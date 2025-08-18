package authentication

import (
	"github.com/gin-gonic/gin"
	authRequests "go-auth-otp-service/src/api/http/requests/authentication"
	response "go-auth-otp-service/src/api/http/responses"
	"go-auth-otp-service/src/pkg/validator"
	"go-auth-otp-service/src/services/authentication"
	"golang.org/x/net/context"

	"net/http"
)

type RegisterController struct {
	RegisterService authentication.IRegisterService
}

func (controller *RegisterController) SendOtp(c *gin.Context) {
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

func (controller *RegisterController) VerifyOtp(c *gin.Context) {
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
	// prepare data for service
	ctx := context.WithValue(context.Background(), "request-ip", c.GetString("request-ip"))
	ctx = context.WithValue(ctx, "request-user-agent", c.GetHeader("User-Agent"))
	// register user
	jwt, err := controller.RegisterService.VerifyRegisterOTPViaRedisKey(ctx, &req)
	if err != nil {
		resp := response.Api(c).SetMessage(err.Error())
		resp.SetLog().Send()
		return
	}
	// Return response.
	response.Api(c).SetMessage("verify-register-request-successful").
		SetStatusCode(http.StatusOK).
		SetData(map[string]interface{}{
			"access_tokens": jwt,
		}).
		SetLog().
		Send()
}
