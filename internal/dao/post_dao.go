package dao

import (
	"gorm.io/gorm"
	"mini-api-golang/internal/models"
)

// PostDAO provides data-access methods for Post records.
type PostDAO struct {
	DB *gorm.DB
}

// NewPostDAO creates a new PostDAO.
func NewPostDAO(db *gorm.DB) *PostDAO {
	return &PostDAO{DB: db}
}

// Create inserts a new post into the database.
func (d *PostDAO) Create(post *models.Post) error {
	return d.DB.Create(post).Error
}

// GetByID retrieves a post by its ID, preloading the author.
func (d *PostDAO) GetByID(id uint) (*models.Post, error) {
	var post models.Post
	err := d.DB.Preload("Author").Preload("Images").First(&post, id).Error
	return &post, err
}

// ListPublic returns all public posts, ordered newest-first, with authors preloaded.
func (d *PostDAO) ListPublic() ([]models.Post, error) {
	var posts []models.Post
	err := d.DB.Preload("Author").Preload("Images").
		Where("visibility = ?", models.VisibilityPublic).
		Order("created_at DESC").
		Find(&posts).Error
	return posts, err
}

// ListForUser returns:
//   - all public posts from everyone, plus
//   - all private posts belonging to viewerID
//
// ordered newest-first, with authors preloaded.
func (d *PostDAO) ListForUser(viewerID uint) ([]models.Post, error) {
	var posts []models.Post
	err := d.DB.Preload("Author").Preload("Images").
		Where("visibility = ? OR (visibility = ? AND author_id = ?)",
			models.VisibilityPublic, models.VisibilityPrivate, viewerID).
		Order("created_at DESC").
		Find(&posts).Error
	return posts, err
}

// ListByAuthorPublic returns all public posts belonging to authorID, newest-first.
func (d *PostDAO) ListByAuthorPublic(authorID uint) ([]models.Post, error) {
	var posts []models.Post
	err := d.DB.Preload("Author").Preload("Images").
		Where("author_id = ? AND visibility = ?", authorID, models.VisibilityPublic).
		Order("created_at DESC").
		Find(&posts).Error
	return posts, err
}

// ListByAuthorAll returns all posts (public + private) belonging to authorID, newest-first.
func (d *PostDAO) ListByAuthorAll(authorID uint) ([]models.Post, error) {
	var posts []models.Post
	err := d.DB.Preload("Author").Preload("Images").
		Where("author_id = ?", authorID).
		Order("created_at DESC").
		Find(&posts).Error
	return posts, err
}

// Delete removes a post by ID.
func (d *PostDAO) Delete(postID uint) error {
	return d.DB.Delete(&models.Post{}, postID).Error
}
