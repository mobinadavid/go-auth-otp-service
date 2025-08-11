package validator

import (
	"errors"
	"github.com/go-playground/locales/ar"
	"github.com/go-playground/locales/en"
	"github.com/go-playground/locales/es"
	"github.com/go-playground/locales/fa"
	"github.com/go-playground/locales/fr"
	"github.com/go-playground/locales/tr"
	ut "github.com/go-playground/universal-translator"
	"github.com/go-playground/validator/v10"
	artranslations "github.com/go-playground/validator/v10/translations/ar"
	entranslations "github.com/go-playground/validator/v10/translations/en"
	estranslations "github.com/go-playground/validator/v10/translations/es"
	fatranslations "github.com/go-playground/validator/v10/translations/fa"
	frtranslations "github.com/go-playground/validator/v10/translations/fr"
	trtranslations "github.com/go-playground/validator/v10/translations/tr"
)

var (
	validate *validator.Validate
	uni      *ut.UniversalTranslator
)

// Validate function that takes a struct and a locale, then performs validation with localized messages.
func Validate(s interface{}, locale string) map[string]string {
	// Initialize the Universal Translator.
	uni = ut.New(en.New(), en.New(), fa.New(), es.New(), fr.New(), ar.New(), tr.New())

	// Get the validator instance.
	validate = validator.New()

	// Get the translator for the given locale, default to English
	translator, found := uni.GetTranslator(locale)
	if !found {
		translator, _ = uni.GetTranslator("fa")
	}

	// Register translations for the locale.
	switch locale {
	case "en":
		_ = entranslations.RegisterDefaultTranslations(validate, translator)
	case "fr":
		_ = frtranslations.RegisterDefaultTranslations(validate, translator)
	case "es":
		_ = estranslations.RegisterDefaultTranslations(validate, translator)
	case "ar":
		_ = artranslations.RegisterDefaultTranslations(validate, translator)
	case "tr":
		_ = trtranslations.RegisterDefaultTranslations(validate, translator)
	default:
		_ = fatranslations.RegisterDefaultTranslations(validate, translator)
	}

	// Attach Custom rules.
	RegisterRules(validate, uni)

	// Perform validation.
	err := validate.Struct(s)

	// Translate errors.
	if err != nil {
		errMap := make(map[string]string)
		var valErrs validator.ValidationErrors
		if errors.As(err, &valErrs) {
			for _, err := range err.(validator.ValidationErrors) {
				errMap[err.Field()] = err.Translate(translator)
			}
			return errMap // Return the translated error messages
		} else {
			errMap["error"] = err.Error()
			return errMap
		}
	}

	// No errors found.
	return nil
}
