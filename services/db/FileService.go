package db

import (
	"log"
	"strings"
	"time"

	"github.com/dhanushs3366/21BIT0623_Backend.git/models"
)

func (s *Store) CreateFileTable() error {
	query := `
		CREATE TABLE IF NOT EXISTS FILES(
			ID INT GENERATED ALWAYS AS IDENTITY PRIMARY KEY,
			USER_ID INT NOT NULL,
			S3_KEY VARCHAR(255) UNIQUE NOT NULL,
			CREATED_AT TIMESTAMP NOT NULL,
			UPDATED_AT TIMESTAMP NOT NULL,
			CONSTRAINT fk_file_user
				FOREIGN KEY(USER_ID)
				REFERENCES USERS(ID)
		)
	`

	_, err := s.db.Exec(query)
	return err
}

func (s *Store) InsertFile(userID uint, s3Key string) error {
	query := `
		INSERT INTO FILES(USER_ID,S3_KEY,CREATED_AT,UPDATED_AT)
		VALUES($1,$2,$3,$4)
	`
	now := time.Now()
	_, err := s.db.Exec(query, userID, s3Key, now, now)

	return err
}

func (s *Store) GetFileKey(fileID string, userID uint) (string, error) {
	query := `
		SELECT S3_KEY FROM FILES 
		WHERE ID=$1 AND USER_ID=$2
	`
	row := s.db.QueryRow(query, fileID, userID)

	var s3Key string
	err := row.Scan(&s3Key)

	if err != nil {
		return "", err
	}
	return s3Key, nil
}

// fet the last uploaded fileID
func (s *Store) GetLatestFileID(userID uint) (*models.File, error) {
	query := `
		SELECT * FROM FILES
		WHERE USER_ID=$1 
		ORDER BY CREATED_AT DESC
		LIMIT 1	
	`
	row := s.db.QueryRow(query, userID)
	var file models.File
	err := row.Scan(&file.ID, &file.UserID, &file.S3Key, &file.CreatedAt, &file.UpdatedAt)
	if err != nil {
		return nil, err
	}

	return &file, nil

}

func (s *Store) DeleteExpiredFiles() ([]string, error) {
	expirationTime := time.Now().Add(-DELETION_TIME_FOR_S3_OBJECTS * time.Hour)
	var fileKeys []string

	selectQuery := `
		SELECT F.S3_KEY FROM FILES F
		INNER JOIN FILE_METADATA M ON F.ID=M.FILE_ID
		WHERE F.UPDATED_AT<=$1
	`
	rows, err := s.db.Query(selectQuery, expirationTime)
	if err != nil {
		return nil, err
	}

	for rows.Next() {
		var fileKey string
		err := rows.Scan(&fileKey)

		if err != nil {
			continue
		}
		fileKeys = append(fileKeys, fileKey)
	}

	deleteQuery := `
		DELETE FROM FILES
		WHERE UPDATED_AT<=$1
	`

	_, err = s.db.Exec(deleteQuery, expirationTime)
	if err != nil {
		return nil, err
	}

	log.Println("Deleted files")
	keyStr := strings.Join(fileKeys, ", ")
	log.Println(keyStr)
	return fileKeys, nil
}
