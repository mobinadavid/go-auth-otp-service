package providers

import (
	"go-auth-otp-service/src/database"
	"go-auth-otp-service/src/repositories"
	"go-auth-otp-service/src/services"
)

func ProvideUserRepository(db *database.Database) *repositories.UserRepository {
	return &repositories.UserRepository{
		DatabaseHandler: db,
	}
}

func ProvideUserService(userRepository *repositories.UserRepository) *services.UserService {
	return &services.UserService{
		UserRepository: userRepository,
	}
}
