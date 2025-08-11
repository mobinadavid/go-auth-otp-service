package services

import (
	"github.com/google/uuid"
	"go-auth-otp-service/src/api/errs"
	"go-auth-otp-service/src/models"
	"go-auth-otp-service/src/repositories"
)

type IUserService interface {
	GetByNationalIdentityCode(nationalIdentityCode string) (*models.UserModel, error)
	Update(user *models.UserModel) (*models.UserModel, error)
	GetByUuid(uuid *uuid.UUID) (*models.UserModel, error)
}

type UserService struct {
	UserRepository repositories.IUserRepository
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
func (service *UserService) Update(user *models.UserModel) (*models.UserModel, error) {
	res, err := service.UserRepository.Update(user)
	if err != nil {
		return nil, errs.SomeThingWentWrong
	}

	return res, nil
}
