package service

import (
	"errors"
	"mini-api-golang/internal/dao"
	"mini-api-golang/internal/models"
)

// PostService provides business logic for post operations.
type PostService struct {
	postDAO *dao.PostDAO
}

// NewPostService creates a new PostService.
func NewPostService(postDAO *dao.PostDAO) *PostService {
	return &PostService{postDAO: postDAO}
}

// Create validates and persists a new post.
func (s *PostService) Create(authorID uint, content string, visibility models.Visibility) (*models.Post, error) {
	if content == "" {
		return nil, errors.New("content cannot be empty")
	}
	if visibility != models.VisibilityPublic && visibility != models.VisibilityPrivate {
		return nil, errors.New("visibility must be 'public' or 'private'")
	}
	post := &models.Post{
		AuthorID:   authorID,
		Content:    content,
		Visibility: visibility,
	}
	if err := s.postDAO.Create(post); err != nil {
		return nil, err
	}
	// Reload to get preloaded Author
	return s.postDAO.GetByID(post.ID)
}

// ListHome returns posts for the home feed based on viewer authentication.
// Unauthenticated viewers (viewerID == 0) see only public posts.
func (s *PostService) ListHome(viewerID uint) ([]models.Post, error) {
	if viewerID == 0 {
		return s.postDAO.ListPublic()
	}
	return s.postDAO.ListForUser(viewerID)
}

// ListByUser returns posts for a given author as seen by the viewer.
// If viewerID equals authorID, all posts are returned; otherwise only public ones.
func (s *PostService) ListByUser(authorID, viewerID uint) ([]models.Post, error) {
	if authorID == viewerID {
		return s.postDAO.ListByAuthorAll(authorID)
	}
	return s.postDAO.ListByAuthorPublic(authorID)
}
