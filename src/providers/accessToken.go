package providers

import (
	authentication2 "go-auth-otp-service/src/api/http/controllers/authentication"
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

func ProvideUserAccessTokenController(accessTokenService *authentication.AccessTokenService) *authentication2.AccessTokenController {
	return &authentication2.AccessTokenController{
		AccessTokenService: accessTokenService,
	}
}
