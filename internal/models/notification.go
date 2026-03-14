package models

import "time"

// NotificationType enumerates the types of in-app notifications.
type NotificationType string

const (
	NotificationTypeLike    NotificationType = "like"
	NotificationTypeRepost  NotificationType = "repost"
	NotificationTypeFollow  NotificationType = "follow"
	NotificationTypeNewPost NotificationType = "new_post"
)

// Notification represents a notification sent to a user.
type Notification struct {
	ID          uint             `gorm:"primaryKey;autoIncrement" json:"id"`
	RecipientID uint             `gorm:"not null;index" json:"recipient_id"`
	ActorID     uint             `gorm:"not null" json:"actor_id"`
	Type        NotificationType `gorm:"size:20;not null" json:"type"`
	PostID      *uint            `gorm:"index" json:"post_id,omitempty"`
	IsRead      bool             `gorm:"default:false" json:"is_read"`
	Actor       User             `gorm:"foreignKey:ActorID" json:"actor,omitempty"`
	Post        *Post            `gorm:"foreignKey:PostID" json:"post,omitempty"`
	CreatedAt   time.Time        `json:"created_at"`
}
