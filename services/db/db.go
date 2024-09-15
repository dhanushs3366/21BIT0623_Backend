package db

import (
	"database/sql"
	"fmt"
	"log"
	"os"

	_ "github.com/lib/pq"
)

var db *sql.DB

// delete objects after hrs
const DELETION_TIME_FOR_S3_OBJECTS = 24

var ErrNoEntityFound error = sql.ErrNoRows

func connect() (*sql.DB, error) {
	DB_USER := os.Getenv("DB_USER")
	DB_PORT := os.Getenv("DB_PORT")
	DB_NAME := os.Getenv("DB_NAME")
	DB_PASSWORD := os.Getenv("DB_PASSWORD")
	DB_HOST := os.Getenv("DB_HOST")

	postgresStr := fmt.Sprintf("user=%s password=%s dbname=%s host=%s port=%s sslmode=disable", DB_USER, DB_PASSWORD, DB_NAME, DB_HOST, DB_PORT)
	var err error

	db, err = sql.Open("postgres", postgresStr)

	if err != nil {
		return nil, err
	}

	err = db.Ping()

	if err != nil {
		return nil, err
	}

	return db, err
}

func Close() error {
	return db.Close()
}

func sync() error {
	store := GetNewStore(db)
	err := store.CreateUserTable()
	if err != nil {
		return err
	}
	log.Print("Created user table")

	err = store.CreateFileTable()
	if err != nil {
		return err
	}
	log.Println("File Table created")

	err = store.CreateFileMetaDataTable()

	if err != nil {
		return err
	}

	log.Println("File meta data table created")
	return nil
}

func Init() (*sql.DB, error) {
	DB, err := connect()
	if err != nil {
		return nil, err
	}

	err = sync()
	if err != nil {
		return nil, err
	}

	return DB, err
}
