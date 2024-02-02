package util

import (
	"api/shared/constants"
	"encoding/csv"
	"fmt"
	"io"
	"os"

	"github.com/aws/aws-sdk-go/aws"
	"github.com/aws/aws-sdk-go/service/dynamodb"
	"github.com/aws/aws-sdk-go/service/dynamodb/dynamodbattribute"
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

// uniqueValuesInCSV takes theCSV file and returns a map where keys are column names
// and values are slices of unique values in those columns.
func UniqueValuesInCSV(file *os.File) (map[string][]string, error) {

	reader := csv.NewReader(file)
	headers, err := reader.Read()
	if err != nil {
		return nil, err
	}

	// Initialize map to hold sets for unique values per column
	uniqueValuesMap := make(map[string]map[string]bool)
	for _, header := range headers {
		uniqueValuesMap[header] = make(map[string]bool)
	}

	// Iterate over each record
	for {
		record, err := reader.Read()
		if err == io.EOF {
			break
		}
		if err != nil {
			return nil, err
		}

		for i, value := range record {
			// Update the set of unique values for the corresponding column
			uniqueValuesMap[headers[i]][value] = true
		}
	}

	// Convert sets to slices
	result := make(map[string][]string)
	for header, valuesSet := range uniqueValuesMap {
		for value := range valuesSet {
			result[header] = append(result[header], value)
		}
	}

	return result, nil
}

func UpdateReportCsvColumns(csvid string, csvColumns map[string][]string) error {
	tableName := os.Getenv(constants.ReportTable)
	dynamoDBClient, err := GetDynamoDBClient(constants.USEast2)

	if err != nil {
		return fmt.Errorf("error getting dynamodb client: %v", err)
	}

	// Marshal the CSVColumns map to a DynamoDB attribute value
	updatedCSVColumns, err := dynamodbattribute.MarshalMap(csvColumns)
	if err != nil {
		return err
	}

	// Prepare update input
	input := &dynamodb.UpdateItemInput{
		TableName: aws.String(tableName),
		Key: map[string]*dynamodb.AttributeValue{
			constants.CSVIDField: {
				S: aws.String(csvid),
			},
		},
		UpdateExpression:          aws.String("set " + constants.CSVColumnsField + " = :v"),
		ExpressionAttributeValues: map[string]*dynamodb.AttributeValue{":v": {M: updatedCSVColumns}},
		ReturnValues:              aws.String("UPDATED_NEW"),
	}

	// Update the item in DynamoDB
	_, err = dynamoDBClient.UpdateItem(input)
	return err
}
