package util

import (
	"api/shared/constants"
	"context"
	"fmt"
	"os"

	"github.com/aws/aws-lambda-go/events"
	"github.com/aws/aws-sdk-go-v2/aws"
	"github.com/aws/aws-sdk-go-v2/config"
	"github.com/aws/aws-sdk-go-v2/service/cognitoidentityprovider"
)

// ExtractUserID extracts the user's Cognito ID from the API Gateway request context
func ExtractUserID(request events.APIGatewayProxyRequest) (string, error) {
	claims, ok := request.RequestContext.Authorizer["claims"].(map[string]interface{})
	if !ok {
		return "", fmt.Errorf("no claims found in request")
	}

	userID, ok := claims["sub"].(string)
	if !ok {
		return "", fmt.Errorf("user ID not found in claims")
	}

	return userID, nil
}

// GetUserNickname fetches the nickname of the user from Cognito User Pool
func GetUserNickname(userID string) (string, error) {
	// Load the AWS default config
	cfg, err := config.LoadDefaultConfig(context.TODO())
	if err != nil {
		return "", err
	}

	// Create a new Cognito Identity Provider client
	client := cognitoidentityprovider.NewFromConfig(cfg)

	// Prepare the request
	input := &cognitoidentityprovider.AdminGetUserInput{
		UserPoolId: aws.String(os.Getenv(constants.UserPoolID)),
		Username:   aws.String(userID),
	}

	// Fetch the user details
	result, err := client.AdminGetUser(context.Background(), input)
	if err != nil {
		return "", err
	}

	// Loop through the user attributes to find the nickname
	for _, attr := range result.UserAttributes {
		if *attr.Name == "nickname" {
			return *attr.Value, nil
		}
	}

	return "", fmt.Errorf("nickname not found for user: " + userID)
}
