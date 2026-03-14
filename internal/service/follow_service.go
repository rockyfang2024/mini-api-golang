package service

import (
	"errors"
	"fmt"

	"go.uber.org/zap"
	"mini-api-golang/internal/dao"
	"mini-api-golang/internal/models"
	"mini-api-golang/pkg/logger"
)

// FollowService manages follow/unfollow operations and notifications.
type FollowService struct {
	followDAO       *dao.FollowDAO
	userDAO         *dao.UserDAO
	notificationDAO *dao.NotificationDAO
	settingsDAO     *dao.UserSettingsDAO
}

// NewFollowService creates a new FollowService.
func NewFollowService(followDAO *dao.FollowDAO, userDAO *dao.UserDAO, notificationDAO *dao.NotificationDAO, settingsDAO *dao.UserSettingsDAO) *FollowService {
	return &FollowService{
		followDAO:       followDAO,
		userDAO:         userDAO,
		notificationDAO: notificationDAO,
		settingsDAO:     settingsDAO,
	}
}

// Follow creates a follow relationship from followerID → followingID.
// Sends a notification to the user being followed.
func (s *FollowService) Follow(followerID, followingID uint) error {
	if followerID == followingID {
		return dao.ErrCannotFollowSelf
	}

	// Verify the target user exists
	if _, err := s.userDAO.GetByID(followingID); err != nil {
		return fmt.Errorf("user not found")
	}

	settings, err := ensureUserSettings(s.settingsDAO, followingID)
	if err != nil {
		return fmt.Errorf("failed to load user settings: %w", err)
	}
	if !settings.AllowFollow {
		return ErrFollowDisabled
	}

	// Check for duplicate
	exists, err := s.followDAO.Exists(followerID, followingID)
	if err != nil {
		return fmt.Errorf("database error: %w", err)
	}
	if exists {
		return dao.ErrAlreadyFollowed
	}

	follow := &models.Follow{
		FollowerID:  followerID,
		FollowingID: followingID,
	}
	if err := s.followDAO.Create(follow); err != nil {
		return fmt.Errorf("failed to create follow: %w", err)
	}

	// Notify the user being followed
	n := &models.Notification{
		RecipientID: followingID,
		ActorID:     followerID,
		Type:        models.NotificationTypeFollow,
	}
	if err := s.notificationDAO.Create(n); err != nil {
		logger.Error("failed to create follow notification", zap.Error(err), zap.Uint("follower_id", followerID), zap.Uint("following_id", followingID))
	}

	return nil
}

// ErrFollowDisabled indicates the user has disabled new followers.
var ErrFollowDisabled = errors.New("follow disabled")

// Unfollow removes a follow relationship from followerID → followingID.
func (s *FollowService) Unfollow(followerID, followingID uint) error {
	return s.followDAO.Delete(followerID, followingID)
}

// IsFollowing returns whether followerID is following followingID.
func (s *FollowService) IsFollowing(followerID, followingID uint) (bool, error) {
	return s.followDAO.Exists(followerID, followingID)
}

// ListFollowers returns paginated followers of the given user.
func (s *FollowService) ListFollowers(userID uint, page, pageSize int) ([]models.Follow, int64, error) {
	if page < 1 {
		page = 1
	}
	if pageSize < 1 || pageSize > 100 {
		pageSize = 20
	}
	return s.followDAO.ListFollowers(userID, page, pageSize)
}

// ListFollowing returns paginated users that the given user is following.
func (s *FollowService) ListFollowing(userID uint, page, pageSize int) ([]models.Follow, int64, error) {
	if page < 1 {
		page = 1
	}
	if pageSize < 1 || pageSize > 100 {
		pageSize = 20
	}
	return s.followDAO.ListFollowing(userID, page, pageSize)
}
