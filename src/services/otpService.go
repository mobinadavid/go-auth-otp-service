package services

import (
	"context"
	"crypto/rand"
	"errors"
	"fmt"
	"github.com/redis/go-redis/v9"
	"go-auth-otp-service/src/api/errs"
	"go-auth-otp-service/src/cache"
	"go-auth-otp-service/src/config"
	"math/big"
	"strconv"
	"time"
)

type IOTPService interface {
	RequestOTP(mobile string) error
	VerifyOTP(mobile, otp string) (bool, error)
	generateOTP(length int) (string, error)
}

type OTPService struct {
}

func (service *OTPService) RequestOTP(mobile string) error {
	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	// get a key for the otp
	key := getRedisKey(mobile)

	// get if otp exist don't let new otp be create
	exists, err := cache.GetInstance().GetClient().Exists(ctx, key).Result()
	if err != nil {
		return errs.SomeThingWentWrong
	}

	if exists != 0 {
		return errs.ErrAuthOTPExists
	}

	otpLength, err := strconv.Atoi(config.GetInstance().Get("OTP_LENGTH"))
	if err != nil {
		return errs.SomeThingWentWrong
	}

	// generate a otp
	otpExpiration, err := time.ParseDuration(config.GetInstance().Get("OTP_EXPIRATION") + "s")
	otp, err := service.generateOTP(otpLength)
	if err != nil {
		return errs.SomeThingWentWrong
	}

	// Send otp
	//here we print it in the console
	//Todo://use notification service
	fmt.Println("this is the otp:", otp)
	// set key:otp in redis
	err = cache.GetInstance().GetClient().Set(ctx, key, otp, otpExpiration).Err()
	if err != nil {
		return errs.SomeThingWentWrong
	}

	return nil
}

func (service *OTPService) VerifyOTP(mobile, otp string) (bool, error) {
	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	// get a key for the otp
	key := getRedisKey(mobile)

	// get the value from redis
	storedOTP, err := cache.GetInstance().GetClient().Get(ctx, key).Result()

	// check otp exist
	if errors.Is(err, redis.Nil) {
		return false, nil
	} else if err != nil {
		return false, errs.SomeThingWentWrong
	}

	// check otp value is valid
	if storedOTP != otp {
		return false, nil
	}

	// remove otp in redis if it's ok
	err = cache.GetInstance().GetClient().Del(ctx, key).Err()
	if err != nil {
		return false, errs.SomeThingWentWrong
	}

	return true, nil
}

func (service *OTPService) generateOTP(length int) (string, error) {
	const charset = "0123456789"
	otp := make([]byte, length)
	for i := range otp {
		num, err := rand.Int(rand.Reader, big.NewInt(int64(len(charset))))
		if err != nil {
			return "", err
		}
		otp[i] = charset[num.Int64()]
	}
	return string(otp), nil
}

func getRedisKey(mobile string) string {
	return fmt.Sprintf("otp-%s", mobile)
}
