package authentication

type RegisterRequest struct {
	NationalIdentityCode string `json:"national_identity_code" validate:"omitempty,iranian-national-identity-code"`
	Mobile               string `json:"mobile" validate:"required,e164"`
}
