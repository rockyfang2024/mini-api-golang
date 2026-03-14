package service

import (
	"mini-api-golang/internal/dao"
	"mini-api-golang/internal/models"
)

// UserSettingsService manages user settings.
type UserSettingsService struct {
	settingsDAO *dao.UserSettingsDAO
}

// NewUserSettingsService creates a new UserSettingsService.
func NewUserSettingsService(settingsDAO *dao.UserSettingsDAO) *UserSettingsService {
	return &UserSettingsService{settingsDAO: settingsDAO}
}

// GetOrCreate returns settings for a user, creating defaults if missing.
func (s *UserSettingsService) GetOrCreate(userID uint) (*models.UserSettings, error) {
	return ensureUserSettings(s.settingsDAO, userID)
}

// Update applies updated settings for a user.
func (s *UserSettingsService) Update(userID uint, updates models.UserSettings) (*models.UserSettings, error) {
	settings, err := ensureUserSettings(s.settingsDAO, userID)
	if err != nil {
		return nil, err
	}
	settings.AllowComments = updates.AllowComments
	settings.AllowFollow = updates.AllowFollow
	settings.OnlyFollowersCanView = updates.OnlyFollowersCanView
	settings.OnlyFollowingCanView = updates.OnlyFollowingCanView
	if err := s.settingsDAO.Update(settings); err != nil {
		return nil, err
	}
	return settings, nil
}
