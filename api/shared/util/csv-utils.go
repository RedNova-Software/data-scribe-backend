package util

import (
	"api/shared/constants"
	"io"
	"os"

	"github.com/aws/aws-sdk-go/aws"
	"github.com/aws/aws-sdk-go/service/s3"
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
