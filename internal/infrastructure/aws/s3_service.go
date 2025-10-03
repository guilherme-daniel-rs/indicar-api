package aws

import (
	"bytes"
	"context"
	"fmt"
	"indicar-api/configs"
	"mime"
	"path/filepath"
	"strings"
	"time"

	"github.com/aws/aws-sdk-go-v2/aws"
	"github.com/aws/aws-sdk-go-v2/config"
	"github.com/aws/aws-sdk-go-v2/service/s3"
)

type S3Service struct {
	client *s3.Client
	Bucket string
}

func NewS3Service() (*S3Service, error) {
	// Use default AWS configuration which will automatically:
	// 1. Use IAM roles if running on EC2/ECS/Lambda
	// 2. Use environment variables if available
	// 3. Use AWS credentials file if available
	// 4. Use EC2 instance metadata if available
	cfg, err := config.LoadDefaultConfig(context.TODO(),
		config.WithRegion(configs.Get().AWS.Region),
	)
	if err != nil {
		return nil, err
	}

	client := s3.NewFromConfig(cfg)

	return &S3Service{
		client: client,
		Bucket: configs.Get().AWS.S3Bucket,
	}, nil
}

// UploadFile uploads a file to S3
func (s *S3Service) UploadFile(key string, data []byte, contentType string) error {
	_, err := s.client.PutObject(context.TODO(), &s3.PutObjectInput{
		Bucket:      aws.String(s.Bucket),
		Key:         aws.String(key),
		Body:        bytes.NewReader(data),
		ContentType: aws.String(contentType),
	})

	return err
}

// GetFileURL returns a public URL for the file
func (s *S3Service) GetFileURL(key string) string {
	return "https://" + s.Bucket + ".s3." + configs.Get().AWS.Region + ".amazonaws.com/" + key
}

// GetPresignedURL generates a pre-signed URL for secure file access
func (s *S3Service) GetPresignedURL(key string, expiration time.Duration) (string, error) {
	presignClient := s3.NewPresignClient(s.client)

	request, err := presignClient.PresignGetObject(context.TODO(), &s3.GetObjectInput{
		Bucket: aws.String(s.Bucket),
		Key:    aws.String(key),
	}, func(opts *s3.PresignOptions) {
		opts.Expires = expiration
	})

	if err != nil {
		return "", err
	}

	return request.URL, nil
}

// ValidateFileType validates if the file type is allowed
func (s *S3Service) ValidateFileType(filename string, allowedTypes []string) error {
	ext := strings.ToLower(filepath.Ext(filename))
	contentType := mime.TypeByExtension(ext)

	if contentType == "" {
		return fmt.Errorf("unable to determine file type for extension: %s", ext)
	}

	for _, allowedType := range allowedTypes {
		if strings.Contains(contentType, allowedType) {
			return nil
		}
	}

	return fmt.Errorf("file type %s is not allowed. Allowed types: %v", contentType, allowedTypes)
}

// ValidateFileSize validates if the file size is within limits
func (s *S3Service) ValidateFileSize(size int64, maxSize int64) error {
	if size > maxSize {
		return fmt.Errorf("file size %d bytes exceeds maximum allowed size of %d bytes", size, maxSize)
	}
	return nil
}

// DeleteFile deletes a file from S3
func (s *S3Service) DeleteFile(key string) error {
	_, err := s.client.DeleteObject(context.TODO(), &s3.DeleteObjectInput{
		Bucket: aws.String(s.Bucket),
		Key:    aws.String(key),
	})
	return err
}

// FileExists checks if a file exists in S3
func (s *S3Service) FileExists(key string) (bool, error) {
	_, err := s.client.HeadObject(context.TODO(), &s3.HeadObjectInput{
		Bucket: aws.String(s.Bucket),
		Key:    aws.String(key),
	})

	if err != nil {
		// Check if it's a "not found" error
		if strings.Contains(err.Error(), "NoSuchKey") || strings.Contains(err.Error(), "NotFound") {
			return false, nil
		}
		return false, err
	}

	return true, nil
}
