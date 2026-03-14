package models

import "time"

// User represents the user model in the application.
type User struct {
	ID           uint      `gorm:"primaryKey;autoIncrement" json:"id"`
	Username     string    `gorm:"size:100;uniqueIndex;not null" json:"username"`
	Email        string    `gorm:"size:200;uniqueIndex;not null" json:"email"`
	PasswordHash string    `gorm:"size:255;not null" json:"-"`
	AvatarURL    string    `gorm:"size:500" json:"avatar_url,omitempty"`
	CreatedAt    time.Time `json:"created_at"`
	UpdatedAt    time.Time `json:"updated_at"`
}