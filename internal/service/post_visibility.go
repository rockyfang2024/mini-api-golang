package service

import "mini-api-golang/internal/dao"

func canViewUserPosts(settingsDAO *dao.UserSettingsDAO, followDAO *dao.FollowDAO, viewerID, authorID uint) (bool, error) {
	if viewerID == authorID {
		return true, nil
	}
	settings, err := ensureUserSettings(settingsDAO, authorID)
	if err != nil {
		return false, err
	}
	if settings.OnlyFollowersCanView {
		if viewerID == 0 {
			return false, nil
		}
		isFollower, err := followDAO.Exists(viewerID, authorID)
		if err != nil {
			return false, err
		}
		if !isFollower {
			return false, nil
		}
	}
	if settings.OnlyFollowingCanView {
		if viewerID == 0 {
			return false, nil
		}
		isFollowed, err := followDAO.Exists(authorID, viewerID)
		if err != nil {
			return false, err
		}
		if !isFollowed {
			return false, nil
		}
	}
	return true, nil
}
