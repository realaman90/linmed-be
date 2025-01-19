package aws

import (
	"fmt"
	"io"

	"github.com/aws/aws-sdk-go/aws"
	"github.com/aws/aws-sdk-go/aws/credentials"
	"github.com/aws/aws-sdk-go/aws/session"
	"github.com/aws/aws-sdk-go/service/s3"
	"github.com/aws/aws-sdk-go/service/s3/s3manager"
)

// S3Client holds the S3 service client
type S3Client struct {
	Uploader *s3manager.Uploader
	S3Client *s3.S3
}

// NewS3Client creates a new S3 client
func NewS3Client(region, Id, key string) (*S3Client, error) {
	sess, err := session.NewSession(&aws.Config{
		Region:      aws.String(region),
		Endpoint:    aws.String("https://s3.us-east-1.amazonaws.com"),
		Credentials: credentials.NewStaticCredentials(Id, key, ""),
	})
	if err != nil {
		return nil, fmt.Errorf("failed to create session: %w", err)
	}

	uploader := s3manager.NewUploader(sess)
	s3Client := s3.New(sess)
	return &S3Client{
		Uploader: uploader,
		S3Client: s3Client,
	}, nil
}

// UploadImage uploads an image to the specified S3 bucket from an io.Reader
func (c *S3Client) UploadImage(bucketName string, file io.Reader, key string) error {

	// Upload the file to S3
	_, err := c.Uploader.Upload(&s3manager.UploadInput{

		Bucket: aws.String(bucketName),
		Key:    aws.String(key),
		Body:   file,
	})
	if err != nil {
		return fmt.Errorf("failed to upload %q to bucket %q, %w", key, bucketName, err)
	}

	return nil
}

// GetImage retrieves an image from the specified S3 bucket using the provided key
func (c *S3Client) GetImage(bucketName string, key string) (io.ReadCloser, error) {
	// Create a request to get the object from S3
	result, err := c.S3Client.GetObject(&s3.GetObjectInput{
		Bucket: aws.String(bucketName),
		Key:    aws.String(key),
	})
	if err != nil {
		return nil, fmt.Errorf("failed to get %q from bucket %q, %w", key, bucketName, err)
	}

	return result.Body, nil
}
