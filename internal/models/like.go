package models

import "time"

// Like represents a user liking a post.
// A user can only like a given post once (unique constraint on UserID + PostID).
type Like struct {
	ID        uint      `gorm:"primaryKey;autoIncrement" json:"id"`
	UserID    uint      `gorm:"not null;index;uniqueIndex:idx_user_post_like" json:"user_id"`
	PostID    uint      `gorm:"not null;index;uniqueIndex:idx_user_post_like" json:"post_id"`
	User      User      `gorm:"foreignKey:UserID" json:"user,omitempty"`
	Post      Post      `gorm:"foreignKey:PostID" json:"-"`
	CreatedAt time.Time `json:"created_at"`
}
