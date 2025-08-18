package providers

import (
	"go-auth-otp-service/src/api/http/middlewares"
	"go-auth-otp-service/src/services"
)

func ProvideRateLimiterService() *services.RateLimitService {
	return &services.RateLimitService{
		Limiter:     services.DefaultLimiter(),
		KeyGetter:   services.DefaultKeyGetter,
		ExcludedKey: nil,
	}
}

func ProvideRateLimiterMiddleware(rateLimiterService *services.RateLimitService) *middlewares.RateLimiterMiddleware {
	return &middlewares.RateLimiterMiddleware{
		RateLimitService: rateLimiterService,
	}
}
