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
	"mini-api-golang/internal/models"
)

// PostImageService handles image uploads for posts.
type PostImageService struct {
	postImageDAO *dao.PostImageDAO
	uploadDir    string
	maxSizeB     int64
}

// NewPostImageService creates a new PostImageService.
func NewPostImageService(postImageDAO *dao.PostImageDAO, uploadDir string, maxSizeMB int64) *PostImageService {
	_ = os.MkdirAll(filepath.Join(uploadDir, "posts"), 0750)
	return &PostImageService{
		postImageDAO: postImageDAO,
		uploadDir:    uploadDir,
		maxSizeB:     maxSizeMB * 1024 * 1024,
	}
}

var allowedPostImageMIMETypes = map[string]string{
	"image/jpeg": ".jpg",
	"image/png":  ".png",
	"image/webp": ".webp",
	"image/gif":  ".gif",
}

// UploadPostImages saves uploaded images and persists records.
func (s *PostImageService) UploadPostImages(postID uint, files []*multipart.FileHeader) ([]models.PostImage, error) {
	images := make([]models.PostImage, 0, len(files))
	savedPaths := make([]string, 0, len(files))
	postDir := filepath.Join(s.uploadDir, "posts")

	for idx, fileHeader := range files {
		if fileHeader.Size > s.maxSizeB {
			cleanupFiles(savedPaths)
			return nil, fmt.Errorf("file size exceeds maximum allowed size of %d MB", s.maxSizeB/(1024*1024))
		}

		src, err := fileHeader.Open()
		if err != nil {
			cleanupFiles(savedPaths)
			return nil, fmt.Errorf("failed to open uploaded file: %w", err)
		}

		detected, err := mimetype.DetectReader(src)
		if err != nil {
			src.Close()
			cleanupFiles(savedPaths)
			return nil, fmt.Errorf("failed to detect file type: %w", err)
		}

		ext, allowed := allowedPostImageMIMETypes[detected.String()]
		if !allowed {
			src.Close()
			cleanupFiles(savedPaths)
			return nil, errors.New("unsupported file type; allowed: jpeg, png, webp, gif")
		}

		if seeker, ok := src.(interface {
			io.Reader
			Seek(offset int64, whence int) (int64, error)
		}); ok {
			if _, err := seeker.Seek(0, io.SeekStart); err != nil {
				src.Close()
				cleanupFiles(savedPaths)
				return nil, fmt.Errorf("failed to seek file: %w", err)
			}
		} else {
			src.Close()
			src, err = fileHeader.Open()
			if err != nil {
				cleanupFiles(savedPaths)
				return nil, fmt.Errorf("failed to re-open uploaded file: %w", err)
			}
		}

		filename := uuid.New().String() + ext
		dstPath := filepath.Join(postDir, filename)
		dst, err := os.Create(dstPath)
		if err != nil {
			src.Close()
			cleanupFiles(savedPaths)
			return nil, fmt.Errorf("failed to create destination file: %w", err)
		}

		if _, err := io.Copy(dst, src); err != nil {
			dst.Close()
			src.Close()
			cleanupFiles(savedPaths)
			return nil, fmt.Errorf("failed to save file: %w", err)
		}
		dst.Close()
		src.Close()

		urlPath := "/uploads/posts/" + filename
		savedPaths = append(savedPaths, dstPath)
		images = append(images, models.PostImage{
			PostID:    postID,
			URL:       urlPath,
			SortOrder: idx + 1,
		})
	}

	if err := s.postImageDAO.CreateBatch(images); err != nil {
		cleanupFiles(savedPaths)
		return nil, fmt.Errorf("failed to save post images: %w", err)
	}

	return images, nil
}

func cleanupFiles(paths []string) {
	for _, path := range paths {
		_ = os.Remove(path)
	}
}
