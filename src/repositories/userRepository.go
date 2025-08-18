package repositories

import (
	"fmt"
	"github.com/google/uuid"
	"go-auth-otp-service/src/database"
	"go-auth-otp-service/src/database/scopes"
	"go-auth-otp-service/src/models"
)

// IUserRepository interface defines the methods to interact with the User data store.
type IUserRepository interface {
	GetList(builder *scopes.BuilderModel) (*scopes.PaginateModel, error)
	GetByUuid(uuid *uuid.UUID) (*models.UserModel, error)
	GetByNationalIdentityCode(nationalIdentityCode string) (*models.UserModel, error)
	GetByMobile(mobile string) (*models.UserModel, error)
	Create(user *models.UserModel) (*models.UserModel, error)
	Update(user *models.UserModel) (*models.UserModel, error)
	Delete(user *models.UserModel) error
}

// UserRepository struct implements the UserRepository interface.
type UserRepository struct {
	DatabaseHandler *database.Database
}

// GetList retrieve all users
func (repository *UserRepository) GetList(builder *scopes.BuilderModel) (*scopes.PaginateModel, error) {
	var results []*models.UserModel

	// Get the database client
	db := repository.DatabaseHandler.GetClient().Model(results)

	// Apply pagination, filtering, and sorting using the BuilderModel
	db, err := builder.QueryBuilderScope(db)
	if err != nil {
		return nil, fmt.Errorf("user list retrieval failed: %s", err.Error())
	}

	// Create the PaginateModel and execute the query
	paginateModel, err := builder.CreatePaginateModel(db, &results)
	if err != nil {
		return nil, fmt.Errorf("user list retrieval failed: %s", err.Error())
	}
	return paginateModel, nil
}

// GetByUuid retrieve a user by uuid
func (repository *UserRepository) GetByUuid(uuid *uuid.UUID) (*models.UserModel, error) {
	var user models.UserModel
	result := repository.DatabaseHandler.GetClient().First(&user, "uuid = ?", uuid)
	if result.Error != nil {
		return nil, fmt.Errorf("user get by uuid failed: %s", result.Error.Error())
	}

	return &user, nil
}

// GetByNationalIdentityCode gets a user by national-identity-code.
func (repository *UserRepository) GetByNationalIdentityCode(nationalIdentityCode string) (*models.UserModel, error) {
	var user models.UserModel
	result := repository.DatabaseHandler.GetClient().First(&user, "national_identity_code = ?", nationalIdentityCode)
	if result.Error != nil {
		return nil, result.Error
	}

	return &user, nil
}

// GetByMobile gets a user by mobile .
func (repository *UserRepository) GetByMobile(mobile string) (*models.UserModel, error) {
	var user models.UserModel
	result := repository.DatabaseHandler.GetClient().First(&user, "mobile = ?", mobile)
	if result.Error != nil {
		return nil, result.Error
	}

	return &user, nil
}

// Create inserts a new User into the database
func (repository *UserRepository) Create(user *models.UserModel) (*models.UserModel, error) {
	result := repository.DatabaseHandler.GetClient().Create(&user)
	if result.Error != nil {
		return nil, fmt.Errorf("user creation failed: %s", result.Error.Error())
	}
	return user, nil
}

// Update update a user
func (repository *UserRepository) Update(user *models.UserModel) (*models.UserModel, error) {
	result := repository.DatabaseHandler.GetClient().Save(&user)
	if result.Error != nil {
		return nil, fmt.Errorf("user update failed: %s", result.Error.Error())
	}
	return user, nil
}

// Delete delete a user
func (repository *UserRepository) Delete(user *models.UserModel) error {
	result := repository.DatabaseHandler.GetClient().Delete(&user)
	if result.Error != nil {
		return fmt.Errorf("user delete failed: %s", result.Error.Error())
	}
	return nil
}
