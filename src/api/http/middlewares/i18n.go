package middlewares

import (
	"github.com/gin-gonic/gin"
	"go-auth-otp-service/src/pkg/i18n"
)

func I18n(context *gin.Context) {
	if context != nil {
		if context.GetHeader("Accept-Language") != "" {
			lng := context.GetHeader("Accept-Language")
			if !isLocaleSupported(lng) {
				lng = "fa"
			}
			context.Set("locale", lng)
		}
	}

	context.Next()
}

// Checks if a locale is supported
func isLocaleSupported(lng string) bool {
	for _, locale := range i18n.Locales {
		if lng == locale {
			return true
		}
	}
	return false
}
