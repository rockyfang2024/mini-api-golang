package service

import (
	"mini-api-golang/internal/dao"
	"mini-api-golang/internal/models"
)

// NotificationService handles notification retrieval and read-status updates.
type NotificationService struct {
	notificationDAO *dao.NotificationDAO
}

// NewNotificationService creates a new NotificationService.
func NewNotificationService(notificationDAO *dao.NotificationDAO) *NotificationService {
	return &NotificationService{notificationDAO: notificationDAO}
}

// List returns a paginated list of notifications for the given user.
func (s *NotificationService) List(recipientID uint, page, pageSize int) ([]models.Notification, int64, error) {
	if page < 1 {
		page = 1
	}
	if pageSize < 1 || pageSize > 100 {
		pageSize = 20
	}
	return s.notificationDAO.ListByRecipient(recipientID, page, pageSize)
}

// MarkRead marks a single notification as read for the given recipient.
func (s *NotificationService) MarkRead(notificationID, recipientID uint) error {
	return s.notificationDAO.MarkRead(notificationID, recipientID)
}

// MarkAllRead marks all unread notifications as read for the given recipient.
func (s *NotificationService) MarkAllRead(recipientID uint) error {
	return s.notificationDAO.MarkAllRead(recipientID)
}
