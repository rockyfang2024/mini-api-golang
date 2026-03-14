package service

import (
	"errors"
	"mini-api-golang/internal/dao"
	"mini-api-golang/internal/models"
	"mini-api-golang/internal/utils"
)

// UserService provides business logic for user operations.
type UserService struct {
	userDAO *dao.UserDAO
}

// NewUserService creates a new UserService with the given DAO.
func NewUserService(userDAO *dao.UserDAO) *UserService {
	return &UserService{userDAO: userDAO}
}

// Register creates a new user with a hashed password.
func (s *UserService) Register(username, email, password string) (*models.User, error) {
	// Check for duplicate username
	if _, err := s.userDAO.GetByUsername(username); err == nil {
		return nil, errors.New("username already exists")
	}

	hash, err := utils.HashPassword(password)
	if err != nil {
		return nil, err
	}

	user := &models.User{
		Username:     username,
		Email:        email,
		PasswordHash: hash,
	}

	if err := s.userDAO.Create(user); err != nil {
		return nil, err
	}

	return user, nil
}

// Login validates credentials and returns the user on success.
func (s *UserService) Login(username, password string) (*models.User, error) {
	user, err := s.userDAO.GetByUsername(username)
	if err != nil {
		return nil, errors.New("invalid credentials")
	}

	if !utils.CheckPassword(password, user.PasswordHash) {
		return nil, errors.New("invalid credentials")
	}

	return user, nil
}

// GetByID retrieves a user by ID.
func (s *UserService) GetByID(id uint) (*models.User, error) {
	return s.userDAO.GetByID(id)
}

// Update modifies an existing user.
func (s *UserService) Update(user *models.User) error {
	return s.userDAO.Update(user)
}

// Delete removes a user by ID.
func (s *UserService) Delete(id uint) error {
	return s.userDAO.Delete(id)
}

// List retrieves all users.
func (s *UserService) List() ([]models.User, error) {
	return s.userDAO.List()
}