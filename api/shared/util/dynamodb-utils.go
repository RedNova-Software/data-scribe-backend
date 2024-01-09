package util

import (
    "github.com/aws/aws-sdk-go/aws"
    "github.com/aws/aws-sdk-go/aws/session"
    "github.com/aws/aws-sdk-go/service/dynamodb"
    "github.com/aws/aws-sdk-go/service/dynamodb/dynamodbattribute"
	"api/shared/constants"
)

func PutItem(item interface{}, tableName string) error {
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


func newDynamoDBClient(region string) (*dynamodb.DynamoDB, error) {
    sess, err := session.NewSession(&aws.Config{
        Region: aws.String(region)},
    )
    if err != nil {
        return nil, err
    }
    return dynamodb.New(sess), nil
}


