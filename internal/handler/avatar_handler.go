package handler

import (
	"net/http"
	"strconv"

	"github.com/gin-gonic/gin"
	"mini-api-golang/config"
	"mini-api-golang/internal/service"
	"mini-api-golang/internal/utils"
)

// AvatarHandler handles avatar upload requests.
type AvatarHandler struct {
	avatarService *service.AvatarService
	cfg           *config.Config
}

// NewAvatarHandler creates a new AvatarHandler.
func NewAvatarHandler(avatarService *service.AvatarService, cfg *config.Config) *AvatarHandler {
	return &AvatarHandler{avatarService: avatarService, cfg: cfg}
}

// UploadAvatar handles POST /api/me/avatar — upload or replace the authenticated user's avatar.
func (h *AvatarHandler) UploadAvatar(c *gin.Context) {
	userIDStr, _ := c.Get("user_id")
	userID, err := strconv.ParseUint(userIDStr.(string), 10, 64)
	if err != nil {
		utils.ErrorResponse(c, http.StatusUnauthorized, "invalid token")
		return
	}

	fileHeader, err := c.FormFile("avatar")
	if err != nil {
		utils.ErrorResponse(c, http.StatusBadRequest, "avatar file is required (form field: avatar)")
		return
	}

	urlPath, err := h.avatarService.UploadAvatar(uint(userID), fileHeader)
	if err != nil {
		utils.ErrorResponse(c, http.StatusBadRequest, err.Error())
		return
	}

	utils.SuccessResponse(c, http.StatusOK, "avatar uploaded", gin.H{"avatar_url": urlPath})
}
