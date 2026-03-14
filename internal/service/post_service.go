package service

import (
	"errors"
	"mini-api-golang/internal/dao"
	"mini-api-golang/internal/models"
)

// PostService provides business logic for post operations.
type PostService struct {
	postDAO         *dao.PostDAO
	followDAO       *dao.FollowDAO
	notificationDAO *dao.NotificationDAO
}

// NewPostService creates a new PostService.
func NewPostService(postDAO *dao.PostDAO, followDAO *dao.FollowDAO, notificationDAO *dao.NotificationDAO) *PostService {
	return &PostService{
		postDAO:         postDAO,
		followDAO:       followDAO,
		notificationDAO: notificationDAO,
	}
}

// Create validates and persists a new post, then notifies all followers of the author.
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
	saved, err := s.postDAO.GetByID(post.ID)
	if err != nil {
		return nil, err
	}

	// Notify all followers about the new post (best-effort, only for public posts)
	if visibility == models.VisibilityPublic {
		s.notifyFollowers(authorID, saved.ID)
	}

	return saved, nil
}

// notifyFollowers sends new_post notifications to all followers of the author
// using a single batch insert for efficiency.
func (s *PostService) notifyFollowers(authorID, postID uint) {
	followerIDs, err := s.followDAO.ListFollowerIDs(authorID)
	if err != nil || len(followerIDs) == 0 {
		return
	}
	notifications := make([]*models.Notification, len(followerIDs))
	for i, fID := range followerIDs {
		pid := postID
		notifications[i] = &models.Notification{
			RecipientID: fID,
			ActorID:     authorID,
			Type:        models.NotificationTypeNewPost,
			PostID:      &pid,
		}
	}
	_ = s.notificationDAO.BatchCreate(notifications)
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

