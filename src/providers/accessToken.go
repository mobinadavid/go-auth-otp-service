package providers

import (
	"go-auth-otp-service/src/database"
	"go-auth-otp-service/src/repositories"
	"go-auth-otp-service/src/services/authentication"
)

func ProvideAccessTokenService(accessTokenRepository *repositories.AccessTokenRepository,
	jwtService *authentication.JwtService,
	UserRepository *repositories.UserRepository) *authentication.AccessTokenService {
	return &authentication.AccessTokenService{
		AccessTokenRepository: accessTokenRepository,
		JwtService:            jwtService,
		UserRepository:        UserRepository,
	}
}

func ProvideAccessTokenRepository(db *database.Database) *repositories.AccessTokenRepository {
	return &repositories.AccessTokenRepository{
		DatabaseHandler: db,
	}
}
