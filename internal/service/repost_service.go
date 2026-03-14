package service

import (
	"fmt"

	"go.uber.org/zap"
	"mini-api-golang/internal/dao"
	"mini-api-golang/internal/models"
	"mini-api-golang/pkg/logger"
)

// RepostService manages repost operations and notification creation.
type RepostService struct {
	repostDAO       *dao.RepostDAO
	postDAO         *dao.PostDAO
	notificationDAO *dao.NotificationDAO
}

// NewRepostService creates a new RepostService.
func NewRepostService(repostDAO *dao.RepostDAO, postDAO *dao.PostDAO, notificationDAO *dao.NotificationDAO) *RepostService {
	return &RepostService{
		repostDAO:       repostDAO,
		postDAO:         postDAO,
		notificationDAO: notificationDAO,
	}
}

// Repost records a repost for the given post by the given user.
// A user may only repost a given post once; subsequent calls return ErrAlreadyReposted.
// It also creates a notification for the post author.
func (s *RepostService) Repost(userID, originalPostID uint) error {
	// Check the original post exists
	post, err := s.postDAO.GetByID(originalPostID)
	if err != nil {
		return fmt.Errorf("post not found")
	}

	// Check for duplicate
	exists, err := s.repostDAO.Exists(userID, originalPostID)
	if err != nil {
		return fmt.Errorf("database error: %w", err)
	}
	if exists {
		return dao.ErrAlreadyReposted
	}

	repost := &models.Repost{
		UserID:         userID,
		OriginalPostID: originalPostID,
	}
	if err := s.repostDAO.Create(repost); err != nil {
		return fmt.Errorf("failed to create repost: %w", err)
	}

	// Notify post author (skip self-repost notifications)
	if post.AuthorID != userID {
		n := &models.Notification{
			RecipientID: post.AuthorID,
			ActorID:     userID,
			Type:        models.NotificationTypeRepost,
			PostID:      &originalPostID,
		}
		if err := s.notificationDAO.Create(n); err != nil {
			logger.Error("failed to create repost notification", zap.Error(err), zap.Uint("post_id", originalPostID), zap.Uint("user_id", userID))
		}
	}

	return nil
}

// IsReposted returns whether the given user has reposted the given post.
func (s *RepostService) IsReposted(userID, originalPostID uint) (bool, error) {
	return s.repostDAO.Exists(userID, originalPostID)
}

// RepostCount returns the number of reposts for a post.
func (s *RepostService) RepostCount(originalPostID uint) (int64, error) {
	return s.repostDAO.CountByPost(originalPostID)
}
