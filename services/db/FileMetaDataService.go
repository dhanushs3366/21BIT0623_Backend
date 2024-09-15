package db

import (
	"log"
	"strings"
	"time"

	"github.com/dhanushs3366/21BIT0623_Backend.git/models"
)

// one file->one metadata fileID
// no need to add updated and created field refer the files table
func (s *Store) CreateFileMetaDataTable() error {
	query := `
		CREATE TABLE IF NOT EXISTS FILE_METADATA(
			ID INT GENERATED ALWAYS AS IDENTITY PRIMARY KEY,
			FILE_ID INT NOT NULL,
			FILE_NAME VARCHAR(255),
			FILE_SIZE INT,
			CONTENT_TYPE VARCHAR(255),
			UPLOAD_DATE TIMESTAMP,
			DESCRIPTION TEXT,
			CONSTRAINT unique_file_id UNIQUE(FILE_ID)
		)
	`

	_, err := s.db.Exec(query)

	return err
}

func (s *Store) InsertMetaData(fileID uint, filename string, fileSize uint, contentType models.FileType, description string) error {
	query := `
		INSERT INTO FILE_METADATA(FILE_ID,FILE_NAME,FILE_SIZE,CONTENT_TYPE,UPLOAD_DATE,DESCRIPTION)
		VALUES($1,$2,$3,$4,$5,$6)
	`
	uploadDate := time.Now()

	_, err := s.db.Exec(query, fileID, filename, fileSize, contentType, uploadDate, description)
	return err
}

// get all the files metadata for a user
func (s *Store) GetMetaData(userID uint) ([]models.FileMetaData, error) {
	var metadata []models.FileMetaData
	query := `
		SELECT M.ID, M.FILE_ID, M.FILE_NAME,M.FILE_SIZE,M.CONTENT_TYPE,M.UPLOAD_DATE,M.DESCRIPTION
			FROM FILE_METADATA M INNER JOIN FILES F 
			ON M.FILE_ID=F.ID
			WHERE F.USER_ID=$1
	`

	rows, err := s.db.Query(query, userID)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	for rows.Next() {
		var data models.FileMetaData
		err := rows.Scan(&data.ID, &data.FileID, &data.FileName, &data.FileSize, &data.ContentType, &data.UploadDate, &data.Description)
		if err != nil {
			continue
		}

		metadata = append(metadata, data)
	}

	return metadata, nil
}

// get the last uploaded file's metadata
func (s *Store) GetLatestMetaData() (*models.FileMetaData, error) {
	query := `
		SELECT M.* FROM FILE_METADATA M
		INNER JOIN FILES F 
		ON F.ID=M.FILE_ID
		ORDER BY F.CREATED_AT DESC
		LIMIT 1

	`
	var metadata models.FileMetaData
	row := s.db.QueryRow(query)
	err := row.Scan(&metadata.ID, &metadata.FileID, &metadata.FileName, &metadata.FileSize, &metadata.ContentType, &metadata.UploadDate, &metadata.Description)
	if err != nil {
		return nil, err
	}
	return &metadata, nil
}

// get the last uploaded file's metadata
func (s *Store) SearchFiles(name, fileType string, startDate, endDate time.Time) ([]models.FileMetaData, error) {
	query := "SELECT * FROM FILE_METADATA WHERE 1=1"
	var args []interface{}
	var conditions []string

	if name != "" {
		conditions = append(conditions, "FILE_NAME ILIKE $1")
		args = append(args, "%"+name+"%")
	}
	if fileType != "" {
		conditions = append(conditions, "CONTENT_TYPE  ILIKE $2")
		args = append(args, "%"+fileType+"%")
	}
	if !startDate.IsZero() && !endDate.IsZero() {
		conditions = append(conditions, "UPLOAD_DATE BETWEEN $3 AND $4")
		args = append(args, startDate, endDate)
	}

	if len(conditions) > 0 {
		query += " AND " + strings.Join(conditions, " AND ")
	}

	log.Println(args...)
	rows, err := s.db.Query(query, args...)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var files []models.FileMetaData
	for rows.Next() {
		var file models.FileMetaData
		err := rows.Scan(&file.ID, &file.FileID, &file.FileName, &file.FileSize, &file.ContentType, &file.UploadDate, &file.Description)
		if err != nil {
			return nil, err
		}
		files = append(files, file)
	}

	if err := rows.Err(); err != nil {
		return nil, err
	}
	return files, nil
}
