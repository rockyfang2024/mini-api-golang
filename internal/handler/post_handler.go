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
	postService   *service.PostService
	likeService   *service.LikeService
	repostService *service.RepostService
}

// NewPostHandler creates a new PostHandler.
func NewPostHandler(postService *service.PostService, likeService *service.LikeService, repostService *service.RepostService) *PostHandler {
	return &PostHandler{
		postService:   postService,
		likeService:   likeService,
		repostService: repostService,
	}
}

// PostResponse wraps a post with aggregated like/repost counts and viewer state.
type PostResponse struct {
	models.Post
	LikeCount   int64 `json:"like_count"`
	RepostCount int64 `json:"repost_count"`
	IsLiked     bool  `json:"is_liked"`
	IsReposted  bool  `json:"is_reposted"`
}

// enrichPost adds like/repost counts and viewer state to a post.
func (h *PostHandler) enrichPost(post models.Post, viewerID uint) PostResponse {
	likeCount, _ := h.likeService.LikeCount(post.ID)
	repostCount, _ := h.repostService.RepostCount(post.ID)

	var isLiked, isReposted bool
	if viewerID > 0 {
		isLiked, _ = h.likeService.IsLiked(viewerID, post.ID)
		isReposted, _ = h.repostService.IsReposted(viewerID, post.ID)
	}

	return PostResponse{
		Post:        post,
		LikeCount:   likeCount,
		RepostCount: repostCount,
		IsLiked:     isLiked,
		IsReposted:  isReposted,
	}
}

// enrichPosts enriches a slice of posts.
func (h *PostHandler) enrichPosts(posts []models.Post, viewerID uint) []PostResponse {
	result := make([]PostResponse, len(posts))
	for i, p := range posts {
		result[i] = h.enrichPost(p, viewerID)
	}
	return result
}

// createPostRequest is the expected body for POST /api/posts.
type createPostRequest struct {
	Content    string            `json:"content"    binding:"required"`
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

	utils.SuccessResponse(c, http.StatusCreated, "post created", h.enrichPost(*post, uint(userID)))
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

	utils.SuccessResponse(c, http.StatusOK, "ok", h.enrichPosts(posts, viewerID))
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

	utils.SuccessResponse(c, http.StatusOK, "ok", h.enrichPosts(posts, viewerID))
}

