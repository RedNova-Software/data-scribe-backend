package util

import (
	"api/shared/constants"
	"sync"
	"time"

	"github.com/aws/aws-sdk-go/aws"
	"github.com/aws/aws-sdk-go/aws/session"
	"github.com/aws/aws-sdk-go/service/s3"
)

var (
	s3Client    *s3.S3
	s3Once      sync.Once
	s3CreateErr error
)

// GetS3Client returns a singleton S3 client
func GetS3Client(region string) (*s3.S3, error) {
	s3Once.Do(func() {
		s3Client, s3CreateErr = newS3Client(region)
	})
	return s3Client, s3CreateErr
}

func newS3Client(region string) (*s3.S3, error) {
	sess, err := session.NewSession(&aws.Config{
		Region: aws.String(region)},
	)
	if err != nil {
		return nil, err
	}
	return s3.New(sess), nil
}

func GeneratePresignedURL(bucketName, objectKey string, contentType string, duration time.Duration) (string, error) {
	s3Client, err := GetS3Client(constants.USEast2)
	if err != nil {
		return "", err
	}
	// Generate the pre-signed URL for put operations
	req, _ := s3Client.PutObjectRequest(&s3.PutObjectInput{
		Bucket:      aws.String(bucketName),
		Key:         aws.String(objectKey),
		ContentType: aws.String(contentType),
	})

	// Create the pre-signed URL
	urlStr, err := req.Presign(duration)
	if err != nil {
		return "", err
	}

	return urlStr, nil
}
