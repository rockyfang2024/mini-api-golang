package handler

import (
	"errors"
	"net/http"
	"strconv"

	"github.com/gin-gonic/gin"
	"mini-api-golang/internal/dao"
	"mini-api-golang/internal/service"
	"mini-api-golang/internal/utils"
)

// LikeHandler handles like/unlike requests.
type LikeHandler struct {
	likeService *service.LikeService
}

// NewLikeHandler creates a new LikeHandler.
func NewLikeHandler(likeService *service.LikeService) *LikeHandler {
	return &LikeHandler{likeService: likeService}
}

// LikePost handles POST /api/posts/:id/like — like a post.
func (h *LikeHandler) LikePost(c *gin.Context) {
	postID, err := parseUintParam(c, "id")
	if err != nil {
		utils.ErrorResponse(c, http.StatusBadRequest, "invalid post id")
		return
	}

	userIDStr, _ := c.Get("user_id")
	userID, err := strconv.ParseUint(userIDStr.(string), 10, 64)
	if err != nil {
		utils.ErrorResponse(c, http.StatusUnauthorized, "invalid token")
		return
	}

	if err := h.likeService.Like(uint(userID), postID); err != nil {
		if errors.Is(err, dao.ErrAlreadyLiked) {
			utils.ErrorResponse(c, http.StatusConflict, "already liked this post")
			return
		}
		utils.ErrorResponse(c, http.StatusBadRequest, err.Error())
		return
	}

	utils.SuccessResponse(c, http.StatusOK, "post liked", nil)
}

// UnlikePost handles DELETE /api/posts/:id/like — remove a like from a post.
func (h *LikeHandler) UnlikePost(c *gin.Context) {
	postID, err := parseUintParam(c, "id")
	if err != nil {
		utils.ErrorResponse(c, http.StatusBadRequest, "invalid post id")
		return
	}

	userIDStr, _ := c.Get("user_id")
	userID, err := strconv.ParseUint(userIDStr.(string), 10, 64)
	if err != nil {
		utils.ErrorResponse(c, http.StatusUnauthorized, "invalid token")
		return
	}

	if err := h.likeService.Unlike(uint(userID), postID); err != nil {
		if errors.Is(err, dao.ErrLikeNotFound) {
			utils.ErrorResponse(c, http.StatusNotFound, "like not found")
			return
		}
		utils.ErrorResponse(c, http.StatusInternalServerError, "failed to remove like")
		return
	}

	utils.SuccessResponse(c, http.StatusOK, "like removed", nil)
}
