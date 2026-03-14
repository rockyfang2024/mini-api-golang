package service

import (
	"errors"

	"gorm.io/gorm"
	"mini-api-golang/internal/dao"
	"mini-api-golang/internal/models"
)

func ensureUserSettings(settingsDAO *dao.UserSettingsDAO, userID uint) (*models.UserSettings, error) {
	settings, err := settingsDAO.GetByUserID(userID)
	if err == nil {
		return settings, nil
	}
	if !errors.Is(err, gorm.ErrRecordNotFound) {
		return nil, err
	}
	settings = &models.UserSettings{
		UserID:        userID,
		AllowComments: true,
		AllowFollow:   true,
	}
	if err := settingsDAO.Create(settings); err != nil {
		return nil, err
	}
	return settings, nil
}
