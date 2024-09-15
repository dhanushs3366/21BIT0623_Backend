package models

import (
	"time"
)

type FileMetaData struct {
	ID          uint      `json:"id"`
	FileID      uint      `json:"file_id"`
	FileName    string    `json:"filename"`
	FileSize    uint64    `json:"file_size"`
	ContentType FileType  `json:"content_type"` // should move out the utils from s3service package
	UploadDate  time.Time `json:"upload_date"`
	Description string    `json:"description,omitempty"`
}

type FileType string

const (
	// Image Types
	JPEG FileType = "image/jpeg"
	PNG  FileType = "image/png"
	GIF  FileType = "image/gif"
	BMP  FileType = "image/bmp"
	SVG  FileType = "image/svg+xml"

	// Video Types
	MP4 FileType = "video/mp4"
	AVI FileType = "video/x-msvideo"
	MOV FileType = "video/quicktime"

	// Audio Types
	MP3 FileType = "audio/mpeg"
	WAV FileType = "audio/wav"
	OGG FileType = "audio/ogg"

	// Document Types
	PDF  FileType = "application/pdf"
	DOCX FileType = "application/vnd.openxmlformats-officedocument.wordprocessingml.document"
	DOC  FileType = "application/msword"
	XLSX FileType = "application/vnd.openxmlformats-officedocument.spreadsheetml.sheet"
	XLS  FileType = "application/vnd.ms-excel"
	PPTX FileType = "application/vnd.openxmlformats-officedocument.presentationml.presentation"
	PPT  FileType = "application/vnd.ms-powerpoint"

	// Archive Types
	ZIP FileType = "application/zip"
	RAR FileType = "application/x-rar-compressed"
	GZ  FileType = "application/gzip"

	// Others
	JSON FileType = "application/json"
	XML  FileType = "application/xml"
	HTML FileType = "text/html"
	TXT  FileType = "text/plain"
)
