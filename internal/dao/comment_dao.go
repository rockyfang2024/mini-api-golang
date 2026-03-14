package dao

import (
	"gorm.io/gorm"
	"mini-api-golang/internal/models"
)

// CommentDAO handles database operations for comments.
type CommentDAO struct {
	DB *gorm.DB
}

// NewCommentDAO creates a new CommentDAO.
func NewCommentDAO(db *gorm.DB) *CommentDAO {
	return &CommentDAO{DB: db}
}

// Create inserts a new comment.
func (d *CommentDAO) Create(comment *models.Comment) error {
	return d.DB.Create(comment).Error
}

// GetByID retrieves a comment by ID with author preloaded.
func (d *CommentDAO) GetByID(id uint) (*models.Comment, error) {
	var comment models.Comment
	err := d.DB.Preload("Author").First(&comment, id).Error
	return &comment, err
}

// ListByPost returns all comments for a post, oldest first.
func (d *CommentDAO) ListByPost(postID uint) ([]models.Comment, error) {
	var comments []models.Comment
	err := d.DB.Preload("Author").
		Where("post_id = ?", postID).
		Order("created_at ASC").
		Find(&comments).Error
	return comments, err
}
