package handler

import (
	"net/http"
	"time"

	"github.com/dhanushs3366/21BIT0623_Backend.git/services"
	"github.com/labstack/echo/v4"
	"golang.org/x/crypto/bcrypt"
)

func (h *Hanlder) register(c echo.Context) error {
	username := c.FormValue("username")
	password := c.FormValue("password")
	email := c.FormValue("email")

	hashedPassword, err := services.HashPassword(password)
	if err != nil {
		return c.JSON(http.StatusInternalServerError, err.Error())
	}

	err = h.store.CreateUser(username, hashedPassword, email)
	if err != nil {
		return c.JSON(http.StatusInternalServerError, err.Error())
	}

	return c.JSON(http.StatusOK, "User created successfully")
}

func (h *Hanlder) login(c echo.Context) error {
	username := c.FormValue("username")
	password := c.FormValue("password")

	user, err := h.store.GetUser(username)

	if err != nil {
		return c.JSON(http.StatusUnauthorized, "user not found")
	}

	err = bcrypt.CompareHashAndPassword([]byte(user.Password), []byte(password))
	if err != nil {
		return c.JSON(http.StatusUnauthorized, "invalid credentials")
	}

	// log in succesful, give token
	tokenStr, err := services.GenerateJWTToken(user)

	if err != nil {
		return c.JSON(http.StatusInternalServerError, err.Error())
	}

	cookie := http.Cookie{
		Name:    "auth_token",
		Value:   tokenStr,
		Path:    "/",
		Expires: time.Now().Add(services.EXPIRY_TIME * time.Hour),
	}

	c.SetCookie(&cookie)

	return c.JSON(http.StatusOK, "login successful")
}
