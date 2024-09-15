package s3service

import (
	"mime/multipart"
	"os"
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
	req, _ := c.client.PutObjectRequest(&s3.PutObjectInput{
		Bucket: aws.String(c.bucketName),
		Key:    aws.String(objectKey),
	})

	presignedURL, err := req.Presign(expiration)

	if err != nil {
		return "", err
	}

	return presignedURL, err
}
