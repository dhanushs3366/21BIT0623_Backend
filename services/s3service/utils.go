package s3service

import (
	"errors"
	"fmt"
	"mime/multipart"

	"github.com/dhanushs3366/21BIT0623_Backend.git/models"
	"github.com/google/uuid"
)

const PUBLIC_DEFAULT_EXPIRATION = 1 //in hrs, default expiry time during upload of a file

func getNewUUID() string {
	uuid := uuid.New()
	return uuid.String()
}
func GetFileType(fileHeader *multipart.FileHeader) (models.FileType, error) {
	mimeType := fileHeader.Header.Get("Content-Type")

	switch mimeType {
	// Image Types
	case "image/jpeg":
		return models.JPEG, nil
	case "image/png":
		return models.PNG, nil
	case "image/gif":
		return models.GIF, nil
	case "image/bmp":
		return models.BMP, nil
	case "image/svg+xml":
		return models.SVG, nil

	// Video Types
	case "video/mp4":
		return models.MP4, nil
	case "video/x-msvideo":
		return models.AVI, nil
	case "video/quicktime":
		return models.MOV, nil

	// Audio Types
	case "audio/mpeg":
		return models.MP3, nil
	case "audio/wav":
		return models.WAV, nil
	case "audio/ogg":
		return models.OGG, nil

	// Document Types
	case "application/pdf":
		return models.PDF, nil
	case "application/vnd.openxmlformats-officedocument.wordprocessingml.document":
		return models.DOCX, nil
	case "application/msword":
		return models.DOC, nil
	case "application/vnd.openxmlformats-officedocument.spreadsheetml.sheet":
		return models.XLSX, nil
	case "application/vnd.ms-excel":
		return models.XLS, nil
	case "application/vnd.openxmlformats-officedocument.presentationml.presentation":
		return models.PPTX, nil
	case "application/vnd.ms-powerpoint":
		return models.PPT, nil

	// Archive Types
	case "application/zip":
		return models.ZIP, nil
	case "application/x-rar-compressed":
		return models.RAR, nil
	case "application/gzip":
		return models.GZ, nil

	// Other Types
	case "application/json":
		return models.JSON, nil
	case "application/xml":
		return models.XML, nil
	case "text/html":
		return models.HTML, nil
	case "text/plain":
		return models.TXT, nil

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
