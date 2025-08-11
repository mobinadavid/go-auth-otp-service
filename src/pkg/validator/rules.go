package validator

import (
	"fmt"
	"github.com/gin-gonic/gin/binding"
	ut "github.com/go-playground/universal-translator"
	"github.com/go-playground/validator/v10"
	"github.com/google/uuid"
	"go-auth-otp-service/src/database"
	"go-auth-otp-service/src/pkg/utils"
	"log"
	"regexp"
	"strconv"
	"strings"
	"time"
)

func RegisterRules(val *validator.Validate, trans *ut.UniversalTranslator) {
	// Map of rule names to their corresponding validation functions
	ruleToFunc := map[string]validator.Func{
		"jalali-date":                    jalaliDateValidation,
		"iranian-national-identity-code": iranianNationalCodeValidation,
		"iranian-company-national-id":    iranianCompanyNationalIdValidation,
		"iranian-mobile":                 iranianMobileValidation,
		"iranian-phone":                  iranianPhoneValidation,
		"iranian-postal-code":            iranianPostalCodeValidation,
		"iranian-bank-card-number":       iranianBankCardNumberValidation,
		"iranian-bank-sheba-number":      iranianBankShebaNumberValidation,
		"is-uuid":                        isValidUuid,
		"exists":                         exists,
		"unique":                         unique,
		"is-rfc3339":                     isRfc3339,
	}

	for ruleName, ruleFunc := range ruleToFunc {

		// Register the validation.
		_ = val.RegisterValidation(ruleName, ruleFunc)

		if v, ok := binding.Validator.Engine().(*validator.Validate); ok {
			err := v.RegisterValidation(ruleName, ruleFunc)
			if err != nil {
				log.Println(err)
			}
		}

	}
}

// jalaliDateValidation validates a date string according to specific rules:
// The year must be between 1300 and 1400, the month between 01 and 12, and the day between 01 and 31.
func jalaliDateValidation(fl validator.FieldLevel) bool {
	matched, _ := regexp.MatchString(`^(13\d\d|14\d\d)/(0[1-9]|1[0-2])/(0[1-9]|[12]\d|3[01])$`, fl.Field().String())
	return matched
}

// iranianNationalCodeValidation validates the iranian national code.
func iranianNationalCodeValidation(fl validator.FieldLevel) bool {
	value := fl.Field().String()
	matched, _ := regexp.MatchString(`^\d{8,10}$`, value)
	if !matched {
		return false
	}

	sequentialMatch, _ := regexp.MatchString(`^[0]{10}|[1]{10}|[2]{10}|[3]{10}|[4]{10}|[5]{10}|[6]{10}|[7]{10}|[8]{10}|[9]{10}$`, value)
	if sequentialMatch {
		return false
	}

	value = fmt.Sprintf("%010s", value)
	var sub int
	for i, char := range value[:9] {
		digit, _ := strconv.Atoi(string(char))
		sub += digit * (10 - i)
	}

	var control int
	if sub%11 < 2 {
		control = sub % 11
	} else {
		control = 11 - (sub % 11)
	}

	lastDigit, _ := strconv.Atoi(string(value[9]))
	return lastDigit == control
}

// iranianMobileValidation validates the iranian mobile number.
// valid example: 9123456789
func iranianMobileValidation(fl validator.FieldLevel) bool {
	matched, _ := regexp.MatchString(`^09[0-9]{9}$`, fl.Field().String())
	return matched
}

// iranianPhoneValidation validates the iranian phone number with area code.
// valid example: 02123456789
func iranianPhoneValidation(fl validator.FieldLevel) bool {
	matched, _ := regexp.MatchString(`^(0[1-9]{2})[2-9][0-9]{7}$`, fl.Field().String())
	return matched
}

// iranianPostalCodeValidation validates the iranian postal code.
func iranianPostalCodeValidation(fl validator.FieldLevel) bool {
	matched, _ := regexp.MatchString(`\b[13-9]{5}[0-9]{5}\b`, fl.Field().String())
	return matched
}

// iranianBankCardNumberValidation validates Iranian bank payment card number.
// depending on 'http://www.aliarash.com/article/creditcart/credit-debit-cart.htm'.
func iranianBankCardNumberValidation(fl validator.FieldLevel) bool {
	// Check if the card number is exactly 16 digits
	match, _ := regexp.MatchString(`^\d{16}$`, fl.Field().String())
	if !match {
		return false
	}

	sum := 0
	for position := 1; position <= 16; position++ {
		digit, _ := strconv.Atoi(string(fl.Field().String()[position-1]))
		if position%2 != 0 {
			digit = digit * 2
		}
		if digit > 9 {
			digit = digit - 9
		}

		sum += digit
	}

	return sum%10 == 0
}

// iranianBankShebaNumberValidation validates an Iranian bank Sheba number (IBAN).
func iranianBankShebaNumberValidation(fl validator.FieldLevel) bool {
	// Remove non-alphanumeric characters and convert to uppercase
	value := regexp.MustCompile(`[\W_]+`).ReplaceAllString(fl.Field().String(), "")
	value = strings.ToUpper(value)

	// Validate length and structure
	if len(value) < 4 || len(value) > 34 || !utils.IsLetter(value[0]) || !utils.IsLetter(value[1]) || !utils.IsDigit(value[2]) || !utils.IsDigit(value[3]) {
		return false
	}

	// Prepare the string for modulo operation
	ibanReplaceChars := ""
	for i := 'A'; i <= 'Z'; i++ {
		ibanReplaceChars += fmt.Sprintf("%c", i)
	}

	ibanReplaceValues := make([]string, 26)
	for i := range ibanReplaceChars {
		ibanReplaceValues[i] = strconv.Itoa(10 + i)
	}
	tmpIBAN := value[4:] + value[:4]
	for i, char := range ibanReplaceChars {
		tmpIBAN = strings.ReplaceAll(tmpIBAN, string(char), ibanReplaceValues[i])
	}

	// Perform the modulo operation
	isValid, _ := strconv.ParseInt(tmpIBAN[:1], 10, 64)
	for i := 1; i < len(tmpIBAN); i++ {
		num, _ := strconv.ParseInt(tmpIBAN[i:i+1], 10, 64)
		isValid = (isValid*10 + num) % 97
	}

	return isValid == 1
}

// isValidUuid Custom validator function to validate UUID format
func isValidUuid(fl validator.FieldLevel) bool {
	_, err := uuid.Parse(fl.Field().String())
	return err == nil
}

func exists(fl validator.FieldLevel) bool {
	value := fl.Field().String()

	// Parse the tag to extract table name and column name
	tag := fl.Param()
	parts := strings.Split(tag, ":")
	if len(parts) != 2 {
		return false
	}

	db := database.GetInstance().GetClient()

	if !db.Migrator().HasColumn(parts[0], parts[1]) {
		return false
	}

	// Todo: Critical Reminder to check against SQL Injections.
	query := fmt.Sprintf("%s = ?", parts[1])
	var count int64
	err := db.Table(parts[0]).Where(query, value).Count(&count)

	if err.Error != nil {
		return false
	}

	return count > 0
}

func unique(fl validator.FieldLevel) bool {
	value := fl.Field().String()

	// Parse the tag to extract table name and column name
	tag := fl.Param()
	parts := strings.Split(tag, ":")
	if len(parts) != 2 {
		return false
	}

	db := database.GetInstance().GetClient()

	if !db.Migrator().HasColumn(parts[0], parts[1]) {
		return false
	}

	query := fmt.Sprintf("%s = ?", parts[1])
	var count int64
	err := db.Table(parts[0]).Where(query, value).Count(&count)

	if err.Error != nil {
		return false
	}

	return count == 0
}

func isRfc3339(fl validator.FieldLevel) bool {
	_, ok := fl.Field().Interface().(time.Time)
	return ok
}

func iranianCompanyNationalIdValidation(fl validator.FieldLevel) bool {
	value := fl.Field().String()

	// The national ID must be exactly 11 digits.
	if matched, _ := regexp.MatchString(`^\d{11}$`, value); !matched {
		return false
	}

	// It should not be sequential (e.g., "11111111111").
	if sequentialMatch, _ := regexp.MatchString(`^([0-9])\1{10}$`, value); sequentialMatch {
		return false
	}

	// Extract the control digit (last digit).
	controlDigit, err := strconv.Atoi(string(value[10]))
	if err != nil {
		return false
	}

	// Extract the tenth digit and add 2.
	tenthDigit, err := strconv.Atoi(string(value[9]))
	if err != nil {
		return false
	}
	adjustedTenthDigit := tenthDigit + 2

	// Calculate the weighted sum.
	weights := []int{29, 27, 23, 19, 17, 29, 27, 23, 19, 17}
	sum := 0
	for i := 0; i < 10; i++ {
		digit, err := strconv.Atoi(string(value[i]))
		if err != nil {
			return false
		}
		sum += (adjustedTenthDigit + digit) * weights[i]
	}

	// Calculate the control number.
	controlNumber := sum % 11
	if controlNumber == 10 {
		controlNumber = 0
	}

	// The ID is valid if the control number matches the last digit.
	return controlNumber == controlDigit
}
