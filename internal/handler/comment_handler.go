package handler

import (
	"errors"
	"net/http"
	"strconv"

	"github.com/gin-gonic/gin"
	"gorm.io/gorm"
	"mini-api-golang/internal/service"
	"mini-api-golang/internal/utils"
)

// CommentHandler handles comment endpoints.
type CommentHandler struct {
	commentService *service.CommentService
}

// NewCommentHandler creates a new CommentHandler.
func NewCommentHandler(commentService *service.CommentService) *CommentHandler {
	return &CommentHandler{commentService: commentService}
}

type commentRequest struct {
	Content string `json:"content" binding:"required"`
}

// CreateComment handles POST /api/posts/:id/comments.
func (h *CommentHandler) CreateComment(c *gin.Context) {
	postID, err := parseUintParam(c, "id")
	if err != nil {
		utils.ErrorResponse(c, http.StatusBadRequest, "invalid post id")
		return
	}

	userID, ok := getUserID(c)
	if !ok {
		utils.ErrorResponse(c, http.StatusUnauthorized, "invalid token")
		return
	}

	var req commentRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		utils.ErrorResponse(c, http.StatusBadRequest, err.Error())
		return
	}

	comment, err := h.commentService.CreateComment(postID, userID, req.Content)
	if err != nil {
		h.handleCommentError(c, err)
		return
	}

	utils.SuccessResponse(c, http.StatusCreated, "comment created", comment)
}

// ReplyComment handles POST /api/comments/:id/replies.
func (h *CommentHandler) ReplyComment(c *gin.Context) {
	commentID, err := parseUintParam(c, "id")
	if err != nil {
		utils.ErrorResponse(c, http.StatusBadRequest, "invalid comment id")
		return
	}

	userID, ok := getUserID(c)
	if !ok {
		utils.ErrorResponse(c, http.StatusUnauthorized, "invalid token")
		return
	}

	var req commentRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		utils.ErrorResponse(c, http.StatusBadRequest, err.Error())
		return
	}

	comment, err := h.commentService.ReplyToComment(commentID, userID, req.Content)
	if err != nil {
		h.handleCommentError(c, err)
		return
	}

	utils.SuccessResponse(c, http.StatusCreated, "comment created", comment)
}

// ListComments handles GET /api/posts/:id/comments.
func (h *CommentHandler) ListComments(c *gin.Context) {
	postID, err := parseUintParam(c, "id")
	if err != nil {
		utils.ErrorResponse(c, http.StatusBadRequest, "invalid post id")
		return
	}

	var viewerID uint
	if userIDStr, exists := c.Get("user_id"); exists {
		if id, err := strconv.ParseUint(userIDStr.(string), 10, 64); err == nil {
			viewerID = uint(id)
		}
	}

	comments, err := h.commentService.ListByPost(postID, viewerID)
	if err != nil {
		h.handleCommentError(c, err)
		return
	}

	utils.SuccessResponse(c, http.StatusOK, "ok", comments)
}

func (h *CommentHandler) handleCommentError(c *gin.Context, err error) {
	switch {
	case errors.Is(err, service.ErrPostsNotVisible):
		utils.ErrorResponse(c, http.StatusForbidden, "posts not visible")
	case errors.Is(err, service.ErrCommentsDisabled):
		utils.ErrorResponse(c, http.StatusForbidden, "comments are disabled")
	case errors.Is(err, gorm.ErrRecordNotFound):
		utils.ErrorResponse(c, http.StatusNotFound, "resource not found")
	default:
		utils.ErrorResponse(c, http.StatusBadRequest, err.Error())
	}
}

func getUserID(c *gin.Context) (uint, bool) {
	userIDStr, exists := c.Get("user_id")
	if !exists {
		return 0, false
	}
	userID, err := strconv.ParseUint(userIDStr.(string), 10, 64)
	if err != nil {
		return 0, false
	}
	return uint(userID), true
}
