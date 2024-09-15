package handler

import (
	"database/sql"
	"fmt"
	"net/http"

	"github.com/dhanushs3366/21BIT0623_Backend.git/services"
	"github.com/dhanushs3366/21BIT0623_Backend.git/services/db"
	redisservice "github.com/dhanushs3366/21BIT0623_Backend.git/services/redisService"
	"github.com/dhanushs3366/21BIT0623_Backend.git/services/s3service"
	"github.com/labstack/echo/v4"
	"github.com/labstack/echo/v4/middleware"
)

type Hanlder struct {
	router *echo.Echo
	store  *db.Store
	s3     *s3service.S3Service
	redis  *redisservice.RedisClient
}

func Init(database *sql.DB) (*Hanlder, error) {

	s3Client, err := s3service.GetNewS3Client()
	if err != nil {
		return nil, err
	}

	rdb, err := redisservice.GetNewRedisClient()
	if err != nil {
		return nil, err
	}
	h := Hanlder{
		router: echo.New(),
		store:  db.GetNewStore(database),
		s3:     s3Client,
		redis:  rdb,
	}

	userGroup := h.router.Group("/user")

	h.router.Use(middleware.Logger())
	h.router.Use(middleware.Recover())

	userGroup.Use(services.ValidateJWT)

	h.router.GET("/", func(c echo.Context) error {
		return c.HTML(http.StatusOK, "<p>Hiiii</p>")
	})

	h.router.POST("/register", h.register)
	h.router.POST("/login", h.login)

	// user routes
	userGroup.GET("/hello", func(c echo.Context) error {
		return c.JSON(http.StatusOK, "hello")
	})
	userGroup.POST("/upload", h.handleFileUpload)
	userGroup.POST("/upload/bulk", h.handleBulkUpload)
	userGroup.GET("/files/metadata", h.getFileMetadata)
	userGroup.GET("/share", h.getPublicURL)
	userGroup.GET("/files/search", h.searchFiles)

	return &h, nil
}

func (h *Hanlder) Run(port string) error {
	err := h.router.Start(fmt.Sprintf(":%s", port))
	return err
}
