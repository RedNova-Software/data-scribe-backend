package util

import (
	"api/shared/constants"
	"bytes"
	"encoding/base64"
	"io"
	"os"

	"github.com/aws/aws-sdk-go/aws"
	"github.com/aws/aws-sdk-go/service/s3"
	"github.com/aws/aws-sdk-go/service/s3/s3manager"
	"github.com/google/uuid"
)

// Returns a handle to the csv file in the local file system
// Downloads it from s3 given its key
func GetCSVFileHandle(s3Key string) (*os.File, error) {
	s3Client, err := GetS3Client(os.Getenv(constants.USEast2))
	if err != nil {
		return nil, err
	}

	const tempFileName = "temp.csv"

	// Create a file to write the S3 Object contents to.
	file, err := os.Create(tempFileName) // Name of the local file to create
	if err != nil {
		return nil, err
	}

	// Write the contents of S3 Object to the file
	downOutput, err := s3Client.GetObject(&s3.GetObjectInput{
		Bucket: aws.String(os.Getenv(constants.S3BucketName)),
		Key:    aws.String(s3Key),
	})
	if err != nil {
		return nil, err
	}

	// Copy the contents of the S3 object to the file
	_, err = io.Copy(file, downOutput.Body)
	return file, err
}

// Uploads the csv to S3 and returns its key
func UploadBase64CSVtoS3(base64Data string) (string, error) {
	// Decode the base64 string
	data, err := base64.StdEncoding.DecodeString(base64Data)
	if err != nil {
		return "", err
	}

	// Get the S3 client (singleton)
	s3Client, err := GetS3Client(os.Getenv(constants.USEast2))
	if err != nil {
		return "", err
	}

	// Create an uploader with the S3 client
	uploader := s3manager.NewUploaderWithClient(s3Client)

	objectKey := uuid.New().String() + ".csv"

	// Upload the file
	_, err = uploader.Upload(&s3manager.UploadInput{
		Bucket: aws.String(os.Getenv(constants.S3BucketName)),
		Key:    aws.String(objectKey),
		Body:   bytes.NewReader(data),
	})

	if err != nil {
		return "", err
	}

	return objectKey, nil
}
