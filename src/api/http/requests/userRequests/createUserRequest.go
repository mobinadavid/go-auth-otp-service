package userRequests

import "github.com/google/uuid"

// CreateRequest struct for validating incoming request data for creating a user
type CreateRequest struct {
	Uuid                 uuid.UUID `json:"uuid" validate:"required,uuid"`
	FirstName            string    `json:"first_name" validate:"required,max=255"`
	LastName             string    `json:"last_name" validate:"required,max=255"`
	NationalIdentityCode string    `json:"national_identity_code" `
	Mobile               string    `json:"mobile" validate:"required,e164"`
	Password             string    `json:"password" validate:"required,max=255,is-strong-password"`
}
