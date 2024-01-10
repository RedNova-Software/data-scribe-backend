package util

import (
    "github.com/aws/aws-sdk-go/aws"
    "github.com/aws/aws-sdk-go/aws/session"
    "github.com/aws/aws-sdk-go/service/dynamodb"
    "github.com/aws/aws-sdk-go/service/dynamodb/dynamodbattribute"
	"api/shared/constants"
    "fmt"
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

    return updateDynamoDBNestedElement(dynamoDBClient, tableName, reportID, updateExpression, expressionAttributeValues)
}


func updateDynamoDBNestedElement(dynamoDBClient *dynamodb.DynamoDB, tableName, reportID, updateExpression string, expressionAttributeValues map[string]*dynamodb.AttributeValue) error {
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


func newDynamoDBClient(region string) (*dynamodb.DynamoDB, error) {
    sess, err := session.NewSession(&aws.Config{
        Region: aws.String(region)},
    )
    if err != nil {
        return nil, err
    }
    return dynamodb.New(sess), nil
}


