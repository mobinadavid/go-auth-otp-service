package authentication

import (
	"errors"
	"github.com/google/uuid"
	"go-auth-otp-service/src/api/errs"
	"go-auth-otp-service/src/database/scopes"
	"go-auth-otp-service/src/hash"
	"go-auth-otp-service/src/models"
	"go-auth-otp-service/src/repositories"
	"time"
)

type TokenType string
type OwnerType string

const (
	AccessToken  TokenType = "access_token"
	RefreshToken TokenType = "refresh_token"
)

type IAccessTokenService interface {
	GetList(builder *scopes.BuilderModel) (*scopes.PaginateModel, error)
	GetActiveTokens(builder *scopes.BuilderModel) (*scopes.PaginateModel, error)
	GetByUuid(accessTokenUuid *uuid.UUID) (*models.AccessTokenModel, error)
	Create(owner interface{}, dto *JwtDTO, ip, userAgent string) (*models.AccessTokenModel, error)
	UpdateLastUsedAt(accessToken *models.AccessTokenModel) (*models.AccessTokenModel, error)
	RefreshAccessTokens(refreshToken, ownerType string) (*JwtDTO, error)
	Validate(tokenString string, tokenType TokenType, ownerType string) (*models.AccessTokenModel, error)
	RevokeTokens(ownerID uint, ownerType string) error
	RevokeTokenByUuid(accessTokenUuid *uuid.UUID, ownerID uint, ownerType string) error
}

type AccessTokenService struct {
	AccessTokenRepository repositories.IAccessTokenRepository
	JwtService            IJwtService
	UserRepository        repositories.IUserRepository
}

func (service *AccessTokenService) GetList(builder *scopes.BuilderModel) (*scopes.PaginateModel, error) {
	res, err := service.AccessTokenRepository.GetList(builder)
	if err != nil {
		return nil, errs.SomeThingWentWrong
	}

	return res, nil
}

func (service *AccessTokenService) GetActiveTokens(builder *scopes.BuilderModel) (*scopes.PaginateModel, error) {
	res, err := service.AccessTokenRepository.GetActiveTokens(builder)
	if err != nil {
		return nil, errs.SomeThingWentWrong
	}

	return res, nil
}

func (service *AccessTokenService) Create(owner interface{}, dto *JwtDTO, ip, userAgent string) (*models.AccessTokenModel, error) {
	accessToken := &models.AccessTokenModel{
		Uuid:                  dto.Uuid,
		AccessToken:           []byte(dto.AccessTokenString),
		AccessTokenExpiresAt:  dto.AccessTokenExpiresAt,
		RefreshToken:          []byte(dto.RefreshTokenString),
		RefreshTokenExpiresAt: dto.RefreshTokenExpiresAt,
		IP:                    ip,
		UserAgent:             userAgent,
	}
	switch owner := owner.(type) {
	case *models.UserModel:
		accessToken.OwnerID = owner.ID
		accessToken.OwnerType = "user"
	default:
		return nil, errors.New("unsupported owner type")
	}

	atOrm, err := service.AccessTokenRepository.Create(accessToken)
	if err != nil {
		return nil, err
	}

	return atOrm, nil
}

func (service *AccessTokenService) GetByUuid(accessTokenUuid *uuid.UUID) (*models.AccessTokenModel, error) {
	accessToken, err := service.AccessTokenRepository.GetByUuid(accessTokenUuid)
	if err != nil {
		return nil, errs.SomeThingWentWrong
	}

	return accessToken, nil
}

func (service *AccessTokenService) UpdateLastUsedAt(accessToken *models.AccessTokenModel) (*models.AccessTokenModel, error) {
	res, err := service.AccessTokenRepository.UpdateLastUsedAt(accessToken, time.Now())
	if err != nil {
		return nil, err
	}

	return res, nil
}

func (service *AccessTokenService) RefreshAccessTokens(refreshToken, ownerType string) (*JwtDTO, error) {
	//validate token
	token, err := service.Validate(refreshToken, RefreshToken, ownerType)
	if err != nil {
		return nil, errs.ErrInvalidRefreshToken
	}

	// generate new jwt
	jwtDto, err := service.JwtService.Generate()
	if err != nil {
		return nil, errs.SomeThingWentWrong
	}

	// update user tokens
	_, err = service.AccessTokenRepository.RefreshAccessTokens(
		token,
		jwtDto.Uuid,
		[]byte(jwtDto.AccessTokenString),
		[]byte(jwtDto.RefreshTokenString),
		jwtDto.AccessTokenExpiresAt,
		jwtDto.RefreshTokenExpiresAt,
	)
	if err != nil {
		return nil, errs.SomeThingWentWrong
	}

	return jwtDto, nil
}

func (service *AccessTokenService) Validate(tokenString string, tokenType TokenType, ownerType string) (*models.AccessTokenModel, error) {
	// Validate the extracted JWT token and retrieve the user claims.
	userClaimed, err := service.JwtService.Validate(tokenString)
	if err != nil {
		return nil, err
	}

	// Additionally, validate the token against the database and check for its expiry.
	claimedUuid, err := uuid.Parse(userClaimed.ID)
	if err != nil {
		return nil, err
	}

	// Retrieve the access token from the database using the parsed UUID.
	token, err := service.AccessTokenRepository.GetByUuid(&claimedUuid)
	if err != nil {
		return nil, errs.RecordNotFound
	}

	if token.OwnerType != ownerType {
		return nil, errs.ErrAuthenticationFailed
	}

	// Verify the hash of the stored token against the provided token to ensure they match.
	var storedHash []byte
	var tokenExpiresAt time.Time

	switch tokenType {
	case AccessToken:
		storedHash = token.AccessToken
		tokenExpiresAt = token.AccessTokenExpiresAt

	case RefreshToken:
		storedHash = token.RefreshToken
		tokenExpiresAt = token.RefreshTokenExpiresAt
	}

	hashCheck, err := hash.VerifyStoredHash(storedHash, tokenString)
	if err != nil || !hashCheck {
		return nil, errs.ErrAuthenticationFailed
	}

	// Check if the token has expired by comparing its expiry timestamp against the current time.
	if tokenExpiresAt.Before(time.Now()) {
		return nil, errs.ErrTokenExpired
	}
	return token, nil
}

func (service *AccessTokenService) RevokeTokens(ownerID uint, ownerType string) error {
	// get list of tokens
	accessTokens, err := service.AccessTokenRepository.GetAll(ownerID, ownerType)
	if err != nil {
		return errs.SomeThingWentWrong
	}

	// delete the tokens
	err = service.AccessTokenRepository.DeleteMany(accessTokens)
	if err != nil {
		return errs.SomeThingWentWrong
	}
	return nil
}

func (service *AccessTokenService) RevokeTokenByUuid(accessTokenUuid *uuid.UUID, ownerID uint, ownerType string) error {
	accessToken, err := service.AccessTokenRepository.GetByUuid(accessTokenUuid)
	if err != nil {
		return errs.RecordNotFound
	}

	if accessToken.OwnerID != ownerID || accessToken.OwnerType != ownerType {
		return errs.RecordNotFound
	}

	err = service.AccessTokenRepository.Delete(accessToken)
	if err != nil {
		return errs.SomeThingWentWrong
	}
	return nil
}
