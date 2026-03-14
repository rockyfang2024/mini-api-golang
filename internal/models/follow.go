package models

import "time"

// Follow represents a follower → following relationship.
// A user can only follow another user once (unique constraint on FollowerID + FollowingID).
type Follow struct {
	ID          uint      `gorm:"primaryKey;autoIncrement" json:"id"`
	FollowerID  uint      `gorm:"not null;index;uniqueIndex:idx_follower_following" json:"follower_id"`
	FollowingID uint      `gorm:"not null;index;uniqueIndex:idx_follower_following" json:"following_id"`
	Follower    User      `gorm:"foreignKey:FollowerID" json:"follower,omitempty"`
	Following   User      `gorm:"foreignKey:FollowingID" json:"following,omitempty"`
	CreatedAt   time.Time `json:"created_at"`
}
