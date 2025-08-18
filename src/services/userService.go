package services

import (
	"github.com/google/uuid"
	"go-auth-otp-service/src/api/errs"
	"go-auth-otp-service/src/api/http/requests/userRequests"
	"go-auth-otp-service/src/database/scopes"
	"go-auth-otp-service/src/models"
	"go-auth-otp-service/src/repositories"
)

type IUserService interface {
	GetList(builder *scopes.BuilderModel) (*scopes.PaginateModel, error)
	Update(user *models.UserModel) (*models.UserModel, error)
	GetByUuid(uuid *uuid.UUID) (*models.UserModel, error)
	Create(request *userRequests.CreateRequest) (*models.UserModel, error)
	GetByNationalIdentityCode(nationalIdentityCode string) (*models.UserModel, error)
	GetByMobile(mobile string) (*models.UserModel, error)
}

type UserService struct {
	UserRepository repositories.IUserRepository
}

func (service *UserService) GetList(builder *scopes.BuilderModel) (*scopes.PaginateModel, error) {
	res, err := service.UserRepository.GetList(builder)
	if err != nil {
		return nil, errs.SomeThingWentWrong
	}

	return res, nil
}

func (service *UserService) GetByUuid(uuid *uuid.UUID) (*models.UserModel, error) {
	res, err := service.UserRepository.GetByUuid(uuid)
	if err != nil {
		return nil, errs.SomeThingWentWrong
	}

	return res, nil
}

func (service *UserService) GetByNationalIdentityCode(nationalIdentityCode string) (*models.UserModel, error) {
	res, err := service.UserRepository.GetByNationalIdentityCode(nationalIdentityCode)
	if err != nil {
		return nil, err
	}
	return res, nil
}

func (service *UserService) GetByMobile(mobile string) (*models.UserModel, error) {
	res, err := service.UserRepository.GetByMobile(mobile)
	if err != nil {
		return nil, err
	}
	return res, nil
}

func (service *UserService) Update(user *models.UserModel) (*models.UserModel, error) {
	res, err := service.UserRepository.Update(user)
	if err != nil {
		return nil, errs.SomeThingWentWrong
	}

	return res, nil
}

func (service *UserService) Create(request *userRequests.CreateRequest) (*models.UserModel, error) {
	user := &models.UserModel{
		Uuid:                 uuid.New(),
		FirstName:            request.FirstName,
		LastName:             request.LastName,
		NationalIdentityCode: request.NationalIdentityCode,
		Mobile:               request.Mobile,
		Password:             []byte(request.Password),
	}
	userOrm, err := service.UserRepository.Create(user)
	if err != nil {
		return nil, errs.SomeThingWentWrong
	}

	return userOrm, nil
}
