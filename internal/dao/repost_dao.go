package dao

import (
	"errors"

	"gorm.io/gorm"
	"mini-api-golang/internal/models"
)

// RepostDAO handles database operations for reposts.
type RepostDAO struct {
	DB *gorm.DB
}

// NewRepostDAO creates a new RepostDAO.
func NewRepostDAO(db *gorm.DB) *RepostDAO {
	return &RepostDAO{DB: db}
}

// Create inserts a repost record.
func (d *RepostDAO) Create(repost *models.Repost) error {
	return d.DB.Create(repost).Error
}

// Exists returns true if the user has already reposted the given post.
func (d *RepostDAO) Exists(userID, originalPostID uint) (bool, error) {
	var count int64
	err := d.DB.Model(&models.Repost{}).
		Where("user_id = ? AND original_post_id = ?", userID, originalPostID).
		Count(&count).Error
	return count > 0, err
}

// CountByPost returns the number of reposts for a given post.
func (d *RepostDAO) CountByPost(originalPostID uint) (int64, error) {
	var count int64
	err := d.DB.Model(&models.Repost{}).
		Where("original_post_id = ?", originalPostID).
		Count(&count).Error
	return count, err
}

// Sentinel errors.
var (
	ErrAlreadyReposted = errors.New("already reposted")
	ErrRepostNotFound  = errors.New("repost not found")
)
