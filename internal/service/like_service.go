package service

import (
	"errors"
	"fmt"

	"go.uber.org/zap"
	"mini-api-golang/internal/dao"
	"mini-api-golang/internal/models"
	"mini-api-golang/pkg/logger"
)

// LikeService manages like/unlike operations and delegates notification creation.
type LikeService struct {
	likeDAO         *dao.LikeDAO
	postDAO         *dao.PostDAO
	notificationDAO *dao.NotificationDAO
}

// NewLikeService creates a new LikeService.
func NewLikeService(likeDAO *dao.LikeDAO, postDAO *dao.PostDAO, notificationDAO *dao.NotificationDAO) *LikeService {
	return &LikeService{
		likeDAO:         likeDAO,
		postDAO:         postDAO,
		notificationDAO: notificationDAO,
	}
}

// Like records a like for the given post by the given user.
// It also creates a notification for the post author if the liker is not the author.
func (s *LikeService) Like(userID, postID uint) error {
	// Check the post exists
	post, err := s.postDAO.GetByID(postID)
	if err != nil {
		return fmt.Errorf("post not found")
	}

	// Check for duplicate
	exists, err := s.likeDAO.Exists(userID, postID)
	if err != nil {
		return fmt.Errorf("database error: %w", err)
	}
	if exists {
		return dao.ErrAlreadyLiked
	}

	like := &models.Like{
		UserID: userID,
		PostID: postID,
	}
	if err := s.likeDAO.Create(like); err != nil {
		return fmt.Errorf("failed to create like: %w", err)
	}

	// Notify post author (skip self-like notifications)
	if post.AuthorID != userID {
		n := &models.Notification{
			RecipientID: post.AuthorID,
			ActorID:     userID,
			Type:        models.NotificationTypeLike,
			PostID:      &postID,
		}
		if err := s.notificationDAO.Create(n); err != nil {
			logger.Error("failed to create like notification", zap.Error(err), zap.Uint("post_id", postID), zap.Uint("user_id", userID))
		}
	}

	return nil
}

// Unlike removes a like for the given post by the given user.
func (s *LikeService) Unlike(userID, postID uint) error {
	err := s.likeDAO.Delete(userID, postID)
	if errors.Is(err, dao.ErrLikeNotFound) {
		return dao.ErrLikeNotFound
	}
	return err
}

// IsLiked returns whether the given user has liked the given post.
func (s *LikeService) IsLiked(userID, postID uint) (bool, error) {
	return s.likeDAO.Exists(userID, postID)
}

// LikeCount returns the number of likes for a post.
func (s *LikeService) LikeCount(postID uint) (int64, error) {
	return s.likeDAO.CountByPost(postID)
}
