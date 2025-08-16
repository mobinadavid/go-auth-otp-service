package authentication

import (
	"errors"
	"github.com/golang-jwt/jwt/v5"
	"github.com/google/uuid"
	"go-auth-otp-service/src/api/errs"
	"go-auth-otp-service/src/config"
	"strconv"
	"time"
)

// IJwtService defines the interface for JWT operations.
type IJwtService interface {
	Generate() (dto *JwtDTO, err error)
	Validate(tokenString string) (*Claims, error)
}

// JwtService implements the IJwtService interface.
type JwtService struct{}

type JwtDTO struct {
	Uuid                  uuid.UUID `json:"-"`
	AccessTokenString     string    `json:"access_token_string"`
	RefreshTokenString    string    `json:"refresh_token_string"`
	RefreshTokenExpiresAt time.Time `json:"refresh_token_expires_at"`
	AccessTokenExpiresAt  time.Time `json:"access_token_expires_at"`
}

// Claims defines the structure of the JWT claims.
type Claims struct {
	jwt.RegisteredClaims
}

// Generate generates an access token and a refresh token.
func (service *JwtService) Generate() (dto *JwtDTO, err error) {
	tokenUuid, _ := uuid.NewUUID()
	// Generate access token
	accessTokenLifetime, _ := strconv.Atoi(config.GetInstance().Get("JWT_ACCESS_TOKEN_LIFETIME"))
	accessTokenExpiresAt := time.Now().Add(time.Duration(accessTokenLifetime) * time.Second)
	accessTokenString, err := generateToken(tokenUuid, accessTokenExpiresAt)
	if err != nil {
		return nil, errs.SomeThingWentWrong
	}

	// Generate refresh token
	refreshTokenLifetime, _ := strconv.Atoi(config.GetInstance().Get("JWT_REFRESH_TOKEN_EXPIRATION"))
	refreshTokenExpiresAt := time.Now().Add(time.Duration(refreshTokenLifetime) * time.Second)
	refreshTokenString, err := generateToken(tokenUuid, refreshTokenExpiresAt)
	if err != nil {
		return dto, errs.SomeThingWentWrong
	}

	return &JwtDTO{
		Uuid:                  tokenUuid,
		AccessTokenString:     accessTokenString,
		RefreshTokenString:    refreshTokenString,
		RefreshTokenExpiresAt: refreshTokenExpiresAt,
		AccessTokenExpiresAt:  accessTokenExpiresAt,
	}, err
}

// Validate validates a token string and returns the claims if the token is valid.
func (service *JwtService) Validate(tokenString string) (*Claims, error) {
	claims := &Claims{}

	token, err := jwt.ParseWithClaims(tokenString, claims, func(token *jwt.Token) (interface{}, error) {
		if _, ok := token.Method.(*jwt.SigningMethodHMAC); !ok {
			return nil, errs.ErrInvalidSigningMethod
		}
		return []byte(config.GetInstance().Get("JWT_SECRET")), nil
	}, jwt.WithAudience(config.GetInstance().Get("APP_HOST")), jwt.WithIssuer(config.GetInstance().Get("APP_NAME")))

	// Handle parsing errors explicitly
	if err != nil {
		// Check if the error is related to token expiration
		if errors.Is(err, jwt.ErrTokenExpired) {
			return nil, errs.ErrTokenExpired
		}
		return nil, err
	}

	if !token.Valid {
		return nil, errs.ErrInvalidToken
	}

	return claims, nil
}

// generateToken creates a token with a specified expiration duration.
func generateToken(uuid uuid.UUID, expiresAt time.Time) (string, error) {
	claims := &Claims{
		RegisteredClaims: jwt.RegisteredClaims{
			ID:        uuid.String(),
			Issuer:    config.GetInstance().Get("APP_NAME"),
			Audience:  []string{config.GetInstance().Get("APP_HOST")},
			IssuedAt:  jwt.NewNumericDate(time.Now()),
			ExpiresAt: jwt.NewNumericDate(expiresAt),
			NotBefore: jwt.NewNumericDate(time.Now()),
		},
	}
	token := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)
	return token.SignedString([]byte(config.GetInstance().Get("JWT_SECRET")))
}
