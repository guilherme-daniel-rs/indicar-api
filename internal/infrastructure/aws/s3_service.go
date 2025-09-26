package aws

import (
	"bytes"
	"context"
	"indicar-api/configs"

	"github.com/aws/aws-sdk-go-v2/aws"
	"github.com/aws/aws-sdk-go-v2/config"
	"github.com/aws/aws-sdk-go-v2/credentials"
	"github.com/aws/aws-sdk-go-v2/service/s3"
)

type S3Service struct {
	client *s3.Client
	Bucket string
}

func NewS3Service() (*S3Service, error) {
	cfg, err := config.LoadDefaultConfig(context.TODO(),
		config.WithRegion(configs.Get().AWS.Region),
		config.WithCredentialsProvider(credentials.NewStaticCredentialsProvider(
			configs.Get().AWS.AccessKeyID,
			configs.Get().AWS.SecretAccessKey,
			"",
		)),
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

func (s *S3Service) UploadFile(key string, data []byte, contentType string) error {
	_, err := s.client.PutObject(context.TODO(), &s3.PutObjectInput{
		Bucket:      aws.String(s.Bucket),
		Key:         aws.String(key),
		Body:        bytes.NewReader(data),
		ContentType: aws.String(contentType),
	})

	return err
}

func (s *S3Service) GetFileURL(key string) string {
	return "https://" + s.Bucket + ".s3." + configs.Get().AWS.Region + ".amazonaws.com/" + key
}
