package s3service

import (
	"errors"
	"fmt"
	"mime/multipart"

	"github.com/dhanushs3366/21BIT0623_Backend.git/models"
	"github.com/google/uuid"
)

func getNewUUID() string {
	uuid := uuid.New()
	return uuid.String()
}
func GetFileType(fileHeader *multipart.FileHeader) (models.FileType, error) {
	mimeType := fileHeader.Header.Get("Content-Type")

	switch mimeType {
	case "image/jpeg":
		return models.JPEG, nil
	case "image/png":
		return models.PNG, nil
	case "video/mp4":
		return models.MP4, nil
	case "audio/mpeg":
		return models.MP3, nil
	default:
		return "", errors.New("unsupported format")
	}
}

func GenerateKeyForS3(header *multipart.FileHeader) (string, error) {
	fileType, err := GetFileType(header)
	if err != nil {
		return "", err
	}
	uuidKey := getNewUUID()

	key := fmt.Sprintf("%s/%s/%s", fileType, header.Filename, uuidKey)

	return key, nil
}
