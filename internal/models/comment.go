package models

import "time"

// Comment represents a comment on a post or another comment.
type Comment struct {
	ID              uint      `gorm:"primaryKey;autoIncrement" json:"id"`
	PostID          uint      `gorm:"not null;index" json:"post_id"`
	ParentCommentID *uint     `gorm:"index" json:"parent_comment_id,omitempty"`
	AuthorID        uint      `gorm:"not null;index" json:"author_id"`
	Author          User      `gorm:"foreignKey:AuthorID" json:"author,omitempty"`
	Content         string    `gorm:"type:text;not null" json:"content"`
	CreatedAt       time.Time `json:"created_at"`
	UpdatedAt       time.Time `json:"updated_at"`
}
