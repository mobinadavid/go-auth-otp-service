package authentication

type RegisterRequest struct {
	//TODO: we can add more fields
	NationalIdentityCode string `json:"national_identity_code" validate:"omitempty,iranian-national-identity-code"`
	Mobile               string `json:"mobile" validate:"required,iranian-mobile"`
}

type VerifyRegisterOTP struct {
	RegisterKey string `json:"register_key" validate:"omitempty"`
	OTP         string `json:"otp" validate:"required"`
}
