package s3service

import (
	"fmt"
	"log"
	"mime/multipart"
	"os"
	"sync"
	"time"

	"github.com/aws/aws-sdk-go/aws"
	"github.com/aws/aws-sdk-go/aws/credentials"
	"github.com/aws/aws-sdk-go/aws/session"
	"github.com/aws/aws-sdk-go/service/s3"
)

type S3Service struct {
	client     *s3.S3
	bucketName string
}

func GetNewS3Client() (*S3Service, error) {
	AWS_ACCESS_KEY := os.Getenv("AWS_ACCESS_KEY")
	AWS_SECRET_KEY := os.Getenv("AWS_SECRET_KEY")
	AWS_REGION := os.Getenv("AWS_REGION")
	AWS_BUCKET_NAME := os.Getenv("AWS_BUCKET_NAME")
	sess, err := session.NewSession(&aws.Config{
		Region:      aws.String(AWS_REGION),
		Credentials: credentials.NewStaticCredentials(AWS_ACCESS_KEY, AWS_SECRET_KEY, ""),
	})

	if err != nil {
		return nil, err
	}

	s3Client := s3.New(sess)

	return &S3Service{
		client:     s3Client,
		bucketName: AWS_BUCKET_NAME,
	}, err
}

func (c *S3Service) PutObject(file multipart.File, header multipart.FileHeader, key string) error {

	_, err := c.client.PutObject(&s3.PutObjectInput{
		Bucket:        aws.String(c.bucketName),
		Key:           aws.String(key),
		Body:          file,
		ContentLength: aws.Int64(header.Size),
		ContentType:   aws.String(header.Header.Get("Content-Type")),
	})

	return err
}

func (c *S3Service) GeneratePresignedURL(objectKey string, expiration time.Duration) (string, error) {
	req, _ := c.client.GetObjectRequest(&s3.GetObjectInput{
		Bucket: aws.String(c.bucketName),
		Key:    aws.String(objectKey),
	})

	presignedURL, err := req.Presign(expiration)

	if err != nil {
		return "", err
	}

	return presignedURL, err
}

func (s *S3Service) DeleteExpiredFile(expiredObjKey string) error {
	input := &s3.DeleteObjectInput{
		Bucket: aws.String(s.bucketName),
		Key:    aws.String(expiredObjKey),
	}

	_, err := s.client.DeleteObject(input)
	if err != nil {
		return fmt.Errorf("failed to delete file from S3: %w", err)
	}

	err = s.client.WaitUntilObjectNotExists(&s3.HeadObjectInput{
		Bucket: aws.String(s.bucketName),
		Key:    aws.String(expiredObjKey),
	})
	if err != nil {
		return fmt.Errorf("failed to wait for file deletion: %w", err)
	}

	return nil
}

func (s *S3Service) DeleteExpiredFiles(expiredKeys []string) error {
	var wg sync.WaitGroup
	errCh := make(chan error, len(expiredKeys))

	for _, key := range expiredKeys {
		wg.Add(1)

		go func(objKey string) {
			defer wg.Done()

			if err := s.DeleteExpiredFile(objKey); err != nil {
				log.Printf("Failed to delete file: %s, error: %v", objKey, err)
				errCh <- err
			}
		}(key)
	}

	wg.Wait()
	close(errCh)

	if len(errCh) > 0 {
		return fmt.Errorf("some files failed to delete")
	}

	return nil
}
