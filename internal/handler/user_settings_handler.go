package handler

import (
	"net/http"
	"strconv"

	"github.com/gin-gonic/gin"
	"mini-api-golang/internal/models"
	"mini-api-golang/internal/service"
	"mini-api-golang/internal/utils"
)

// UserSettingsHandler manages user settings endpoints.
type UserSettingsHandler struct {
	settingsService *service.UserSettingsService
}

// NewUserSettingsHandler creates a new UserSettingsHandler.
func NewUserSettingsHandler(settingsService *service.UserSettingsService) *UserSettingsHandler {
	return &UserSettingsHandler{settingsService: settingsService}
}

// GetSettings handles GET /api/settings.
func (h *UserSettingsHandler) GetSettings(c *gin.Context) {
	userID, ok := currentUserID(c)
	if !ok {
		utils.ErrorResponse(c, http.StatusUnauthorized, "invalid token")
		return
	}

	settings, err := h.settingsService.GetOrCreate(userID)
	if err != nil {
		utils.ErrorResponse(c, http.StatusInternalServerError, "failed to load settings")
		return
	}

	utils.SuccessResponse(c, http.StatusOK, "ok", settings)
}

type updateSettingsRequest struct {
	AllowComments        bool `json:"allow_comments"`
	AllowFollow          bool `json:"allow_follow"`
	OnlyFollowersCanView bool `json:"only_followers_can_view"`
	OnlyFollowingCanView bool `json:"only_following_can_view"`
}

// UpdateSettings handles PUT /api/settings.
func (h *UserSettingsHandler) UpdateSettings(c *gin.Context) {
	userID, ok := currentUserID(c)
	if !ok {
		utils.ErrorResponse(c, http.StatusUnauthorized, "invalid token")
		return
	}

	var req updateSettingsRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		utils.ErrorResponse(c, http.StatusBadRequest, err.Error())
		return
	}

	updates := models.UserSettings{
		AllowComments:        req.AllowComments,
		AllowFollow:          req.AllowFollow,
		OnlyFollowersCanView: req.OnlyFollowersCanView,
		OnlyFollowingCanView: req.OnlyFollowingCanView,
	}

	settings, err := h.settingsService.Update(userID, updates)
	if err != nil {
		utils.ErrorResponse(c, http.StatusInternalServerError, "failed to update settings")
		return
	}

	utils.SuccessResponse(c, http.StatusOK, "settings updated", settings)
}

func currentUserID(c *gin.Context) (uint, bool) {
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
