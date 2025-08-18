package services

import (
	"fmt"
	"github.com/gin-gonic/gin"
	"github.com/gin-gonic/gin/binding"
	"github.com/ulule/limiter/v3"
	sredis "github.com/ulule/limiter/v3/drivers/store/redis"
	"go-auth-otp-service/src/api/errs"
	authRequests "go-auth-otp-service/src/api/http/requests/authentication"
	response "go-auth-otp-service/src/api/http/responses"
	"go-auth-otp-service/src/cache"
	"go-auth-otp-service/src/config"
	"log"
	"net/http"
	"strconv"
	"time"
)

type IRateLimitService interface {
	GetLimiter() *limiter.Limiter
	GetExcludedKey() func(string) bool
	KeyFunc(*gin.Context) string
	ExcludedKeyFunc(string) bool
	SetLimiter(limiter *limiter.Limiter) IRateLimitService
	SetKey(keyGetter func(*gin.Context) string) IRateLimitService
	SetExcludedKey(excludedKey func(string) bool) IRateLimitService
	GetOnExceedHandler(c *gin.Context, resetIn string) bool
	PostOnExceedHandler(c *gin.Context, resetIn string) bool
}

type RateLimitService struct {
	Limiter     *limiter.Limiter
	KeyGetter   func(*gin.Context) string
	ExcludedKey func(string) bool
}

func (s *RateLimitService) KeyFunc(ctx *gin.Context) string {
	return s.KeyGetter(ctx)
}

func (s *RateLimitService) ExcludedKeyFunc(key string) bool {
	return s.ExcludedKey(key)
}

func (s *RateLimitService) GetExcludedKey() func(string) bool {
	return s.ExcludedKey
}

func (s *RateLimitService) GetLimiter() *limiter.Limiter {
	return s.Limiter
}

func (s *RateLimitService) SetLimiter(limiter *limiter.Limiter) IRateLimitService {
	s.Limiter = limiter
	return s
}

func (s *RateLimitService) SetKey(keyGetter func(*gin.Context) string) IRateLimitService {
	s.KeyGetter = keyGetter
	return s
}

func (s *RateLimitService) SetExcludedKey(excludedKey func(string) bool) IRateLimitService {
	s.ExcludedKey = excludedKey
	return s
}

func (s *RateLimitService) GetOnExceedHandler(c *gin.Context, resetIn string) bool {
	response.Api(c).SetStatusCode(http.StatusTooManyRequests).
		SetMessage(errs.TooManyRequest.Error()).SetLog().Send()
	return false
}

func (s *RateLimitService) PostOnExceedHandler(c *gin.Context, resetIn string) bool {
	response.Api(c).SetStatusCode(http.StatusTooManyRequests).
		SetMessage(fmt.Sprintf(errs.TooManyRequest.Error())).
		SetLog().Send()
	return false
}

func DefaultLimiter() *limiter.Limiter {
	strRate := config.GetInstance().Get("RATE_LIMITER_DEFAULT_LIMIT")
	strPeriod := config.GetInstance().Get("RATE_LIMITER_DEFAULT_PERIOD_PER_SECOND")

	if strRate == "" {
		strRate = "60"
	}

	if strPeriod == "" {
		strPeriod = "60"
	}

	rate, err := strconv.Atoi(strRate)
	if err != nil {
		log.Fatal(err)
	}

	period, err := strconv.Atoi(strPeriod)
	if err != nil {
		log.Fatal(err)
	}

	limiterRate := limiter.Rate{
		Period: time.Duration(period) * time.Second,
		Limit:  int64(rate),
	}

	store, err := sredis.NewStore(cache.GetInstance().GetClient())
	if err != nil {
		log.Fatal(err)
	}

	return limiter.New(store, limiterRate)
}

func DefaultKeyGetter(c *gin.Context) string {
	return fmt.Sprintf("default-%s", c.GetString("request-ip"))
}

func CriticalLimiter() *limiter.Limiter {
	strRate := config.GetInstance().Get("RATE_LIMITER_CRITICAL_LIMIT")
	strPeriod := config.GetInstance().Get("RATE_LIMITER_CRITICAL_PERIOD_PER_SECOND")

	if strRate == "" {
		strRate = "3"
	}

	if strPeriod == "" {
		strPeriod = "60"
	}

	rate, err := strconv.Atoi(strRate)
	if err != nil {
		log.Fatal(err)
	}

	period, err := strconv.Atoi(strPeriod)
	if err != nil {
		log.Fatal(err)
	}

	limiterRate := limiter.Rate{
		Period: time.Duration(period) * time.Second,
		Limit:  int64(rate),
	}

	store, err := sredis.NewStore(cache.GetInstance().GetClient())
	if err != nil {
		log.Fatal(err)
	}

	return limiter.New(store, limiterRate)
}

func RegisterCriticalLimiter() *limiter.Limiter {
	strRate := config.GetInstance().Get("RATE_LIMITER_REGISTER_CRITICAL_LIMIT")
	strPeriod := config.GetInstance().Get("RATE_LIMITER_REGISTER_CRITICAL_PERIOD_PER_SECOND")
	return SetLimiter(strRate, strPeriod)
}

func LoginCriticalLimiter() *limiter.Limiter {
	strRate := config.GetInstance().Get("RATE_LIMITER_LOGIN_CRITICAL_LIMIT")
	strPeriod := config.GetInstance().Get("RATE_LIMITER_LOGIN_CRITICAL_PERIOD_PER_SECOND")
	return SetLimiter(strRate, strPeriod)
}

func AdminLoginCriticalLimiter() *limiter.Limiter {
	strRate := config.GetInstance().Get("RATE_LIMITER_ADMIN_LOGIN_CRITICAL_LIMIT")
	strPeriod := config.GetInstance().Get("RATE_LIMITER_ADMIN_LOGIN_CRITICAL_PERIOD_PER_SECOND")
	return SetLimiter(strRate, strPeriod)
}

func SetLimiter(strRate, strPeriod string) *limiter.Limiter {
	if strRate == "" {
		strRate = "3"
	}

	if strPeriod == "" {
		strPeriod = "60"
	}

	rate, err := strconv.Atoi(strRate)
	if err != nil {
		log.Fatal(err)
	}

	period, err := strconv.Atoi(strPeriod)
	if err != nil {
		log.Fatal(err)
	}

	limiterRate := limiter.Rate{
		Period: time.Duration(period) * time.Second,
		Limit:  int64(rate),
	}

	store, err := sredis.NewStore(cache.GetInstance().GetClient())
	if err != nil {
		log.Fatal(err)
	}

	return limiter.New(store, limiterRate)
}

func CriticalKeyGetter(c *gin.Context) string {
	return fmt.Sprintf("critical-%s", c.GetString("request-ip"))
}

func CriticalVerifyOtpKeySetter(c *gin.Context) string {
	return fmt.Sprintf("critical-verify-login-%s", c.GetString("request-ip"))
}

func GenericCriticalKeyGetter(key string) func(*gin.Context) string {
	return func(c *gin.Context) string {
		return fmt.Sprintf("%s-%s", key, c.GetString("request-ip"))
	}
}

func OtpKeyGetter(c *gin.Context) string {
	var req authRequests.AuthSendOtpRequest
	fmt.Println("hi1")
	if err := c.ShouldBindBodyWith(&req, binding.JSON); err == nil && req.Mobile != "" {
		fmt.Println("hi2")

		return fmt.Sprintf("otp-%s", req.Mobile)
	}
	fmt.Println("hi3")

	return fmt.Sprintf("otp-ip-%s", c.ClientIP())
}
