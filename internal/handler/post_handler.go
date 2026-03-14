package handler

import (
	"net/http"
	"strconv"

	"github.com/gin-gonic/gin"
	"mini-api-golang/internal/models"
	"mini-api-golang/internal/service"
	"mini-api-golang/internal/utils"
)

// PostHandler holds dependencies for post HTTP handlers.
type PostHandler struct {
	postService *service.PostService
}

// NewPostHandler creates a new PostHandler.
func NewPostHandler(postService *service.PostService) *PostHandler {
	return &PostHandler{postService: postService}
}

// createPostRequest is the expected body for POST /api/posts.
type createPostRequest struct {
	Content    string           `json:"content"    binding:"required"`
	Visibility models.Visibility `json:"visibility" binding:"required,oneof=public private"`
}

// CreatePost handles POST /api/posts — create a new post (requires auth).
func (h *PostHandler) CreatePost(c *gin.Context) {
	var req createPostRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		utils.ErrorResponse(c, http.StatusBadRequest, err.Error())
		return
	}

	userIDStr, _ := c.Get("user_id")
	userID, err := strconv.ParseUint(userIDStr.(string), 10, 64)
	if err != nil {
		utils.ErrorResponse(c, http.StatusUnauthorized, "invalid token")
		return
	}

	post, err := h.postService.Create(uint(userID), req.Content, req.Visibility)
	if err != nil {
		utils.ErrorResponse(c, http.StatusBadRequest, err.Error())
		return
	}

	utils.SuccessResponse(c, http.StatusCreated, "post created", post)
}

// ListPosts handles GET /api/posts — home feed (visibility filtered by auth state).
func (h *PostHandler) ListPosts(c *gin.Context) {
	var viewerID uint
	if userIDStr, exists := c.Get("user_id"); exists {
		if id, err := strconv.ParseUint(userIDStr.(string), 10, 64); err == nil {
			viewerID = uint(id)
		}
	}

	posts, err := h.postService.ListHome(viewerID)
	if err != nil {
		utils.ErrorResponse(c, http.StatusInternalServerError, "failed to fetch posts")
		return
	}

	utils.SuccessResponse(c, http.StatusOK, "ok", posts)
}

// ListUserPosts handles GET /api/users/:id/posts.
// Returns all posts if :id equals the logged-in user; otherwise only public posts.
func (h *PostHandler) ListUserPosts(c *gin.Context) {
	authorID, err := parseUintParam(c, "id")
	if err != nil {
		utils.ErrorResponse(c, http.StatusBadRequest, "invalid user id")
		return
	}

	var viewerID uint
	if userIDStr, exists := c.Get("user_id"); exists {
		if id, err := strconv.ParseUint(userIDStr.(string), 10, 64); err == nil {
			viewerID = uint(id)
		}
	}

	posts, err := h.postService.ListByUser(authorID, viewerID)
	if err != nil {
		utils.ErrorResponse(c, http.StatusInternalServerError, "failed to fetch posts")
		return
	}

	utils.SuccessResponse(c, http.StatusOK, "ok", posts)
}
