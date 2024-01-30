package util

import (
	"sync"

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
