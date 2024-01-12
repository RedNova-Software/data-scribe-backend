package util

import (
	"api/shared/constants"
	"api/shared/models"
	"fmt"
	"github.com/aws/aws-sdk-go/aws"
	"github.com/aws/aws-sdk-go/aws/session"
	"github.com/aws/aws-sdk-go/service/dynamodb"
	"github.com/aws/aws-sdk-go/service/dynamodb/dynamodbattribute"
    "log"
)

func PutItem(tableName string, item interface{}) error {
    dynamoDBClient, err := newDynamoDBClient(string(constants.USEast2))
    if err != nil {
        return err
    }
    av, err := dynamodbattribute.MarshalMap(item)
    if err != nil {
        return err
    }

    input := &dynamodb.PutItemInput{
        Item:      av,
        TableName: aws.String(tableName),
    }
    _, err = dynamoDBClient.PutItem(input)
    if err != nil {
        return err
    }
    return nil
}

func AddElementToReportList(tableName, reportID, listPath string, element interface{}) error {
    dynamoDBClient, err := newDynamoDBClient(string(constants.USEast2))
    if err != nil {
        log.Printf("Error creating new DynamoDB client: %v", err)
        return err
    }

    var updateExpression string
    expressionAttributeValues := map[string]*dynamodb.AttributeValue{}

    marshaledElement, err := dynamodbattribute.Marshal(element)
    if err != nil {
        log.Printf("Error marshaling element: %v", err)
        return err
    }
    log.Printf(marshaledElement.String())
    updateExpression = fmt.Sprintf("SET %s = list_append(if_not_exists(%s, :emptyList), :elem)", listPath, listPath)
    expressionAttributeValues[":emptyList"] = &dynamodb.AttributeValue{L: []*dynamodb.AttributeValue{}}
    expressionAttributeValues[":elem"] = &dynamodb.AttributeValue{L: []*dynamodb.AttributeValue{marshaledElement}}

    return updateDynamoDBElement(dynamoDBClient, tableName, reportID, updateExpression, expressionAttributeValues)
}


func updateDynamoDBElement(dynamoDBClient *dynamodb.DynamoDB, tableName, reportID, updateExpression string, expressionAttributeValues map[string]*dynamodb.AttributeValue) error {
	input := &dynamodb.UpdateItemInput{
        TableName: aws.String(tableName),
        Key: map[string]*dynamodb.AttributeValue{
            constants.ReportIDField.String(): {
                S: aws.String(reportID),
            },
        },
        UpdateExpression:          aws.String(updateExpression),
        ExpressionAttributeValues: expressionAttributeValues,
        ReturnValues:              aws.String("UPDATED_NEW"),
    }

	_, err := dynamoDBClient.UpdateItem(input)
	if err != nil {
		log.Printf("Error updating DynamoDB: %v", err)
	}

	return err
}

func GetItem(tableName, keyName, keyValue string) (item map[string]*dynamodb.AttributeValue, err error) {
	dynamoDBClient, err := newDynamoDBClient(string(constants.USEast2))

	if err != nil {
		return nil, err
	}

	// Create a DynamoDB input structure for the GetItem operation.
	input := &dynamodb.GetItemInput{
		TableName: aws.String(tableName),
		Key: map[string]*dynamodb.AttributeValue{
			keyName: {
				S: aws.String(keyValue),
			},
		},
	}

	result, err := dynamoDBClient.GetItem(input)
	// Execute the GetItem operation.
	if err != nil {
		return nil, fmt.Errorf("error getting item from DynamoDB: %v", err)
	}

	// Check if the item was found.
	if result.Item == nil {
		return nil, nil // Item not found
	}

	return result.Item, nil
}

func GetAllItems(tableName, projectionExpression string) ([]models.Report, error) {
	dynamoDBClient, err := newDynamoDBClient(string(constants.USEast2))

	if err != nil {
		return nil, err
	}

	// Create a DynamoDB ScanInput with the ProjectionExpression
	input := &dynamodb.ScanInput{
		TableName:            aws.String(tableName),
		ProjectionExpression: aws.String(projectionExpression),
	}

	result, err := dynamoDBClient.Scan(input)
	if err != nil {
		return nil, fmt.Errorf("error scanning DynamoDB table: %v", err)
	}

	reports := []models.Report{}

	for _, item := range result.Items {
		var report models.Report
		err = dynamodbattribute.UnmarshalMap(item, &report)
		if err != nil {
			return nil, fmt.Errorf("error unmarshalling DynamoDB item: %v", err)
		}

		reports = append(reports, report)
	}

	return reports, nil
}

func newDynamoDBClient(region string) (*dynamodb.DynamoDB, error) {
	sess, err := session.NewSession(&aws.Config{
		Region: aws.String(region)},
	)
	if err != nil {
		return nil, err
	}
	return dynamodb.New(sess), nil
}
