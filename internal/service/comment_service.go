package service

import (
	"errors"
	"strings"

	"mini-api-golang/internal/dao"
	"mini-api-golang/internal/models"
)

// CommentService provides business logic for comments.
type CommentService struct {
	commentDAO  *dao.CommentDAO
	postDAO     *dao.PostDAO
	settingsDAO *dao.UserSettingsDAO
	followDAO   *dao.FollowDAO
}

// NewCommentService creates a new CommentService.
func NewCommentService(commentDAO *dao.CommentDAO, postDAO *dao.PostDAO, settingsDAO *dao.UserSettingsDAO, followDAO *dao.FollowDAO) *CommentService {
	return &CommentService{
		commentDAO:  commentDAO,
		postDAO:     postDAO,
		settingsDAO: settingsDAO,
		followDAO:   followDAO,
	}
}

// CreateComment adds a comment to a post.
func (s *CommentService) CreateComment(postID, authorID uint, content string) (*models.Comment, error) {
	content = strings.TrimSpace(content)
	if content == "" {
		return nil, errors.New("content cannot be empty")
	}

	post, err := s.postDAO.GetByID(postID)
	if err != nil {
		return nil, err
	}

	if err := s.ensureCanComment(post, authorID); err != nil {
		return nil, err
	}

	comment := &models.Comment{
		PostID:   postID,
		AuthorID: authorID,
		Content:  content,
	}
	if err := s.commentDAO.Create(comment); err != nil {
		return nil, err
	}

	return s.commentDAO.GetByID(comment.ID)
}

// ReplyToComment adds a reply to an existing comment.
func (s *CommentService) ReplyToComment(parentID, authorID uint, content string) (*models.Comment, error) {
	content = strings.TrimSpace(content)
	if content == "" {
		return nil, errors.New("content cannot be empty")
	}

	parent, err := s.commentDAO.GetByID(parentID)
	if err != nil {
		return nil, err
	}

	post, err := s.postDAO.GetByID(parent.PostID)
	if err != nil {
		return nil, err
	}

	if err := s.ensureCanComment(post, authorID); err != nil {
		return nil, err
	}

	comment := &models.Comment{
		PostID:          parent.PostID,
		ParentCommentID: &parent.ID,
		AuthorID:        authorID,
		Content:         content,
	}
	if err := s.commentDAO.Create(comment); err != nil {
		return nil, err
	}

	return s.commentDAO.GetByID(comment.ID)
}

// ListByPost returns all comments for a post.
func (s *CommentService) ListByPost(postID uint, viewerID uint) ([]models.Comment, error) {
	post, err := s.postDAO.GetByID(postID)
	if err != nil {
		return nil, err
	}

	if err := s.ensureCanViewPost(post, viewerID); err != nil {
		return nil, err
	}

	return s.commentDAO.ListByPost(postID)
}

func (s *CommentService) ensureCanViewPost(post *models.Post, viewerID uint) error {
	if post.Visibility == models.VisibilityPrivate && post.AuthorID != viewerID {
		return ErrPostsNotVisible
	}
	canView, err := canViewUserPosts(s.settingsDAO, s.followDAO, viewerID, post.AuthorID)
	if err != nil {
		return err
	}
	if !canView {
		return ErrPostsNotVisible
	}
	return nil
}

func (s *CommentService) ensureCanComment(post *models.Post, viewerID uint) error {
	if err := s.ensureCanViewPost(post, viewerID); err != nil {
		return err
	}
	settings, err := ensureUserSettings(s.settingsDAO, post.AuthorID)
	if err != nil {
		return err
	}
	if !settings.AllowComments {
		return ErrCommentsDisabled
	}
	return nil
}

// ErrCommentsDisabled indicates comments are disabled for the post's author.
var ErrCommentsDisabled = errors.New("comments are disabled")
