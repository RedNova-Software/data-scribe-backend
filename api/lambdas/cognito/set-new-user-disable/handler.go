package main

import (
	"context"
	"fmt"

	"github.com/aws/aws-lambda-go/events"
	"github.com/aws/aws-lambda-go/lambda"
	"github.com/aws/aws-sdk-go/aws"
	"github.com/aws/aws-sdk-go/aws/session"
	"github.com/aws/aws-sdk-go/service/cognitoidentityprovider"
)

func Handler(ctx context.Context, event events.CognitoEventUserPoolsHeader) (events.CognitoEventUserPoolsHeader, error) {
	fmt.Printf("Event: %+v\n", event)

	if event.TriggerSource == "PostConfirmation_ConfirmSignUp" {
		sess := session.Must(session.NewSession())
		cognitoClient := cognitoidentityprovider.New(sess)

		params := &cognitoidentityprovider.AdminDisableUserInput{
			UserPoolId: aws.String(event.UserPoolID),
			Username:   aws.String(event.UserName),
		}

		_, err := cognitoClient.AdminDisableUser(params)
		if err != nil {
			fmt.Printf("Error disabling user: %v", err)
			return event, err // Return the error to indicate failure
		}

		fmt.Println("User disabled successfully")
	}
	return event, nil
}

func main() {
	lambda.Start(Handler)
}
