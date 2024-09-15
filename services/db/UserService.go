package db

import (
	"database/sql"
	"errors"
	"time"

	"github.com/dhanushs3366/21BIT0623_Backend.git/models"
)

type Store struct {
	db *sql.DB
}

func GetNewStore(db *sql.DB) *Store {
	return &Store{db: db}
}

func (s *Store) CreateUserTable() error {
	query := `
		CREATE TABLE IF NOT EXISTS USERS (
			ID INT GENERATED ALWAYS AS IDENTITY PRIMARY KEY,
			USERNAME VARCHAR(255) UNIQUE NOT NULL,
			PASSWORD VARCHAR(255) NOT NULL,
			EMAIL VARCHAR(255) UNIQUE NOT NULL,
			CREATED_AT TIMESTAMP NOT NULL,
			UPDATED_AT TIMESTAMP NOT NULL
		)
	`

	_, err := s.db.Exec(query)

	return err
}

func (s *Store) CreateUser(username, password, email string) error {
	query := `
		INSERT INTO USERS (USERNAME,PASSWORD,EMAIL,CREATED_AT,UPDATED_AT)
		VALUES($1,$2,$3,$4,$5)
	`
	now := time.Now()
	_, err := s.db.Exec(query, username, password, email, now, now)

	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			return ErrNoEntityFound
		}
		return err
	}

	return nil
}

func (s *Store) GetUser(username string) (*models.User, error) {
	query := `
		SELECT * FROM USERS U 
		WHERE U.USERNAME=$1
	`

	row := s.db.QueryRow(query, username)
	var user models.User
	err := row.Scan(&user.ID, &user.Username, &user.Password, &user.Email, &user.CreatedAt, &user.UpdatedAt)
	if err != nil {
		return nil, err
	}
	return &user, err
}
