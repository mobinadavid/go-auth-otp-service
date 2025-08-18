package authentication

import (
	"github.com/gin-gonic/gin"
	"github.com/google/uuid"
	"go-auth-otp-service/src/api/errs"
	authentication_request "go-auth-otp-service/src/api/http/requests/authentication"
	response "go-auth-otp-service/src/api/http/responses"
	"go-auth-otp-service/src/database/scopes"
	"go-auth-otp-service/src/pkg/validator"
	"go-auth-otp-service/src/services/authentication"
	"net/http"
)

type AccessTokenController struct {
	AccessTokenService authentication.IAccessTokenService
}

func (controller *AccessTokenController) GetList(c *gin.Context) {
	// get the query builder
	builder, exists := c.Get("query_parameters_builder")
	if !exists {
		response.Api(c).SetStatusCode(http.StatusUnprocessableEntity).SetMessage(errs.SomeThingWentWrong.Error()).SetLog().Send()
		return
	}

	// fetch information to query builder
	builderModel := builder.(*scopes.BuilderModel)
	builderModel.Filters["owner_type"] = c.GetString("authenticated-user-type")
	builderModel.Filters["owner_id"] = c.GetUint("authenticated-user-id")

	// get list of access token
	data, err := controller.AccessTokenService.GetList(builderModel)
	if err != nil {
		response.Api(c).SetStatusCode(http.StatusNotFound).SetMessage(err.Error()).SetLog().Send()
		return
	}

	// return response
	response.Api(c).SetMessage("request-successful").
		SetStatusCode(http.StatusOK).
		SetData(map[string]interface{}{
			"access_tokens": data,
		}).SetLog().Send()
}

func (controller *AccessTokenController) GetActiveTokens(c *gin.Context) {
	// get the query builder
	builder, exists := c.Get("query_parameters_builder")
	if !exists {
		response.Api(c).SetStatusCode(http.StatusUnprocessableEntity).SetMessage(errs.SomeThingWentWrong.Error()).SetLog().Send()
		return
	}

	// fetch information to query builder
	builderModel := builder.(*scopes.BuilderModel)
	builderModel.Filters["owner_type"] = c.GetString("authenticated-user-type")
	builderModel.Filters["owner_id"] = c.GetUint("authenticated-user-id")

	// get list of access token
	data, err := controller.AccessTokenService.GetActiveTokens(builderModel)
	if err != nil {
		response.Api(c).SetStatusCode(http.StatusNotFound).SetMessage(err.Error()).SetLog().Send()
		return
	}

	// return response
	response.Api(c).SetMessage("request-successful").
		SetStatusCode(http.StatusOK).
		SetData(map[string]interface{}{
			"access_tokens": data,
		}).SetLog().Send()
}

func (controller *AccessTokenController) GetByUuid(c *gin.Context) {
	// Get the UUID from the URL parameter and parse it
	uuidStr := c.Param("uuid")
	id, err := uuid.Parse(uuidStr)
	if err != nil {
		response.Api(c).SetMessage(errs.InvalidUuid.Error()).SetLog().Send()
		return
	}

	// get a token by uuid
	data, err := controller.AccessTokenService.GetByUuid(&id)
	if err != nil {
		response.Api(c).SetStatusCode(http.StatusNotFound).SetMessage(err.Error()).SetLog().Send()
		return
	}

	// send response
	response.Api(c).SetMessage("request-successful").
		SetStatusCode(http.StatusOK).
		SetData(map[string]interface{}{
			"access_token": data,
		}).SetLog().Send()
}

func (controller *AccessTokenController) RefreshAccessToken(c *gin.Context) {
	// Get refresh token from cookies
	var req authentication_request.RefreshAccessTokenRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		response.Api(c).SetLog().Send()
		return
	}

	// validate request
	if err := validator.Validate(&req, c.GetString("locale")); err != nil {
		response.Api(c).SetErrors(err).Send()
		return
	}

	jwt, err := controller.AccessTokenService.RefreshAccessTokens(req.RefreshToken, "user")
	if err != nil {
		response.Api(c).SetMessage(err.Error()).SetStatusCode(http.StatusUnauthorized).SetLog().Send()
		return
	}

	// Return response
	response.Api(c).SetMessage("request-successful").
		SetStatusCode(http.StatusOK).
		SetData(map[string]interface{}{
			"access_tokens": jwt,
		}).SetLog().Send()
}

func (controller *AccessTokenController) RevokeTokens(c *gin.Context) {
	// revoke tokens
	err := controller.AccessTokenService.RevokeTokens(c.GetUint("authenticated-user-id"), c.GetString("authenticated-user-type"))
	if err != nil {
		response.Api(c).SetMessage(err.Error()).SetLog().Send()
		return
	}

	// send response
	response.Api(c).SetMessage("request-successful").
		SetStatusCode(http.StatusOK).
		SetLog().
		Send()
}

func (controller *AccessTokenController) RevokeTokenByUUID(c *gin.Context) {
	// Get the UUID from the URL parameter and parse it
	uuidStr := c.Param("uuid")
	id, err := uuid.Parse(uuidStr)
	if err != nil {
		response.Api(c).SetMessage(errs.InvalidUuid.Error()).SetLog().Send()
		return
	}

	// revoke the token by it's uuid
	err = controller.AccessTokenService.RevokeTokenByUuid(&id, c.GetUint("authenticated-user-id"), c.GetString("authenticated-user-type"))
	if err != nil {
		response.Api(c).SetStatusCode(http.StatusNotFound).SetMessage(err.Error()).SetLog().Send()
		return
	}

	// send response
	response.Api(c).SetMessage("request-successful").
		SetStatusCode(http.StatusOK).
		SetLog().
		Send()
}

func (controller *AccessTokenController) RevokeCurrentToken(c *gin.Context) {
	// Get the UUID from the URL parameter
	uuidStr := c.GetString("access-token-uuid")
	id, err := uuid.Parse(uuidStr)
	if err != nil {
		response.Api(c).SetMessage(errs.InvalidUuid.Error()).SetLog().Send()
		return
	}

	// revoke the token by it's uuid
	err = controller.AccessTokenService.RevokeTokenByUuid(&id, c.GetUint("authenticated-user-id"), c.GetString("authenticated-user-type"))
	if err != nil {
		response.Api(c).SetStatusCode(http.StatusNotFound).SetMessage(err.Error()).SetLog().Send()
		return
	}

	// send response
	response.Api(c).SetMessage("request-successful").
		SetStatusCode(http.StatusOK).
		SetLog().
		Send()
}
