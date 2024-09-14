package services

import (
	"errors"
	"net/http"
	"os"
	"time"

	"github.com/dhanushs3366/21BIT0623_Backend.git/models"
	"github.com/golang-jwt/jwt/v5"
	"github.com/labstack/echo/v4"
	"golang.org/x/crypto/bcrypt"
)

const (
	HASHING_ROUNDS = 14
	EXPIRY_TIME    = 24 //in hrs
)

type UserClaims struct {
	ID       uint   `json:"user_id"`
	Username string `json:"username"`
	jwt.RegisteredClaims
}

func HashPassword(password string) (string, error) {
	bytes, err := bcrypt.GenerateFromPassword([]byte(password), HASHING_ROUNDS)

	if err != nil {
		return "", err
	}

	return string(bytes), nil
}

func GenerateJWTToken(user *models.User) (string, error) {
	JWT_SECRET := os.Getenv("JWT_SECRET")
	expirationTime := time.Now().Add(EXPIRY_TIME * time.Hour)

	claims := &UserClaims{
		ID:       user.ID,
		Username: user.Username,
		RegisteredClaims: jwt.RegisteredClaims{
			ExpiresAt: jwt.NewNumericDate(expirationTime),
			IssuedAt:  jwt.NewNumericDate(time.Now()),
		},
	}

	token := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)
	tokenStr, err := token.SignedString([]byte(JWT_SECRET))

	if err != nil {
		return "", err
	}

	return tokenStr, nil
}

// a middlerware to validate the user
func ValidateJWT(next echo.HandlerFunc) echo.HandlerFunc {
	return func(c echo.Context) error {
		JWT_SECRET := []byte(os.Getenv("JWT_SECRET"))
		cookie, err := c.Cookie("auth_token")

		if err != nil {
			if errors.Is(err, http.ErrNoCookie) {
				return c.JSON(http.StatusUnauthorized, err)
			}
			return c.JSON(http.StatusBadRequest, err)
		}
		tokenStr := cookie.Value
		claims := &UserClaims{}

		token, err := jwt.ParseWithClaims(tokenStr, claims, func(t *jwt.Token) (interface{}, error) {
			return JWT_SECRET, nil
		})

		if err != nil {
			if errors.Is(err, jwt.ErrSignatureInvalid) {
				return c.JSON(http.StatusUnauthorized, err.Error())
			}
			return c.JSON(http.StatusBadRequest, err.Error())
		}

		if !token.Valid {
			return c.JSON(http.StatusUnauthorized, "user unauthorized")
		}

		return next(c)
	}
}
