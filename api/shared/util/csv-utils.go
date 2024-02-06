package util

import (
	"api/shared/constants"
	"api/shared/models"
	"bytes"
	"encoding/csv"
	"encoding/json"
	"fmt"
	"io"
	"os"
	"path/filepath"

	"github.com/aws/aws-sdk-go/aws"
	"github.com/aws/aws-sdk-go/service/dynamodb"
	"github.com/aws/aws-sdk-go/service/s3"
)

// Returns a handle to the csv file in the local file system
// Downloads it from s3 given its key
func GetCSVFileHandle(s3Key string) (*os.File, error) {
	s3Client, err := GetS3Client(constants.USEast2)
	if err != nil {
		return nil, err
	}

	tmpDir := "/tmp"
	tempFileName := filepath.Join(tmpDir, "temp.csv")

	// Create a file to write the S3 Object contents to.
	file, err := os.Create(tempFileName)
	if err != nil {
		return nil, err
	}

	// Ensure the file is closed in case of an error after this point
	defer func() {
		if err != nil {
			file.Close() // ignore error; Write error takes precedence
		}
	}()

	downOutput, err := s3Client.GetObject(&s3.GetObjectInput{
		Bucket: aws.String(os.Getenv(constants.CsvBucketName)),
		Key:    aws.String(s3Key),
	})
	if err != nil {
		return nil, err
	}
	defer downOutput.Body.Close()

	_, err = io.Copy(file, downOutput.Body)
	if err != nil {
		return nil, err
	}

	// Seek to the beginning of the file to allow for reading again
	_, err = file.Seek(0, io.SeekStart)
	if err != nil {
		return nil, err
	}

	// At this point, don't close the file here; it's intended to be used by the caller
	return file, nil
}

// uniqueValuesInCSV takes theCSV file and returns a map where keys are column names
// and values are slices of unique values in those columns.
func GetUniqueColumnValuesMapInCSV(file *os.File) (models.CsvDataColumnUniqueValuesMap, error) {
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
	result := make(models.CsvDataColumnUniqueValuesMap)
	for header, valuesSet := range uniqueValuesMap {
		for value := range valuesSet {
			result[header] = append(result[header], value)
		}
	}

	return result, nil
}

func UpdateReportCsvColumns(csvid string, csvColumns models.CsvDataColumnUniqueValuesMap) error {
	// Serialize csvColumns to JSON
	jsonData, err := json.Marshal(csvColumns)
	if err != nil {
		return fmt.Errorf("error marshaling csvColumns to JSON: %v", err)
	}

	// Upload JSON to S3
	s3Key, err := uploadColumnDataToS3(csvid+".json", jsonData)
	if err != nil {
		return fmt.Errorf("error uploading csvColumns to S3: %v", err)
	}

	// Store s3Key in DynamoDB
	err = updateDynamoDBWithColumnDataS3Key(csvid, s3Key)
	if err != nil {
		return fmt.Errorf("error updating DynamoDB with S3 key: %v", err)
	}

	return nil
}

// getJSONFromS3 fetches a JSON object from S3 and unmarshals it into a struct.
func GetColumnValuesMapFromS3(s3Key string) (*models.CsvDataColumnUniqueValuesMap, error) {
	s3Client, err := GetS3Client(constants.USEast2)
	if err != nil {
		return nil, err
	}

	// Request the file
	result, err := s3Client.GetObject(&s3.GetObjectInput{
		Bucket: aws.String(os.Getenv(constants.ColumnDataBucketName)),
		Key:    aws.String(s3Key),
	})
	if err != nil {
		return nil, fmt.Errorf("failed to get object: %v", err)
	}
	defer result.Body.Close()

	// Read the S3 object's body
	body, err := io.ReadAll(result.Body)
	if err != nil {
		return nil, fmt.Errorf("failed to read object body: %v", err)
	}

	// Unmarshal the JSON into the Report struct
	var columnValuesMap models.CsvDataColumnUniqueValuesMap
	if err := json.Unmarshal(body, &columnValuesMap); err != nil {
		return nil, fmt.Errorf("failed to unmarshal JSON: %v", err)
	}

	return &columnValuesMap, nil
}

func uploadColumnDataToS3(s3Key string, data []byte) (string, error) {
	s3Client, err := GetS3Client(constants.USEast2)
	if err != nil {
		return "", err
	}

	bucketName := os.Getenv(constants.ColumnDataBucketName)

	// Upload the file to S3
	_, err = s3Client.PutObject(&s3.PutObjectInput{
		Bucket: aws.String(bucketName),
		Key:    aws.String(s3Key),
		Body:   bytes.NewReader(data),
	})
	if err != nil {
		return "", err
	}

	return s3Key, nil
}

func updateDynamoDBWithColumnDataS3Key(csvid, s3Key string) error {
	tableName := os.Getenv(constants.ReportTable)
	dynamoDBClient, err := GetDynamoDBClient(constants.USEast2)
	if err != nil {
		return err
	}

	// Step 1: Query to find the primary key using the GSI
	primaryKey, err := queryPrimaryKeyByCSVID(dynamoDBClient, tableName, csvid)
	if err != nil {
		return fmt.Errorf("error querying primary key by CSVID: %v", err)
	}

	// Step 2: Update the item using the primary key
	input := &dynamodb.UpdateItemInput{
		TableName: aws.String(tableName),
		Key: map[string]*dynamodb.AttributeValue{
			constants.ReportIDField: { // Replace "primaryKeyField" with your table's actual primary key attribute name
				S: aws.String(primaryKey),
			},
		},
		UpdateExpression: aws.String("set " + constants.CSVColumnsS3KeyField + " = :v"),
		ExpressionAttributeValues: map[string]*dynamodb.AttributeValue{
			":v": {S: aws.String(s3Key)},
		},
		ReturnValues: aws.String("UPDATED_NEW"),
	}

	_, err = dynamoDBClient.UpdateItem(input)
	if err != nil {
		return fmt.Errorf("error updating item: %v", err)
	}

	_, err = dynamoDBClient.UpdateItem(input)
	return err
}

// queryPrimaryKeyByCSVID queries the DynamoDB table to find the primary key (id) of an item with a specific csvid.
func queryPrimaryKeyByCSVID(dynamoDBClient *dynamodb.DynamoDB, tableName, csvid string) (string, error) {
	// Replace "CSVIDIndex" with the actual name of the GSI
	// Replace "csvid" with the actual name of the GSI partition key if different
	// Replace "id" with the actual primary key attribute name of your table
	queryInput := &dynamodb.QueryInput{
		TableName: aws.String(tableName),
		IndexName: aws.String(constants.CSVIDField), // The name of the GSI
		KeyConditions: map[string]*dynamodb.Condition{
			constants.CSVIDField: { // The name of the GSI partition key
				ComparisonOperator: aws.String("EQ"),
				AttributeValueList: []*dynamodb.AttributeValue{
					{
						S: aws.String(csvid),
					},
				},
			},
		},
		ProjectionExpression: aws.String(constants.ReportIDField), // Specify the primary key field to be returned
		Limit:                aws.Int64(1),                        // Assuming csvid is unique, we only expect one item
	}

	result, err := dynamoDBClient.Query(queryInput)
	if err != nil {
		return "", err
	}

	if len(result.Items) == 0 {
		return "", nil // No item found matching the csvid
	}

	attrValue, exists := result.Items[0][constants.ReportIDField]
	if !exists {
		return "", nil // Attribute not found in the item
	}

	// Now safely access the .S pointer to get the string value
	primaryKey := attrValue.S
	if primaryKey == nil {
		return "", nil // Primary key is nil
	}

	return *primaryKey, nil
}
