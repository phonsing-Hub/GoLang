package helper

import (
	"fmt"
	"github.com/go-playground/validator/v10"
	"github.com/google/uuid"
	"github.com/phonsing-Hub/GoLang/internal/utils/response"
	"mime/multipart"
	"os"
	"path/filepath"
	"time"
)

var VLD = validator.New()
var allowedImageTypes = map[string]bool{
	"image/png":  true,
	"image/jpeg": true,
	"image/webp": true,
}

const MaxAvatarSize = 2 * 1024 * 1024 // 2MB

func StringToDate(dateString string, layout string) (time.Time, error) {
	t, err := time.Parse(layout, dateString)
	if err != nil {
		return time.Time{}, response.Fail(nil, "INVALID_DATE_FORMAT", "Invalid date format", 400)
	}
	return t, nil
}

func PrepareAvatarPath(savePath string) error {
	fullPath := filepath.Join("./static/uploads", savePath)
	return os.MkdirAll(filepath.Dir(fullPath), 0755)
}

func ValidateAndRenameAvatar(file *multipart.FileHeader) (string, error) {
	if file.Size > MaxAvatarSize {
		return "", fmt.Errorf("file too large: max size is 2MB")
	}

	mimeType := file.Header.Get("Content-Type")
	if !allowedImageTypes[mimeType] {
		return "", fmt.Errorf("unsupported file type: %s", mimeType)
	}

	ext := filepath.Ext(file.Filename)
	if ext == "" {
		switch mimeType {
		case "image/jpeg":
			ext = ".jpg"
		case "image/png":
			ext = ".png"
		case "image/webp":
			ext = ".webp"
		default:
			ext = ".img"
		}
	}

	newName := fmt.Sprintf("avatars/%s%s", uuid.New().String(), ext)
	return newName, nil
}
