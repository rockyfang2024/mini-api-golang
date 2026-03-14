package handler

import (
	"net/http"
	"strconv"

	"github.com/gin-gonic/gin"
	"mini-api-golang/internal/service"
	"mini-api-golang/internal/utils"
)

// NotificationHandler handles notification-related requests.
type NotificationHandler struct {
	notificationService *service.NotificationService
}

// NewNotificationHandler creates a new NotificationHandler.
func NewNotificationHandler(notificationService *service.NotificationService) *NotificationHandler {
	return &NotificationHandler{notificationService: notificationService}
}

// ListNotifications handles GET /api/notifications — get paginated notifications for the current user.
func (h *NotificationHandler) ListNotifications(c *gin.Context) {
	userIDStr, _ := c.Get("user_id")
	userID, err := strconv.ParseUint(userIDStr.(string), 10, 64)
	if err != nil {
		utils.ErrorResponse(c, http.StatusUnauthorized, "invalid token")
		return
	}

	page, _ := strconv.Atoi(c.DefaultQuery("page", "1"))
	pageSize, _ := strconv.Atoi(c.DefaultQuery("page_size", "20"))

	notifications, total, err := h.notificationService.List(uint(userID), page, pageSize)
	if err != nil {
		utils.ErrorResponse(c, http.StatusInternalServerError, "failed to fetch notifications")
		return
	}

	utils.SuccessResponse(c, http.StatusOK, "ok", gin.H{
		"notifications": notifications,
		"total":         total,
		"page":          page,
		"page_size":     pageSize,
	})
}

// MarkNotificationRead handles PUT /api/notifications/:id/read — mark a notification as read.
func (h *NotificationHandler) MarkNotificationRead(c *gin.Context) {
	notifID, err := parseUintParam(c, "id")
	if err != nil {
		utils.ErrorResponse(c, http.StatusBadRequest, "invalid notification id")
		return
	}

	userIDStr, _ := c.Get("user_id")
	userID, err := strconv.ParseUint(userIDStr.(string), 10, 64)
	if err != nil {
		utils.ErrorResponse(c, http.StatusUnauthorized, "invalid token")
		return
	}

	if err := h.notificationService.MarkRead(notifID, uint(userID)); err != nil {
		utils.ErrorResponse(c, http.StatusInternalServerError, "failed to mark notification as read")
		return
	}

	utils.SuccessResponse(c, http.StatusOK, "notification marked as read", nil)
}

// MarkAllNotificationsRead handles PUT /api/notifications/read-all — mark all notifications as read.
func (h *NotificationHandler) MarkAllNotificationsRead(c *gin.Context) {
	userIDStr, _ := c.Get("user_id")
	userID, err := strconv.ParseUint(userIDStr.(string), 10, 64)
	if err != nil {
		utils.ErrorResponse(c, http.StatusUnauthorized, "invalid token")
		return
	}

	if err := h.notificationService.MarkAllRead(uint(userID)); err != nil {
		utils.ErrorResponse(c, http.StatusInternalServerError, "failed to mark all notifications as read")
		return
	}

	utils.SuccessResponse(c, http.StatusOK, "all notifications marked as read", nil)
}
