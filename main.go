package main

import (
	"os"

	"github.com/dhanushs3366/21BIT0623_Backend.git/handler"
	"github.com/dhanushs3366/21BIT0623_Backend.git/services/db"
	"github.com/joho/godotenv"
)

func main() {
	defer db.Close()
	err := godotenv.Load()
	if err != nil {
		panic(err)
	}

	SERVER_PORT := os.Getenv("SERVER_PORT")

	DB, err := db.Init()

	if err != nil {
		panic(err)
	}

	h := handler.Init(DB)

	err = h.Run(SERVER_PORT)

	if err != nil {
		panic(err)
	}
}
