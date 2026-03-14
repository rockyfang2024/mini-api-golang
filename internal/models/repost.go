package models

import "time"

// Repost represents a user reposting (forwarding) a post.
// A user can only repost a given post once (unique constraint on UserID + OriginalPostID).
type Repost struct {
	ID             uint      `gorm:"primaryKey;autoIncrement" json:"id"`
	UserID         uint      `gorm:"not null;index;uniqueIndex:idx_user_post_repost" json:"user_id"`
	OriginalPostID uint      `gorm:"not null;index;uniqueIndex:idx_user_post_repost" json:"original_post_id"`
	User           User      `gorm:"foreignKey:UserID" json:"user,omitempty"`
	OriginalPost   Post      `gorm:"foreignKey:OriginalPostID" json:"-"`
	CreatedAt      time.Time `json:"created_at"`
}
