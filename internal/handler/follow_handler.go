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

// FollowHandler handles follow/unfollow and follower/following list requests.
type FollowHandler struct {
	followService *service.FollowService
}

// NewFollowHandler creates a new FollowHandler.
func NewFollowHandler(followService *service.FollowService) *FollowHandler {
	return &FollowHandler{followService: followService}
}

// FollowUser handles POST /api/users/:id/follow — follow a user.
func (h *FollowHandler) FollowUser(c *gin.Context) {
	targetID, err := parseUintParam(c, "id")
	if err != nil {
		utils.ErrorResponse(c, http.StatusBadRequest, "invalid user id")
		return
	}

	userIDStr, _ := c.Get("user_id")
	userID, err := strconv.ParseUint(userIDStr.(string), 10, 64)
	if err != nil {
		utils.ErrorResponse(c, http.StatusUnauthorized, "invalid token")
		return
	}

	if err := h.followService.Follow(uint(userID), targetID); err != nil {
		switch {
		case errors.Is(err, dao.ErrCannotFollowSelf):
			utils.ErrorResponse(c, http.StatusBadRequest, "cannot follow yourself")
		case errors.Is(err, dao.ErrAlreadyFollowed):
			utils.ErrorResponse(c, http.StatusConflict, "already following this user")
		case errors.Is(err, service.ErrFollowDisabled):
			utils.ErrorResponse(c, http.StatusForbidden, "user has disabled followers")
		default:
			utils.ErrorResponse(c, http.StatusBadRequest, err.Error())
		}
		return
	}

	utils.SuccessResponse(c, http.StatusOK, "followed user", nil)
}

// UnfollowUser handles DELETE /api/users/:id/follow — unfollow a user.
func (h *FollowHandler) UnfollowUser(c *gin.Context) {
	targetID, err := parseUintParam(c, "id")
	if err != nil {
		utils.ErrorResponse(c, http.StatusBadRequest, "invalid user id")
		return
	}

	userIDStr, _ := c.Get("user_id")
	userID, err := strconv.ParseUint(userIDStr.(string), 10, 64)
	if err != nil {
		utils.ErrorResponse(c, http.StatusUnauthorized, "invalid token")
		return
	}

	if err := h.followService.Unfollow(uint(userID), targetID); err != nil {
		if errors.Is(err, dao.ErrFollowNotFound) {
			utils.ErrorResponse(c, http.StatusNotFound, "not following this user")
			return
		}
		utils.ErrorResponse(c, http.StatusInternalServerError, "failed to unfollow user")
		return
	}

	utils.SuccessResponse(c, http.StatusOK, "unfollowed user", nil)
}

// ListFollowers handles GET /api/users/:id/followers — get paginated followers of a user.
func (h *FollowHandler) ListFollowers(c *gin.Context) {
	userID, err := parseUintParam(c, "id")
	if err != nil {
		utils.ErrorResponse(c, http.StatusBadRequest, "invalid user id")
		return
	}

	page, _ := strconv.Atoi(c.DefaultQuery("page", "1"))
	pageSize, _ := strconv.Atoi(c.DefaultQuery("page_size", "20"))

	follows, total, err := h.followService.ListFollowers(userID, page, pageSize)
	if err != nil {
		utils.ErrorResponse(c, http.StatusInternalServerError, "failed to fetch followers")
		return
	}

	utils.SuccessResponse(c, http.StatusOK, "ok", gin.H{
		"followers":  follows,
		"total":      total,
		"page":       page,
		"page_size":  pageSize,
	})
}

// ListFollowing handles GET /api/users/:id/following — get paginated users that a user is following.
func (h *FollowHandler) ListFollowing(c *gin.Context) {
	userID, err := parseUintParam(c, "id")
	if err != nil {
		utils.ErrorResponse(c, http.StatusBadRequest, "invalid user id")
		return
	}

	page, _ := strconv.Atoi(c.DefaultQuery("page", "1"))
	pageSize, _ := strconv.Atoi(c.DefaultQuery("page_size", "20"))

	follows, total, err := h.followService.ListFollowing(userID, page, pageSize)
	if err != nil {
		utils.ErrorResponse(c, http.StatusInternalServerError, "failed to fetch following")
		return
	}

	utils.SuccessResponse(c, http.StatusOK, "ok", gin.H{
		"following":  follows,
		"total":      total,
		"page":       page,
		"page_size":  pageSize,
	})
}
