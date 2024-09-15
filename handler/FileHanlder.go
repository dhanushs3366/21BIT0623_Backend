package handler

import (
	"fmt"
	"net/http"
	"time"

	"github.com/dhanushs3366/21BIT0623_Backend.git/services"
	"github.com/dhanushs3366/21BIT0623_Backend.git/services/s3service"
	"github.com/labstack/echo/v4"
)

func (h *Hanlder) uploadFile(c echo.Context) error {
	file, err := c.FormFile("file")
	description := c.FormValue("description")
	if err != nil {
		return c.JSON(http.StatusBadRequest, err.Error())
	}

	src, err := file.Open()
	if err != nil {
		return c.JSON(http.StatusInternalServerError, err.Error())
	}

	defer src.Close()

	key, err := s3service.GenerateKeyForS3(file)
	if err != nil {
		return err
	}
	fileType, err := s3service.GetFileType(file)
	if err != nil {
		return c.JSON(http.StatusBadRequest, err.Error())
	}

	err = h.s3.PutObject(src, *file, key)

	if err != nil {
		return c.JSON(http.StatusInternalServerError, err.Error())
	}

	userID, err := services.GetUserIDFromToken(c)
	if err != nil {
		return c.JSON(http.StatusInternalServerError, err.Error())
	}

	err = h.store.InsertFile(userID, key)
	if err != nil {
		return c.JSON(http.StatusInternalServerError, err.Error())
	}
	fileID, err := h.store.GetLatestFileID(userID)

	if err != nil {
		return c.JSON(http.StatusInternalServerError, err.Error())
	}

	err = h.store.InsertMetaData(fileID, file.Filename, uint(file.Size), fileType, description)
	if err != nil {
		return c.JSON(http.StatusInternalServerError, err.Error())
	}
	return c.JSON(http.StatusOK, fmt.Sprintf("uploaded the file %s", key))
}

func (h *Hanlder) getPreSignedURL(c echo.Context) error {
	fileID := c.QueryParam("fileID")

	// fileID is of string type the sql package handles the conversion
	// fileID is string converting it to uint would add more unecessary err checking
	s3Key, err := h.store.GetFileKey(fileID)

	if err != nil {
		return c.JSON(http.StatusInternalServerError, err.Error())
	}

	preSignedURL, err := h.s3.GeneratePresignedURL(s3Key, time.Hour)

	if err != nil {
		return c.JSON(http.StatusInternalServerError, err.Error())
	}
	return c.JSON(http.StatusOK, map[string]string{"url": preSignedURL})
}

func (h *Hanlder) getFileMetadata(c echo.Context) error {
	userID, err := services.GetUserIDFromToken(c)
	if err != nil {
		return c.JSON(http.StatusUnauthorized, err.Error())
	}

	metadata, err := h.store.GetMetaData(userID)

	if err != nil {
		return c.JSON(http.StatusInternalServerError, err.Error())
	}
	return c.JSON(http.StatusOK, metadata)
}
