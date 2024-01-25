package util

import (
	"sync"

	"github.com/aws/aws-sdk-go/aws"
	"github.com/aws/aws-sdk-go/aws/session"
	"github.com/aws/aws-sdk-go/service/dynamodb"
)

var (
	client    *dynamodb.DynamoDB
	once      sync.Once
	createErr error
)

// Use this to get client as it returns a singleton
func GetDynamoDBClient(region string) (*dynamodb.DynamoDB, error) {
	once.Do(func() {
		client, createErr = newDynamoDBClient(region)
	})
	return client, createErr
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
