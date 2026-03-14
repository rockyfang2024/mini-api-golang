package models

import "time"

// UserSettings captures privacy and interaction settings for a user.
type UserSettings struct {
	ID                   uint      `gorm:"primaryKey;autoIncrement" json:"id"`
	UserID               uint      `gorm:"uniqueIndex;not null" json:"user_id"`
	User                 User      `gorm:"foreignKey:UserID" json:"-"`
	AllowComments        bool      `gorm:"default:true" json:"allow_comments"`
	AllowFollow          bool      `gorm:"default:true" json:"allow_follow"`
	OnlyFollowersCanView bool      `gorm:"default:false" json:"only_followers_can_view"`
	OnlyFollowingCanView bool      `gorm:"default:false" json:"only_following_can_view"`
	CreatedAt            time.Time `json:"created_at"`
	UpdatedAt            time.Time `json:"updated_at"`
}
