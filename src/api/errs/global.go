package errs

import (
	"errors"
)

// common
var (
	SomeThingWentWrong           = errors.New("something-went-wrong")
	InvalidUuid                  = errors.New("uuid-is-invalid")
	RecordNotFound               = errors.New("record-not-found")
	CantChangeSuperAdmin         = errors.New("you-can-not-change-super-admin-title")
	InvalidMethodForNotification = errors.New("not-supported-method-to-send-notification")
	BankNotExists                = errors.New("bank-not-exists")
	NotImplementedYet            = errors.New("feature-not-implemented-yet")
)

// authenticate
var (
	ErrAuthenticationFailed   = errors.New("auth-failed")
	RegisterFailed            = errors.New("register-failed")
	RecoverPasswordFailed     = errors.New("recover-password-failed")
	ChangePasswordFailed      = errors.New("change-password-failed")
	PasswordNotMatch          = errors.New("invalid-password-match")
	PasswordShouldBeNew       = errors.New("new-password-is-equal-to-current-password")
	ErrDeactivatedUser        = errors.New("user-is-not-active")
	ErrRegisterTimeOut        = errors.New("register-time-out")
	ErrRecoverPasswordTimeOut = errors.New("recover-password-time-out")
	ErrChangePasswordTimeOut  = errors.New("change-password-time-out")
	ErrLoginTimeOut           = errors.New("login-time-out")
	ErrUserAlreadyExists      = errors.New("user-already-exists")
	ErrAdminHasTwoFactorAuth  = errors.New("user-has-two-factor-auth")
	ErrSameUserActivation     = errors.New("user-has-the-same-activation")
	Invalid2FaCode            = errors.New("invalid-two-factor-code")
	NotActiveTwoFactor        = errors.New("two-factor-authentication-is-not-active")
	InvalidRecoveryCode       = errors.New("invalid-recovery-code")
	OTPIsNotValid             = errors.New("otp-is-not-valid")
)

// token
var (
	RefreshTokenMissing     = errors.New("refresh-token-is-missing")
	ErrInvalidRefreshToken  = errors.New("invalid-refresh-token")
	ErrTokenExpired         = errors.New("token-expired")
	ErrInvalidToken         = errors.New("invalid-token")
	ErrInvalidSigningMethod = errors.New("unexpected-signing-method")
)

// otp
var (
	ErrOTPRequired   = errors.New("auth-otp-sent")
	ErrOTPInvalid    = errors.New("auth-otp-invalid")
	ErrAuthOTPExists = errors.New("auth-otp-exists")
	FailedToSendOTP  = errors.New("failed-to-send-top")
)
