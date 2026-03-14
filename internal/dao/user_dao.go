package dao

import (
	"gorm.io/gorm"
	"mini-api-golang/internal/models"
)

// UserDAO provides methods to access User data.
type UserDAO struct {
	DB *gorm.DB
}

// NewUserDAO initializes a new UserDAO.
func NewUserDAO(db *gorm.DB) *UserDAO {
	return &UserDAO{DB: db}
}

// Create adds a new User to the database.
func (u *UserDAO) Create(user *models.User) error {
	return u.DB.Create(user).Error
}

// GetByID retrieves a User by ID.
func (u *UserDAO) GetByID(id uint) (*models.User, error) {
	var user models.User
	err := u.DB.First(&user, id).Error
	return &user, err
}

// GetByUsername retrieves a User by username.
func (u *UserDAO) GetByUsername(username string) (*models.User, error) {
	var user models.User
	err := u.DB.Where("username = ?", username).First(&user).Error
	return &user, err
}

// Update modifies an existing User.
func (u *UserDAO) Update(user *models.User) error {
	return u.DB.Save(user).Error
}

// Delete removes a User from the database by ID.
func (u *UserDAO) Delete(id uint) error {
	return u.DB.Delete(&models.User{}, id).Error
}

// List retrieves all Users from the database.
func (u *UserDAO) List() ([]models.User, error) {
	var users []models.User
	err := u.DB.Find(&users).Error
	return users, err
}