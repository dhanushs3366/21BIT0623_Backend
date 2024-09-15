package handler

import (
	"fmt"
	"mime/multipart"
	"net/http"
	"sync"
	"time"

	"github.com/dhanushs3366/21BIT0623_Backend.git/services"
	"github.com/dhanushs3366/21BIT0623_Backend.git/services/s3service"
	"github.com/labstack/echo/v4"
)

func (h *Hanlder) handleFileUpload(c echo.Context) error {
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

	metadata, err := h.store.GetLatestMetaData()

	if err != nil {
		return c.JSON(http.StatusInternalServerError, err.Error())
	}

	// Cache meta data
	err = h.redis.Add(fmt.Sprintf("user:%d:file:%s", userID, fileID), metadata)

	if err != nil {
		return c.JSON(http.StatusInternalServerError, err.Error())
	}

	return c.JSON(http.StatusOK, "Cached and uploaded the file")
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

func (h *Hanlder) uploadFile(file *multipart.FileHeader, description string, userID uint, errChan chan error, wg *sync.WaitGroup) {
	defer wg.Done()

	src, err := file.Open()
	if err != nil {
		errChan <- fmt.Errorf("error opening file: %w", err)
		return
	}
	defer src.Close()

	key, err := s3service.GenerateKeyForS3(file)
	if err != nil {
		errChan <- fmt.Errorf("error generating key for file: %w", err)
		return
	}

	fileType, err := s3service.GetFileType(file)
	if err != nil {
		errChan <- fmt.Errorf("error determining file type: %w", err)
		return
	}

	err = h.s3.PutObject(src, *file, key)
	if err != nil {
		errChan <- fmt.Errorf("error uploading file to S3: %w", err)
		return
	}

	err = h.store.InsertFile(userID, key)
	if err != nil {
		errChan <- fmt.Errorf("error inserting file into DB: %w", err)
		return
	}

	fileID, err := h.store.GetLatestFileID(userID)
	if err != nil {
		errChan <- fmt.Errorf("error retrieving latest file ID: %w", err)
		return
	}

	err = h.store.InsertMetaData(fileID, file.Filename, uint(file.Size), fileType, description)
	if err != nil {
		errChan <- fmt.Errorf("error inserting file metadata: %w", err)
		return
	}

	errChan <- nil
}

func (h *Hanlder) handleBulkUpload(c echo.Context) error {
	form, err := c.MultipartForm()
	if err != nil {
		return c.JSON(http.StatusBadRequest, err.Error())
	}

	files := form.File["files"]
	description := c.FormValue("description")
	userID, err := services.GetUserIDFromToken(c)
	if err != nil {
		return c.JSON(http.StatusUnauthorized, err.Error())
	}
	var wg sync.WaitGroup
	errChan := make(chan error, len(files))

	for _, file := range files {
		wg.Add(1) // add counter
		go h.uploadFile(file, description, userID, errChan, &wg)
	}

	wg.Wait()
	close(errChan)

	for err := range errChan {
		if err != nil {
			return c.JSON(http.StatusInternalServerError, fmt.Sprintf("Error uploading file: %v", err))
		}
	}

	return c.JSON(http.StatusOK, "All files uploaded successfully!")
}
