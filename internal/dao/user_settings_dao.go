package dao

import (
	"errors"

	"gorm.io/gorm"
	"mini-api-golang/internal/models"
)

// UserSettingsDAO handles database operations for user settings.
type UserSettingsDAO struct {
	DB *gorm.DB
}

// NewUserSettingsDAO creates a new UserSettingsDAO.
func NewUserSettingsDAO(db *gorm.DB) *UserSettingsDAO {
	return &UserSettingsDAO{DB: db}
}

// GetByUserID retrieves settings for a user.
func (d *UserSettingsDAO) GetByUserID(userID uint) (*models.UserSettings, error) {
	var settings models.UserSettings
	err := d.DB.Where("user_id = ?", userID).First(&settings).Error
	if errors.Is(err, gorm.ErrRecordNotFound) {
		return nil, err
	}
	return &settings, err
}

// Create inserts new settings.
func (d *UserSettingsDAO) Create(settings *models.UserSettings) error {
	return d.DB.Create(settings).Error
}

// Update updates settings.
func (d *UserSettingsDAO) Update(settings *models.UserSettings) error {
	return d.DB.Save(settings).Error
}
