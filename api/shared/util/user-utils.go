package util

import (
	"api/shared/constants"
	"api/shared/models"
	"fmt"
	"os"

	"github.com/aws/aws-lambda-go/events"
	"github.com/aws/aws-sdk-go/aws"
	"github.com/aws/aws-sdk-go/service/cognitoidentityprovider"
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
	// Create a new Cognito Identity Provider client
	client, err := GetCognitoClient(constants.USEast2)
	if err != nil {
		return "", err
	}

	// Prepare the request
	input := &cognitoidentityprovider.AdminGetUserInput{
		UserPoolId: aws.String(os.Getenv(constants.UserPoolID)),
		Username:   aws.String(userID),
	}

	// Fetch the user details
	result, err := client.AdminGetUser(input)
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

func GetAllUsers() ([]models.User, error) {
	client, err := GetCognitoClient(constants.USEast2)
	if err != nil {
		return nil, err
	}

	input := &cognitoidentityprovider.ListUsersInput{
		UserPoolId: aws.String(os.Getenv(constants.UserPoolID)),
	}

	var users []models.User
	err = client.ListUsersPages(input, func(page *cognitoidentityprovider.ListUsersOutput, lastPage bool) bool {
		for _, user := range page.Users {
			var userID, userNickName string

			for _, attr := range user.Attributes {
				if *attr.Name == constants.CognitoAttrSub {
					userID = *attr.Value
				} else if *attr.Name == constants.CognitoAttrNickName {
					userNickName = *attr.Value
				}
			}

			users = append(users, models.User{
				UserID:       userID,
				UserNickName: userNickName,
			})
		}
		return !lastPage
	})

	if err != nil {
		return nil, err
	}

	return users, nil
}
