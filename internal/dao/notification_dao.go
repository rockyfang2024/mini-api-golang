package dao

import (
	"gorm.io/gorm"
	"mini-api-golang/internal/models"
)

// NotificationDAO handles database operations for notifications.
type NotificationDAO struct {
	DB *gorm.DB
}

// NewNotificationDAO creates a new NotificationDAO.
func NewNotificationDAO(db *gorm.DB) *NotificationDAO {
	return &NotificationDAO{DB: db}
}

// Create inserts a new notification.
func (d *NotificationDAO) Create(n *models.Notification) error {
	return d.DB.Create(n).Error
}

// BatchCreate inserts multiple notifications in a single database call.
func (d *NotificationDAO) BatchCreate(notifications []*models.Notification) error {
	if len(notifications) == 0 {
		return nil
	}
	return d.DB.Create(notifications).Error
}

// ListByRecipient returns paginated notifications for a given recipient, newest first.
func (d *NotificationDAO) ListByRecipient(recipientID uint, page, pageSize int) ([]models.Notification, int64, error) {
	var notifications []models.Notification
	var total int64

	query := d.DB.Model(&models.Notification{}).Where("recipient_id = ?", recipientID)

	if err := query.Count(&total).Error; err != nil {
		return nil, 0, err
	}

	offset := (page - 1) * pageSize
	err := query.
		Preload("Actor").
		Preload("Post").
		Order("created_at DESC").
		Offset(offset).
		Limit(pageSize).
		Find(&notifications).Error

	return notifications, total, err
}

// MarkRead marks a single notification as read (only for the given recipient).
func (d *NotificationDAO) MarkRead(notificationID, recipientID uint) error {
	return d.DB.Model(&models.Notification{}).
		Where("id = ? AND recipient_id = ?", notificationID, recipientID).
		Update("is_read", true).Error
}

// MarkAllRead marks all unread notifications as read for a recipient.
func (d *NotificationDAO) MarkAllRead(recipientID uint) error {
	return d.DB.Model(&models.Notification{}).
		Where("recipient_id = ? AND is_read = ?", recipientID, false).
		Update("is_read", true).Error
}
