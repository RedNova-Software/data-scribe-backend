package main

import (
	"api/shared/util"
	"context"
	"fmt"

	"github.com/aws/aws-lambda-go/events"
	"github.com/aws/aws-lambda-go/lambda"
)

func Handler(ctx context.Context, s3Event events.S3Event) {
	for _, record := range s3Event.Records {
		s3Record := record.S3
		bucket := s3Record.Bucket.Name
		key := s3Record.Object.Key

		fmt.Printf("Processing file: %s from bucket: %s\n", key, bucket)

		// Load CSV file from S3
		file, err := util.GetCSVFileHandle(key)
		if err != nil {
			fmt.Println("Error loading CSV from S3:", err)
			return
		}

		// Process the CSV file
		uniqueValues, err := util.GetUniqueColumnValuesMapInCSV(file)
		if err != nil {
			fmt.Println("Error processing CSV file:", err)
			return
		}

		// Update DynamoDB
		err = util.UpdateReportCsvColumns(key, uniqueValues)
		if err != nil {
			fmt.Println("Error updating DynamoDB:", err)
			return
		}

		// This will let the polling function know that the csv has been updated successfully
		err = util.SetOperationCompleted(key)
		if err != nil {
			fmt.Println("Error setting operation completed:", err)
			return
		}
	}
}

func main() {
	lambda.Start(Handler)
}
