package models

type File struct {
	ID         uint   `json:"id"`
	UserID     uint   `json:"user_id"`
	MetadataID uint   `json:"metadata_id"`
	S3Key      string `json:"s3_key"`
}
