package authentication

import (
	"go-auth-otp-service/src/api/errs"
	"go-auth-otp-service/src/api/http/requests/authentication"
	"go-auth-otp-service/src/services"
	"gorm.io/gorm"
)

type RegisterService struct {
	UserService services.IUserService
	OTPService  services.IOTPService
}

type IRegisterService interface {
	SaveStateAndSendOTP(req *authentication.RegisterRequest) error
}

func (service *RegisterService) SaveStateAndSendOTP(req *authentication.RegisterRequest) error {
	user, err := service.UserService.GetByNationalIdentityCode(req.NationalIdentityCode)
	if err != nil && err.Error() != gorm.ErrRecordNotFound.Error() {
		return errs.RegisterFailed
	}
	if user != nil {
		return errs.RegisterFailed
	}

	err = service.OTPService.RequestOTP(req.Mobile)
	if err != nil {
		return err
	}

	return nil
}
