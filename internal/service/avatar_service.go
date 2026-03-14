package service

import (
	"errors"
	"fmt"
	"io"
	"mime/multipart"
	"os"
	"path/filepath"

	"github.com/gabriel-vasile/mimetype"
	"github.com/google/uuid"
	"mini-api-golang/internal/dao"
)

// AvatarService handles avatar file uploads for users.
type AvatarService struct {
	userDAO   *dao.UserDAO
	uploadDir string
	maxSizeB  int64
}

// NewAvatarService creates a new AvatarService.
// It also ensures the avatars subdirectory exists so we don't create it on
// every upload request.
func NewAvatarService(userDAO *dao.UserDAO, uploadDir string, maxSizeMB int64) *AvatarService {
	// Eagerly create the avatar directory so upload requests don't need to.
	_ = os.MkdirAll(filepath.Join(uploadDir, "avatars"), 0750)
	return &AvatarService{
		userDAO:   userDAO,
		uploadDir: uploadDir,
		maxSizeB:  maxSizeMB * 1024 * 1024,
	}
}

// allowedMIMETypes contains the permitted MIME types for avatar uploads.
var allowedMIMETypes = map[string]string{
	"image/jpeg": ".jpg",
	"image/png":  ".png",
	"image/webp": ".webp",
	"image/gif":  ".gif",
}

// UploadAvatar saves an uploaded file to disk and updates the user's AvatarURL.
// It validates file size and detects the true MIME type by reading file content
// (not just the extension) to prevent spoofing attacks.
// Returns the relative URL path for the saved file.
func (s *AvatarService) UploadAvatar(userID uint, fileHeader *multipart.FileHeader) (string, error) {
	// Validate file size
	if fileHeader.Size > s.maxSizeB {
		return "", fmt.Errorf("file size exceeds maximum allowed size of %d MB", s.maxSizeB/(1024*1024))
	}

	// Open source file to detect MIME type from content
	src, err := fileHeader.Open()
	if err != nil {
		return "", fmt.Errorf("failed to open uploaded file: %w", err)
	}
	defer src.Close()

	// Detect MIME type by reading the first bytes of the file content
	detected, err := mimetype.DetectReader(src)
	if err != nil {
		return "", fmt.Errorf("failed to detect file type: %w", err)
	}

	ext, allowed := allowedMIMETypes[detected.String()]
	if !allowed {
		return "", errors.New("unsupported file type; allowed: jpeg, png, webp, gif")
	}

	// Seek back to the beginning of the file after MIME detection
	type readSeeker interface {
		io.Reader
		Seek(offset int64, whence int) (int64, error)
	}
	if rs, ok := src.(readSeeker); ok {
		if _, err := rs.Seek(0, io.SeekStart); err != nil {
			return "", fmt.Errorf("failed to seek file: %w", err)
		}
	} else {
		// Fallback: re-open the file
		src.Close()
		src, err = fileHeader.Open()
		if err != nil {
			return "", fmt.Errorf("failed to re-open uploaded file: %w", err)
		}
	}

	// Sanitise: derive extension from detected MIME, ignoring the client-provided filename

	// Generate a unique filename to avoid collisions
	filename := uuid.New().String() + ext
	avatarDir := filepath.Join(s.uploadDir, "avatars")
	dstPath := filepath.Join(avatarDir, filename)

	// Create destination file
	dst, err := os.Create(dstPath)
	if err != nil {
		return "", fmt.Errorf("failed to create destination file: %w", err)
	}
	defer dst.Close()

	if _, err := io.Copy(dst, src); err != nil {
		return "", fmt.Errorf("failed to save file: %w", err)
	}

	// Build a URL path that will be served by the static file handler
	urlPath := "/uploads/avatars/" + filename

	// Persist the avatar URL to the user record
	user, err := s.userDAO.GetByID(userID)
	if err != nil {
		return "", fmt.Errorf("user not found: %w", err)
	}
	user.AvatarURL = urlPath
	if err := s.userDAO.Update(user); err != nil {
		return "", fmt.Errorf("failed to update user avatar: %w", err)
	}

	return urlPath, nil
}
