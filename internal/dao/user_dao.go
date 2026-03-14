package dao

import (
    "gorm.io/gorm"
)

// User represents the user model.
type User struct {
    ID    uint   `gorm:"primaryKey"`
    Name  string `gorm:"size:100;not null"`
    Email string `gorm:"size:100;unique;not null"`
}

// UserDAO defines methods for interacting with User data.
type UserDAO struct {
    DB *gorm.DB
}

// NewUserDAO initializes a new UserDAO.
func NewUserDAO(db *gorm.DB) *UserDAO {
    return &UserDAO{DB: db}
}

// Create adds a new User to the database.
func (u *UserDAO) Create(user *User) error {
    return u.DB.Create(user).Error
}

// Get retrieves a User by ID.
func (u *UserDAO) Get(id uint) (*User, error) {
    var user User
    err := u.DB.First(&user, id).Error
    return &user, err
}

// Update modifies an existing User.
func (u *UserDAO) Update(user *User) error {
    return u.DB.Save(user).Error
}

// Delete removes a User from the database.
func (u *UserDAO) Delete(id uint) error {
    return u.DB.Delete(&User{}, id).Error
}

// List retrieves all Users from the database.
func (u *UserDAO) List() ([]User, error) {
    var users []User
    err := u.DB.Find(&users).Error
    return users, err
}