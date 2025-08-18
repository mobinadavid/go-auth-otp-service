package controllers

import (
	"github.com/gin-gonic/gin"
	"github.com/google/uuid"
	"go-auth-otp-service/src/api/errs"
	response "go-auth-otp-service/src/api/http/responses"
	"go-auth-otp-service/src/database/scopes"
	"go-auth-otp-service/src/services"
	"net/http"
)

type UserController struct {
	UserService services.IUserService
}

func (controller *UserController) GetList(c *gin.Context) {
	builder, exists := c.Get("query_parameters_builder")
	if !exists {
		response.Api(c).SetStatusCode(http.StatusUnprocessableEntity).SetMessage(errs.SomeThingWentWrong.Error()).SetLog().Send()
		return
	}

	// prepare data for service
	builderModel := builder.(*scopes.BuilderModel)

	// get user list
	users, err := controller.UserService.GetList(builderModel)
	if err != nil {
		response.Api(c).SetStatusCode(http.StatusNotFound).SetMessage(err.Error()).SetLog().Send()
		return
	}

	// send response
	response.Api(c).
		SetMessage("request-successful").
		SetStatusCode(http.StatusOK).
		SetData(map[string]interface{}{
			"users": users,
		}).SetLog().Send()
}

func (controller *UserController) GetByUuid(c *gin.Context) {
	// Get the UUID from the URL parameter
	uuidStr := c.Param("uuid")
	id, err := uuid.Parse(uuidStr)
	if err != nil {
		response.Api(c).SetMessage(errs.InvalidUuid.Error()).SetLog().Send()
		return
	}

	// Use the UserRepository to find the user by UUID
	user, err := controller.UserService.GetByUuid(&id)
	if err != nil {
		response.Api(c).SetStatusCode(http.StatusNotFound).SetMessage(err.Error()).SetLog().Send()
		return
	}

	response.Api(c).
		SetMessage("request-successful").
		SetStatusCode(http.StatusOK).
		SetData(map[string]interface{}{
			"user": user,
		}).SetLog().Send()
}
