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
	ErrDeactivatedAdmin       = errors.New("admin-is-not-active")
	ErrDeactivatedUser        = errors.New("user-is-not-active")
	ErrRegisterTimeOut        = errors.New("register-time-out")
	ErrRecoverPasswordTimeOut = errors.New("recover-password-time-out")
	ErrChangePasswordTimeOut  = errors.New("change-password-time-out")
	ErrLoginTimeOut           = errors.New("login-time-out")
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

// captcha
var (
	CaptchaServiceIsTemporaryDown = errors.New("captcha-service-is-temporary-down")
	CaptchaIsNotValid             = errors.New("captcha-is-not-valid")
)

// rate limiter
var (
	TooManyRequest = errors.New("too-many-request")
)

// sejam
var (
	SejamServiceIsTemporaryDown = errors.New("sejam-service-is-temporary-down")
	InvalidOtpSejam             = errors.New("sejam-otp-invalid")
	TooManyRequestForSejam      = errors.New("too-many-request-for-sejam")
	UserNotFoundInSejam         = errors.New("user-not-found-in-sejam")
)

// farabourse
var (
	ProjectNotExists = errors.New("project-not-exists")
	FaraBourseIsDown = errors.New("fara-bourse-is-down")
)

// Shareholder or Company Document
var (
	ShareholderPercentageExceeded = errors.New("shareholder-percentage-exceeded")
	CompanyDocumentDuplicateEntry = errors.New("company-document-duplicate-entry")
)

// file
var (
	InvalidFormData           = errors.New("invalid-form-data")
	CannotOpenFile            = errors.New("cannot-open-file")
	CDNServiceIsTemporaryDown = errors.New("cdn-service-is-temporary-down")
)

// project
var (
	ErrProjectNotEditable              = errors.New("project-not-editable")
	ErrPhaseDateGreaterThanLastPhase   = errors.New("phase-date-greater-than-last-phase")
	ErrPhaseDateIsNotBetweenLastPhases = errors.New("phase-date-is-not-between-last-phases")
	ErrIncompleteInformation           = errors.New("incomplete-information")
	ErrIncorrectInputInformation       = errors.New("incorrect-input-information")
	ErrDuplicateField                  = errors.New("duplicate-field")
	ErrNotFoundProject                 = errors.New("project-not-found")
	ErrStatusChange                    = errors.New("status-change")
)

// investment
var (
	ErrProfileNotCompleted       = errors.New("profile-is-not-completed")
	ErrContributionOutOfRange    = errors.New("contribution-out-range")
	ErrTotalContributionExceeded = errors.New("total-contribution-exceeded")
)

// ipg
var (
	IPGIsTemporaryDown = errors.New("ipg-is-temporary-down")
)

// wallet
var (
	MinimumAmountThreshold          = errors.New("amount-is-less-than-threshold")
	MaximumAmountThreshold          = errors.New("amount-is-great-than-threshold")
	ShebaIsNotActive                = errors.New("sheba-is-not-active")
	ShebaIsNotValid                 = errors.New("sheba-is-not-valid")
	NotEnoughMoneyForWallet         = errors.New("not-enough-money-for-wallet")
	WalletNotActive                 = errors.New("wallet-not-active")
	NotEnoughMoneyToStartInvestment = errors.New("not-enough-money-to-start-investment")
)
