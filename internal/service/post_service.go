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
	settingsDAO     *dao.UserSettingsDAO
}

// NewPostService creates a new PostService.
func NewPostService(postDAO *dao.PostDAO, followDAO *dao.FollowDAO, notificationDAO *dao.NotificationDAO, settingsDAO *dao.UserSettingsDAO) *PostService {
	return &PostService{
		postDAO:         postDAO,
		followDAO:       followDAO,
		notificationDAO: notificationDAO,
		settingsDAO:     settingsDAO,
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
		posts, err := s.postDAO.ListPublic()
		if err != nil {
			return nil, err
		}
		return s.filterVisiblePosts(posts, viewerID)
	}
	posts, err := s.postDAO.ListForUser(viewerID)
	if err != nil {
		return nil, err
	}
	return s.filterVisiblePosts(posts, viewerID)
}

// ListByUser returns posts for a given author as seen by the viewer.
// If viewerID equals authorID, all posts are returned; otherwise only public ones.
func (s *PostService) ListByUser(authorID, viewerID uint) ([]models.Post, error) {
	if authorID == viewerID {
		return s.postDAO.ListByAuthorAll(authorID)
	}
	canView, err := canViewUserPosts(s.settingsDAO, s.followDAO, viewerID, authorID)
	if err != nil {
		return nil, err
	}
	if !canView {
		return nil, ErrPostsNotVisible
	}
	return s.postDAO.ListByAuthorPublic(authorID)
}

// Delete removes a post by ID.
func (s *PostService) Delete(postID uint) error {
	return s.postDAO.Delete(postID)
}

// CanViewPost determines whether a viewer can see a specific post.
func (s *PostService) CanViewPost(post *models.Post, viewerID uint) (bool, error) {
	if post.Visibility == models.VisibilityPrivate && post.AuthorID != viewerID {
		return false, nil
	}
	return canViewUserPosts(s.settingsDAO, s.followDAO, viewerID, post.AuthorID)
}

func (s *PostService) filterVisiblePosts(posts []models.Post, viewerID uint) ([]models.Post, error) {
	visible := make([]models.Post, 0, len(posts))
	for _, post := range posts {
		if post.Visibility == models.VisibilityPrivate && post.AuthorID != viewerID {
			continue
		}
		canView, err := canViewUserPosts(s.settingsDAO, s.followDAO, viewerID, post.AuthorID)
		if err != nil {
			return nil, err
		}
		if !canView {
			continue
		}
		visible = append(visible, post)
	}
	return visible, nil
}

// ErrPostsNotVisible indicates the viewer is not permitted to see the user's posts.
var ErrPostsNotVisible = errors.New("posts not visible")
