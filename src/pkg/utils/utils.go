package utils

import (
	"crypto/rand"
	"encoding/base64"
	"fmt"
	"math/big"
	mRand "math/rand"
	"reflect"
	"strings"
)

// GenerateSalt generates a new salt of the given length.
func GenerateSalt(length uint32) ([]byte, error) {
	salt := make([]byte, length)
	_, err := rand.Read(salt)
	if err != nil {
		return nil, err
	}
	return salt, nil
}

func GenerateRandomString(n int) string {
	letterRunes := []rune("abcdefghijklmnopqrstuvwxyzABCDEFGHIJKLMNOPQRSTUVWXYZ")
	b := make([]rune, n)
	for i := range b {
		b[i] = letterRunes[mRand.Intn(len(letterRunes))]
	}
	return string(b)
}

func GenerateRandomCodes(number int) ([]string, error) {
	const charset = "ABCDEFGHIJKLMNOPQRSTUVWXYZabcdefghijklmnopqrstuvwxyz0123456789"
	const codeLength = 5
	var randomCodes []string

	for i := 0; i < number; i++ {
		var part1, part2 string

		for j := 0; j < codeLength; j++ {
			randomIndex1, _ := rand.Int(rand.Reader, big.NewInt(int64(len(charset))))
			randomIndex2, _ := rand.Int(rand.Reader, big.NewInt(int64(len(charset))))
			part1 += string(charset[randomIndex1.Int64()])
			part2 += string(charset[randomIndex2.Int64()])
		}

		randomCode := fmt.Sprintf("%s-%s", part1, part2)
		randomCodes = append(randomCodes, randomCode)
	}

	return randomCodes, nil
}

// IsLetter function to check if a character is a letter
func IsLetter(c byte) bool {
	return ('A' <= c && c <= 'Z') || ('a' <= c && c <= 'z')
}

// IsDigit function to check if a character is a digit
func IsDigit(c byte) bool {
	return '0' <= c && c <= '9'
}

// IsUpdateRequestEmpty checks if the update request has at least one non-nil value
func IsUpdateRequestEmpty(req interface{}) bool {
	reqValue := reflect.ValueOf(req)

	// Iterate through the fields of the update request
	for i := 0; i < reqValue.NumField(); i++ {
		fieldValue := reqValue.Field(i)

		// Check if the field is nil (i.e., not updated)
		if fieldValue.IsValid() && !fieldValue.IsNil() {
			return false
		}
	}

	return true
}

func Base64ToBigInt(b64 string) *big.Int {
	data, _ := base64.StdEncoding.DecodeString(b64)
	return new(big.Int).SetBytes(data)
}

// GetStructFieldNames returns the field names of a struct.
func GetStructFieldNames(model interface{}) map[string]bool {
	fieldNames := make(map[string]bool)
	val := reflect.ValueOf(model)
	for i := 0; i < val.Type().NumField(); i++ {
		field := val.Type().Field(i)
		jsonTag := field.Tag.Get("json")
		if jsonTag != "" && jsonTag != "-" {
			jsonField := strings.Split(jsonTag, ",")[0]
			fieldNames[jsonField] = true
		}
	}
	return fieldNames
}
