package handler

import (
	"fmt"
	"mime/multipart"
	"net/http"
	"sync"
	"time"

	"github.com/dhanushs3366/21BIT0623_Backend.git/models"
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
	latesFile, err := h.store.GetLatestFileID(userID)

	if err != nil {
		return c.JSON(http.StatusInternalServerError, err.Error())
	}

	err = h.store.InsertMetaData(latesFile.ID, file.Filename, uint(file.Size), fileType, description)
	if err != nil {
		return c.JSON(http.StatusInternalServerError, err.Error())
	}

	metadata, err := h.store.GetLatestMetaData()

	if err != nil {
		return c.JSON(http.StatusInternalServerError, err.Error())
	}

	// Cache meta data
	err = h.redis.Add(fmt.Sprintf("user:%d:file:%d", userID, latesFile.ID), []models.FileMetaData{*metadata})

	if err != nil {
		return c.JSON(http.StatusInternalServerError, err.Error())
	}

	presignedURL, err := h.s3.GeneratePresignedURL(latesFile.S3Key, s3service.PUBLIC_DEFAULT_EXPIRATION*time.Hour)

	if err != nil {
		return c.JSON(http.StatusInternalServerError, err.Error())
	}
	return c.JSON(http.StatusOK, map[string]string{"url": presignedURL})
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

	latesFile, err := h.store.GetLatestFileID(userID)
	if err != nil {
		errChan <- fmt.Errorf("error retrieving latest file ID: %w", err)
		return
	}

	err = h.store.InsertMetaData(latesFile.ID, file.Filename, uint(file.Size), fileType, description)
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

func (h *Hanlder) getPublicURL(c echo.Context) error {
	// in seconds
	expirationTimeStr := c.QueryParam("expirationTime")
	fileID := c.QueryParam("fileID")
	userID, err := services.GetUserIDFromToken(c)

	if err != nil {
		return c.JSON(http.StatusUnauthorized, err.Error())
	}

	if expirationTimeStr == "" || fileID == "" {
		return c.JSON(http.StatusBadRequest, map[string]string{"error": "Error in query params"})
	}
	expirationTime, err := time.ParseDuration(expirationTimeStr + "s")

	if err != nil {
		return c.JSON(http.StatusBadRequest, map[string]string{"error": "invalid time expression"})
	}
	objKey, err := h.store.GetFileKey(fileID, userID)
	if err != nil {
		return c.JSON(http.StatusInternalServerError, map[string]string{"error": "Error getting file key"})
	}
	//  get pre signed url
	URL, err := h.s3.GeneratePresignedURL(objKey, expirationTime)
	if err != nil {
		return c.JSON(http.StatusInternalServerError, map[string]string{"error": "Error generating pre-signed URL"})
	}

	return c.JSON(http.StatusOK, map[string]string{"url": URL})
}

func (h *Hanlder) searchFiles(c echo.Context) error {
	fileName := c.QueryParam("name")
	fileTypeStr := c.QueryParam("type")
	fromDateStr := c.QueryParam("fromDate")
	toDateStr := c.QueryParam("toDate")

	userID, err := services.GetUserIDFromToken(c)
	if err != nil {
		return c.JSON(http.StatusUnauthorized, err.Error())
	}

	cacheKey := fmt.Sprintf("user-%d:fileType-%s:fileName-%s:fromDate-%s:toDate-%s", userID, fileTypeStr, fileName, fromDateStr, toDateStr)

	// Check the cache first
	cachedMetadata, err := h.redis.Get(cacheKey)
	if err == nil {
		return c.JSON(http.StatusOK, cachedMetadata)
	}

	var startDate, endDate time.Time
	if fromDateStr != "" {
		startDate, err = time.Parse(time.RFC3339, fromDateStr)
		if err != nil {
			return c.JSON(http.StatusBadRequest, "Invalid fromDate format")
		}
	}
	if toDateStr != "" {
		endDate, err = time.Parse(time.RFC3339, toDateStr)
		if err != nil {
			return c.JSON(http.StatusBadRequest, "Invalid toDate format")
		}
	}

	files, err := h.store.SearchFiles(fileName, fileTypeStr, startDate, endDate)
	if err != nil {
		return c.JSON(http.StatusInternalServerError, err.Error())
	}

	err = h.redis.Add(cacheKey, files)
	if err != nil {
		return c.JSON(http.StatusInternalServerError, "Failed to cache data")
	}

	return c.JSON(http.StatusOK, files)
}
