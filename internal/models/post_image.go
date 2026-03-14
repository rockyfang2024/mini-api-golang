package models

import "time"

// PostImage represents an uploaded image attached to a post.
type PostImage struct {
	ID        uint      `gorm:"primaryKey;autoIncrement" json:"id"`
	PostID    uint      `gorm:"not null;index" json:"post_id"`
	Post      Post      `gorm:"foreignKey:PostID" json:"-"`
	URL       string    `gorm:"size:500;not null" json:"url"`
	SortOrder int       `gorm:"not null;default:0" json:"sort_order"`
	CreatedAt time.Time `json:"created_at"`
}
