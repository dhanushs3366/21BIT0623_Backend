package handler

import (
	"database/sql"
	"fmt"
	"net/http"

	"github.com/dhanushs3366/21BIT0623_Backend.git/services"
	"github.com/dhanushs3366/21BIT0623_Backend.git/services/db"
	"github.com/labstack/echo/v4"
	"github.com/labstack/echo/v4/middleware"
)

type Hanlder struct {
	router *echo.Echo
	store  *db.Store
}

func Init(database *sql.DB) *Hanlder {
	h := Hanlder{
		router: echo.New(),
		store:  db.GetNewStore(database),
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

	return &h
}

func (h *Hanlder) Run(port string) error {
	err := h.router.Start(fmt.Sprintf(":%s", port))
	return err
}
