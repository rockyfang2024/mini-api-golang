package service_test

import (
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
	"gorm.io/driver/sqlite"
	"gorm.io/gorm"

	"mini-api-golang/internal/dao"
	"mini-api-golang/internal/models"
	"mini-api-golang/internal/service"
)

// setupTestDB creates an in-memory SQLite DB and migrates all models.
func setupTestDB(t *testing.T) *gorm.DB {
	t.Helper()
	db, err := gorm.Open(sqlite.Open(":memory:"), &gorm.Config{})
	require.NoError(t, err)

	err = db.AutoMigrate(
		&models.User{},
		&models.Post{},
		&models.Like{},
		&models.Repost{},
		&models.Notification{},
		&models.Follow{},
	)
	require.NoError(t, err)
	return db
}

// createUser is a helper to insert a test user into the DB.
func createUser(t *testing.T, db *gorm.DB, username string) *models.User {
	t.Helper()
	user := &models.User{
		Username:     username,
		Email:        username + "@example.com",
		PasswordHash: "hashed",
	}
	require.NoError(t, db.Create(user).Error)
	return user
}

// createPost is a helper to insert a test post into the DB.
func createPost(t *testing.T, db *gorm.DB, authorID uint) *models.Post {
	t.Helper()
	post := &models.Post{
		AuthorID:   authorID,
		Content:    "test content",
		Visibility: models.VisibilityPublic,
	}
	require.NoError(t, db.Create(post).Error)
	return post
}

// ─── Like uniqueness ────────────────────────────────────────────────────────

func TestLike_Uniqueness(t *testing.T) {
	db := setupTestDB(t)

	likeDAO := dao.NewLikeDAO(db)
	postDAO := dao.NewPostDAO(db)
	notifDAO := dao.NewNotificationDAO(db)
	likeSvc := service.NewLikeService(likeDAO, postDAO, notifDAO)

	user := createUser(t, db, "alice")
	author := createUser(t, db, "bob")
	post := createPost(t, db, author.ID)

	// First like should succeed
	err := likeSvc.Like(user.ID, post.ID)
	assert.NoError(t, err)

	// Second like on the same post should return ErrAlreadyLiked
	err = likeSvc.Like(user.ID, post.ID)
	assert.ErrorIs(t, err, dao.ErrAlreadyLiked)
}

func TestLike_Count(t *testing.T) {
	db := setupTestDB(t)

	likeDAO := dao.NewLikeDAO(db)
	postDAO := dao.NewPostDAO(db)
	notifDAO := dao.NewNotificationDAO(db)
	likeSvc := service.NewLikeService(likeDAO, postDAO, notifDAO)

	author := createUser(t, db, "author")
	user1 := createUser(t, db, "user1")
	user2 := createUser(t, db, "user2")
	post := createPost(t, db, author.ID)

	require.NoError(t, likeSvc.Like(user1.ID, post.ID))
	require.NoError(t, likeSvc.Like(user2.ID, post.ID))

	count, err := likeSvc.LikeCount(post.ID)
	assert.NoError(t, err)
	assert.Equal(t, int64(2), count)
}

// ─── Unlike ─────────────────────────────────────────────────────────────────

func TestUnlike_Success(t *testing.T) {
	db := setupTestDB(t)

	likeDAO := dao.NewLikeDAO(db)
	postDAO := dao.NewPostDAO(db)
	notifDAO := dao.NewNotificationDAO(db)
	likeSvc := service.NewLikeService(likeDAO, postDAO, notifDAO)

	user := createUser(t, db, "alice")
	author := createUser(t, db, "bob")
	post := createPost(t, db, author.ID)

	require.NoError(t, likeSvc.Like(user.ID, post.ID))

	// Unlike should succeed
	err := likeSvc.Unlike(user.ID, post.ID)
	assert.NoError(t, err)

	// Count should now be 0
	count, _ := likeSvc.LikeCount(post.ID)
	assert.Equal(t, int64(0), count)
}

func TestUnlike_NotLiked(t *testing.T) {
	db := setupTestDB(t)

	likeDAO := dao.NewLikeDAO(db)
	postDAO := dao.NewPostDAO(db)
	notifDAO := dao.NewNotificationDAO(db)
	likeSvc := service.NewLikeService(likeDAO, postDAO, notifDAO)

	user := createUser(t, db, "alice")
	author := createUser(t, db, "bob")
	post := createPost(t, db, author.ID)

	// Unliking something that was never liked should return ErrLikeNotFound
	err := likeSvc.Unlike(user.ID, post.ID)
	assert.ErrorIs(t, err, dao.ErrLikeNotFound)
}

// ─── Repost uniqueness ──────────────────────────────────────────────────────

func TestRepost_Uniqueness(t *testing.T) {
	db := setupTestDB(t)

	repostDAO := dao.NewRepostDAO(db)
	postDAO := dao.NewPostDAO(db)
	notifDAO := dao.NewNotificationDAO(db)
	repostSvc := service.NewRepostService(repostDAO, postDAO, notifDAO)

	user := createUser(t, db, "alice")
	author := createUser(t, db, "bob")
	post := createPost(t, db, author.ID)

	// First repost should succeed
	err := repostSvc.Repost(user.ID, post.ID)
	assert.NoError(t, err)

	// Second repost should return ErrAlreadyReposted
	err = repostSvc.Repost(user.ID, post.ID)
	assert.ErrorIs(t, err, dao.ErrAlreadyReposted)
}

// ─── Follow / Unfollow ──────────────────────────────────────────────────────

func TestFollow_Success(t *testing.T) {
	db := setupTestDB(t)

	followDAO := dao.NewFollowDAO(db)
	userDAO := dao.NewUserDAO(db)
	notifDAO := dao.NewNotificationDAO(db)
	followSvc := service.NewFollowService(followDAO, userDAO, notifDAO)

	alice := createUser(t, db, "alice")
	bob := createUser(t, db, "bob")

	err := followSvc.Follow(alice.ID, bob.ID)
	assert.NoError(t, err)

	// Verify relationship
	ok, err := followSvc.IsFollowing(alice.ID, bob.ID)
	assert.NoError(t, err)
	assert.True(t, ok)
}

func TestFollow_Duplicate(t *testing.T) {
	db := setupTestDB(t)

	followDAO := dao.NewFollowDAO(db)
	userDAO := dao.NewUserDAO(db)
	notifDAO := dao.NewNotificationDAO(db)
	followSvc := service.NewFollowService(followDAO, userDAO, notifDAO)

	alice := createUser(t, db, "alice")
	bob := createUser(t, db, "bob")

	require.NoError(t, followSvc.Follow(alice.ID, bob.ID))

	// Duplicate follow should return ErrAlreadyFollowed
	err := followSvc.Follow(alice.ID, bob.ID)
	assert.ErrorIs(t, err, dao.ErrAlreadyFollowed)
}

func TestFollow_Self(t *testing.T) {
	db := setupTestDB(t)

	followDAO := dao.NewFollowDAO(db)
	userDAO := dao.NewUserDAO(db)
	notifDAO := dao.NewNotificationDAO(db)
	followSvc := service.NewFollowService(followDAO, userDAO, notifDAO)

	alice := createUser(t, db, "alice")

	err := followSvc.Follow(alice.ID, alice.ID)
	assert.ErrorIs(t, err, dao.ErrCannotFollowSelf)
}

func TestUnfollow_Success(t *testing.T) {
	db := setupTestDB(t)

	followDAO := dao.NewFollowDAO(db)
	userDAO := dao.NewUserDAO(db)
	notifDAO := dao.NewNotificationDAO(db)
	followSvc := service.NewFollowService(followDAO, userDAO, notifDAO)

	alice := createUser(t, db, "alice")
	bob := createUser(t, db, "bob")

	require.NoError(t, followSvc.Follow(alice.ID, bob.ID))

	err := followSvc.Unfollow(alice.ID, bob.ID)
	assert.NoError(t, err)

	ok, _ := followSvc.IsFollowing(alice.ID, bob.ID)
	assert.False(t, ok)
}

func TestUnfollow_NotFollowing(t *testing.T) {
	db := setupTestDB(t)

	followDAO := dao.NewFollowDAO(db)
	userDAO := dao.NewUserDAO(db)
	notifDAO := dao.NewNotificationDAO(db)
	followSvc := service.NewFollowService(followDAO, userDAO, notifDAO)

	alice := createUser(t, db, "alice")
	bob := createUser(t, db, "bob")

	err := followSvc.Unfollow(alice.ID, bob.ID)
	assert.ErrorIs(t, err, dao.ErrFollowNotFound)
}

func TestFollow_MutualFollow(t *testing.T) {
	db := setupTestDB(t)

	followDAO := dao.NewFollowDAO(db)
	userDAO := dao.NewUserDAO(db)
	notifDAO := dao.NewNotificationDAO(db)
	followSvc := service.NewFollowService(followDAO, userDAO, notifDAO)

	alice := createUser(t, db, "alice")
	bob := createUser(t, db, "bob")

	require.NoError(t, followSvc.Follow(alice.ID, bob.ID))
	require.NoError(t, followSvc.Follow(bob.ID, alice.ID))

	aliceFollowsBob, _ := followSvc.IsFollowing(alice.ID, bob.ID)
	bobFollowsAlice, _ := followSvc.IsFollowing(bob.ID, alice.ID)
	assert.True(t, aliceFollowsBob)
	assert.True(t, bobFollowsAlice)
}

// ─── Notifications ──────────────────────────────────────────────────────────

func TestNotification_LikeGeneratesNotification(t *testing.T) {
	db := setupTestDB(t)

	likeDAO := dao.NewLikeDAO(db)
	postDAO := dao.NewPostDAO(db)
	notifDAO := dao.NewNotificationDAO(db)
	notifSvc := service.NewNotificationService(notifDAO)
	likeSvc := service.NewLikeService(likeDAO, postDAO, notifDAO)

	author := createUser(t, db, "author")
	liker := createUser(t, db, "liker")
	post := createPost(t, db, author.ID)

	require.NoError(t, likeSvc.Like(liker.ID, post.ID))

	notifs, total, err := notifSvc.List(author.ID, 1, 20)
	assert.NoError(t, err)
	assert.Equal(t, int64(1), total)
	assert.Equal(t, models.NotificationTypeLike, notifs[0].Type)
	assert.Equal(t, liker.ID, notifs[0].ActorID)
}

func TestNotification_RepostGeneratesNotification(t *testing.T) {
	db := setupTestDB(t)

	repostDAO := dao.NewRepostDAO(db)
	postDAO := dao.NewPostDAO(db)
	notifDAO := dao.NewNotificationDAO(db)
	notifSvc := service.NewNotificationService(notifDAO)
	repostSvc := service.NewRepostService(repostDAO, postDAO, notifDAO)

	author := createUser(t, db, "author")
	reposter := createUser(t, db, "reposter")
	post := createPost(t, db, author.ID)

	require.NoError(t, repostSvc.Repost(reposter.ID, post.ID))

	notifs, total, err := notifSvc.List(author.ID, 1, 20)
	assert.NoError(t, err)
	assert.Equal(t, int64(1), total)
	assert.Equal(t, models.NotificationTypeRepost, notifs[0].Type)
}

func TestNotification_MarkRead(t *testing.T) {
	db := setupTestDB(t)

	likeDAO := dao.NewLikeDAO(db)
	postDAO := dao.NewPostDAO(db)
	notifDAO := dao.NewNotificationDAO(db)
	notifSvc := service.NewNotificationService(notifDAO)
	likeSvc := service.NewLikeService(likeDAO, postDAO, notifDAO)

	author := createUser(t, db, "author")
	liker := createUser(t, db, "liker")
	post := createPost(t, db, author.ID)

	require.NoError(t, likeSvc.Like(liker.ID, post.ID))

	notifs, _, _ := notifSvc.List(author.ID, 1, 20)
	assert.False(t, notifs[0].IsRead)

	err := notifSvc.MarkRead(notifs[0].ID, author.ID)
	assert.NoError(t, err)

	// Re-fetch and verify
	updated, _, _ := notifSvc.List(author.ID, 1, 20)
	assert.True(t, updated[0].IsRead)
}

func TestNotification_MarkAllRead(t *testing.T) {
	db := setupTestDB(t)

	likeDAO := dao.NewLikeDAO(db)
	postDAO := dao.NewPostDAO(db)
	notifDAO := dao.NewNotificationDAO(db)
	notifSvc := service.NewNotificationService(notifDAO)
	likeSvc := service.NewLikeService(likeDAO, postDAO, notifDAO)

	author := createUser(t, db, "author")
	liker1 := createUser(t, db, "liker1")
	liker2 := createUser(t, db, "liker2")
	post := createPost(t, db, author.ID)

	require.NoError(t, likeSvc.Like(liker1.ID, post.ID))
	require.NoError(t, likeSvc.Like(liker2.ID, post.ID))

	err := notifSvc.MarkAllRead(author.ID)
	assert.NoError(t, err)

	notifs, _, _ := notifSvc.List(author.ID, 1, 20)
	for _, n := range notifs {
		assert.True(t, n.IsRead)
	}
}

func TestNotification_FollowerReceivesNewPostNotification(t *testing.T) {
	db := setupTestDB(t)

	followDAO := dao.NewFollowDAO(db)
	userDAO := dao.NewUserDAO(db)
	postDAO := dao.NewPostDAO(db)
	notifDAO := dao.NewNotificationDAO(db)
	notifSvc := service.NewNotificationService(notifDAO)
	followSvc := service.NewFollowService(followDAO, userDAO, notifDAO)
	postSvc := service.NewPostService(postDAO, followDAO, notifDAO)

	author := createUser(t, db, "author")
	follower := createUser(t, db, "follower")

	// follower follows author
	require.NoError(t, followSvc.Follow(follower.ID, author.ID))

	// author creates a public post
	_, err := postSvc.Create(author.ID, "hello world", models.VisibilityPublic)
	assert.NoError(t, err)

	// follower should have received a new_post notification
	// (Note: the follow notification is also in the list)
	notifs, _, err := notifSvc.List(follower.ID, 1, 20)
	assert.NoError(t, err)

	found := false
	for _, n := range notifs {
		if n.Type == models.NotificationTypeNewPost && n.ActorID == author.ID {
			found = true
		}
	}
	assert.True(t, found, "follower should have received a new_post notification")
}

func TestNotification_SelfLikeNoNotification(t *testing.T) {
	db := setupTestDB(t)

	likeDAO := dao.NewLikeDAO(db)
	postDAO := dao.NewPostDAO(db)
	notifDAO := dao.NewNotificationDAO(db)
	notifSvc := service.NewNotificationService(notifDAO)
	likeSvc := service.NewLikeService(likeDAO, postDAO, notifDAO)

	user := createUser(t, db, "alice")
	post := createPost(t, db, user.ID)

	// Like own post — should not generate a notification
	require.NoError(t, likeSvc.Like(user.ID, post.ID))

	_, total, _ := notifSvc.List(user.ID, 1, 20)
	assert.Equal(t, int64(0), total)
}
