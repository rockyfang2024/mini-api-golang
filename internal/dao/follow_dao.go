package dao

import (
	"errors"

	"gorm.io/gorm"
	"mini-api-golang/internal/models"
)

// FollowDAO handles database operations for follow relationships.
type FollowDAO struct {
	DB *gorm.DB
}

// NewFollowDAO creates a new FollowDAO.
func NewFollowDAO(db *gorm.DB) *FollowDAO {
	return &FollowDAO{DB: db}
}

// Create inserts a follow record.
func (d *FollowDAO) Create(follow *models.Follow) error {
	return d.DB.Create(follow).Error
}

// Delete removes a follow relationship. Returns ErrFollowNotFound if it doesn't exist.
func (d *FollowDAO) Delete(followerID, followingID uint) error {
	result := d.DB.
		Where("follower_id = ? AND following_id = ?", followerID, followingID).
		Delete(&models.Follow{})
	if result.Error != nil {
		return result.Error
	}
	if result.RowsAffected == 0 {
		return ErrFollowNotFound
	}
	return nil
}

// Exists returns true if followerID is already following followingID.
func (d *FollowDAO) Exists(followerID, followingID uint) (bool, error) {
	var count int64
	err := d.DB.Model(&models.Follow{}).
		Where("follower_id = ? AND following_id = ?", followerID, followingID).
		Count(&count).Error
	return count > 0, err
}

// ListFollowers returns all followers of a given user (paginated).
func (d *FollowDAO) ListFollowers(userID uint, page, pageSize int) ([]models.Follow, int64, error) {
	var follows []models.Follow
	var total int64

	query := d.DB.Model(&models.Follow{}).Where("following_id = ?", userID)

	if err := query.Count(&total).Error; err != nil {
		return nil, 0, err
	}

	offset := (page - 1) * pageSize
	err := query.
		Preload("Follower").
		Order("created_at DESC").
		Offset(offset).
		Limit(pageSize).
		Find(&follows).Error

	return follows, total, err
}

// ListFollowing returns all users that a given user is following (paginated).
func (d *FollowDAO) ListFollowing(userID uint, page, pageSize int) ([]models.Follow, int64, error) {
	var follows []models.Follow
	var total int64

	query := d.DB.Model(&models.Follow{}).Where("follower_id = ?", userID)

	if err := query.Count(&total).Error; err != nil {
		return nil, 0, err
	}

	offset := (page - 1) * pageSize
	err := query.
		Preload("Following").
		Order("created_at DESC").
		Offset(offset).
		Limit(pageSize).
		Find(&follows).Error

	return follows, total, err
}

// ListFollowerIDs returns all follower user IDs for a given user (used for bulk notifications).
func (d *FollowDAO) ListFollowerIDs(userID uint) ([]uint, error) {
	var follows []models.Follow
	err := d.DB.Model(&models.Follow{}).
		Select("follower_id").
		Where("following_id = ?", userID).
		Find(&follows).Error
	if err != nil {
		return nil, err
	}
	ids := make([]uint, len(follows))
	for i, f := range follows {
		ids[i] = f.FollowerID
	}
	return ids, nil
}

// Sentinel errors.
var (
	ErrFollowNotFound  = errors.New("follow relationship not found")
	ErrAlreadyFollowed = errors.New("already following this user")
	ErrCannotFollowSelf = errors.New("cannot follow yourself")
)
