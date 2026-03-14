package models

import "time"

// Visibility represents the visibility level of a post.
type Visibility string

const (
	VisibilityPublic  Visibility = "public"
	VisibilityPrivate Visibility = "private"
)

// Post represents a microblog post.
type Post struct {
	ID         uint       `gorm:"primaryKey;autoIncrement" json:"id"`
	AuthorID   uint       `gorm:"not null;index"           json:"author_id"`
	Author     User       `gorm:"foreignKey:AuthorID"      json:"author,omitempty"`
	Content    string     `gorm:"type:text;not null"       json:"content"`
	Visibility Visibility `gorm:"size:10;not null;default:public" json:"visibility"`
	CreatedAt  time.Time  `json:"created_at"`
	UpdatedAt  time.Time  `json:"updated_at"`
}
