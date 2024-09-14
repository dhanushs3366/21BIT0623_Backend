package s3service

import (
	"errors"
	"fmt"
	"path/filepath"
	"strings"

	"github.com/google/uuid"
)

type FileType string

const (
	JPEG FileType = "JPEG"
	PNG  FileType = "PNG"
)

func getNewUUID() string {
	uuid := uuid.New()
	return uuid.String()
}
func getFileType(filename string) (FileType, error) {
	ext := strings.ToLower(filepath.Ext(filename))
	switch ext {
	case ".jpeg", ".jpg":
		return JPEG, nil
	case ".png":
		return PNG, nil
	default:
		return "", errors.New("unsupported file type")
	}
}

func generateKeyForS3(filename string) (string, error) {
	fileType, err := getFileType(filename)
	if err != nil {
		return "", err
	}
	uuidKey := getNewUUID()

	key := fmt.Sprintf("%s/%s/%s", fileType, filename, uuidKey)

	return key, nil
}
