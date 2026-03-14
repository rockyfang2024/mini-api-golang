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

// RepostHandler handles repost requests.
type RepostHandler struct {
	repostService *service.RepostService
}

// NewRepostHandler creates a new RepostHandler.
func NewRepostHandler(repostService *service.RepostService) *RepostHandler {
	return &RepostHandler{repostService: repostService}
}

// RepostPost handles POST /api/posts/:id/repost — repost a post.
// A user may only repost a given post once.
func (h *RepostHandler) RepostPost(c *gin.Context) {
	originalPostID, err := parseUintParam(c, "id")
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

	if err := h.repostService.Repost(uint(userID), originalPostID); err != nil {
		if errors.Is(err, dao.ErrAlreadyReposted) {
			utils.ErrorResponse(c, http.StatusConflict, "already reposted this post")
			return
		}
		utils.ErrorResponse(c, http.StatusBadRequest, err.Error())
		return
	}

	utils.SuccessResponse(c, http.StatusOK, "post reposted", nil)
}
