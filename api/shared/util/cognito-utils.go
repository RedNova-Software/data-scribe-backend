package util

import (
	"sync"

	"github.com/aws/aws-sdk-go/aws"
	"github.com/aws/aws-sdk-go/aws/session"
	"github.com/aws/aws-sdk-go/service/cognitoidentityprovider"
)

var (
	cognitoClient *cognitoidentityprovider.CognitoIdentityProvider
	cognitoOnce   sync.Once
	cognitoErr    error
)

// GetCognitoClient returns a singleton CognitoIdentityProvider client
func GetCognitoClient(region string) (*cognitoidentityprovider.CognitoIdentityProvider, error) {
	cognitoOnce.Do(func() {
		cognitoClient, cognitoErr = newCognitoClient(region)
	})
	return cognitoClient, cognitoErr
}

func newCognitoClient(region string) (*cognitoidentityprovider.CognitoIdentityProvider, error) {
	sess, err := session.NewSession(&aws.Config{
		Region: aws.String(region)},
	)
	if err != nil {
		return nil, err
	}
	return cognitoidentityprovider.New(sess), nil
}
