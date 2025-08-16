package authentication

import (
	"context"
	"encoding/json"
	"github.com/google/uuid"
	"github.com/redis/go-redis/v9"
	"go-auth-otp-service/src/api/errs"
	"go-auth-otp-service/src/api/http/requests/authentication"
	"go-auth-otp-service/src/api/http/requests/userRequests"
	"go-auth-otp-service/src/cache"
	"go-auth-otp-service/src/config"
	"go-auth-otp-service/src/services"
	"gorm.io/gorm"
	"strconv"
	"time"
)

type RegisterService struct {
	UserService services.IUserService
	OTPService  services.IOTPService
}

type IRegisterService interface {
	SaveStateAndSendOTP(req *authentication.RegisterRequest) (string, error)
	VerifyRegisterOTPViaRedisKey(req *authentication.VerifyRegisterOTP) error
}

func (service *RegisterService) SaveStateAndSendOTP(req *authentication.RegisterRequest) (string, error) {
	user, err := service.UserService.GetByNationalIdentityCode(req.NationalIdentityCode)
	if err != nil && err.Error() != gorm.ErrRecordNotFound.Error() {
		return "", errs.RegisterFailed
	}
	if user != nil {
		return "", errs.ErrUserAlreadyExists
	}
	// marshal the req to save in redis
	reqData, err := json.Marshal(req)
	if err != nil {
		return "", errs.RegisterFailed
	}

	// get expire time
	expiration, err := strconv.Atoi(config.GetInstance().Get("REGISTER_SAVE_STATE_LIFETIME"))
	if err != nil {
		expiration = 120
	}

	// Save the request data in Redis
	key := uuid.New().String()
	err = cache.GetInstance().GetClient().Set(context.Background(), key, reqData, time.Duration(expiration)*time.Second).Err()
	if err != nil {
		return "", errs.RegisterFailed
	}
	err = service.OTPService.RequestOTP(req.Mobile)
	if err != nil {
		return "", err
	}

	return key, nil
}

func (service *RegisterService) VerifyRegisterOTPViaRedisKey(req *authentication.VerifyRegisterOTP) error {
	res, err := cache.GetInstance().GetClient().Get(context.Background(), req.RegisterKey).Result()
	if err != nil {
		if err == redis.Nil {
			return errs.RegisterFailed
		}
		return errs.SomeThingWentWrong
	}

	var resp authentication.RegisterRequest
	err = json.Unmarshal([]byte(res), &resp)
	if err != nil {
		return errs.SomeThingWentWrong
	}

	var otpIsValid bool
	otpIsValid, err = service.OTPService.VerifyOTP(resp.Mobile, req.OTP)
	if err != nil {
		return errs.SomeThingWentWrong
	}

	if !otpIsValid {
		return errs.ErrOTPInvalid
	}

	_, err = service.UserService.Create(&userRequests.CreateRequest{
		//Todo:add more fields
		NationalIdentityCode: resp.NationalIdentityCode,
		Mobile:               resp.Mobile,
	})
	if err != nil {
		return errs.SomeThingWentWrong
	}

	return nil
}
