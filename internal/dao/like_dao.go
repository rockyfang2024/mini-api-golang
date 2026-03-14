package dao

import (
	"errors"

	"gorm.io/gorm"
	"mini-api-golang/internal/models"
)

// LikeDAO handles database operations for likes.
type LikeDAO struct {
	DB *gorm.DB
}

// NewLikeDAO creates a new LikeDAO.
func NewLikeDAO(db *gorm.DB) *LikeDAO {
	return &LikeDAO{DB: db}
}

// Create inserts a like record. Returns ErrDuplicateLike if the user already liked the post.
func (d *LikeDAO) Create(like *models.Like) error {
	result := d.DB.Create(like)
	return result.Error
}

// Delete removes a like by userID and postID. Returns ErrLikeNotFound if not found.
func (d *LikeDAO) Delete(userID, postID uint) error {
	result := d.DB.Where("user_id = ? AND post_id = ?", userID, postID).Delete(&models.Like{})
	if result.Error != nil {
		return result.Error
	}
	if result.RowsAffected == 0 {
		return ErrLikeNotFound
	}
	return nil
}

// Exists returns true if the user has liked the given post.
func (d *LikeDAO) Exists(userID, postID uint) (bool, error) {
	var count int64
	err := d.DB.Model(&models.Like{}).
		Where("user_id = ? AND post_id = ?", userID, postID).
		Count(&count).Error
	return count > 0, err
}

// CountByPost returns the number of likes for a given post.
func (d *LikeDAO) CountByPost(postID uint) (int64, error) {
	var count int64
	err := d.DB.Model(&models.Like{}).Where("post_id = ?", postID).Count(&count).Error
	return count, err
}

// Sentinel errors.
var (
	ErrLikeNotFound  = errors.New("like not found")
	ErrAlreadyLiked  = errors.New("already liked")
)
