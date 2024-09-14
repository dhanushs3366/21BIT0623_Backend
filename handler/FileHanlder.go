package handler

import (
	"net/http"

	"github.com/labstack/echo/v4"
)

func (h *Hanlder) uploadFile(c echo.Context) error {
	file, err := c.FormFile("file")

	if err != nil {
		return c.JSON(http.StatusBadRequest, err.Error())
	}

	src, err := file.Open()
	if err != nil {
		return c.JSON(http.StatusInternalServerError, err.Error())
	}

	defer src.Close()

	err = h.s3.PutObject(src, *file)

	if err != nil {
		return c.JSON(http.StatusInternalServerError, err.Error())
	}

	return c.JSON(http.StatusOK, "uploaded the file")
}
