package db

import (
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

func (s *Store) InsertMetaData(fileID, filename string, fileSize uint, contentType models.FileType, description string) error {
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
