package dao

import (
	"gorm.io/gorm"
	"mini-api-golang/internal/models"
)

// PostImageDAO handles database operations for post images.
type PostImageDAO struct {
	DB *gorm.DB
}

// NewPostImageDAO creates a new PostImageDAO.
func NewPostImageDAO(db *gorm.DB) *PostImageDAO {
	return &PostImageDAO{DB: db}
}

// CreateBatch inserts multiple post images.
func (d *PostImageDAO) CreateBatch(images []models.PostImage) error {
	if len(images) == 0 {
		return nil
	}
	return d.DB.Create(&images).Error
}

// DeleteByPostID removes images for a post.
func (d *PostImageDAO) DeleteByPostID(postID uint) error {
	return d.DB.Where("post_id = ?", postID).Delete(&models.PostImage{}).Error
}
