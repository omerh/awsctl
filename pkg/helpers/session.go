package helpers

import (
	"fmt"
	"os"

	"github.com/aws/aws-sdk-go/aws"
	"github.com/aws/aws-sdk-go/aws/session"
)

// var awsSession *session.Session

// InitAwsSession initialize aws session
func InitAwsSession(region string) (*session.Session, error) {

	config := aws.NewConfig().WithRegion(region)
	config = config.WithCredentialsChainVerboseErrors(true)

	// Create session options with SharedConfigState enabled to support AWS SSO profiles
	sessOpts := session.Options{
		Config:            *config,
		SharedConfigState: session.SharedConfigEnable,
	}

	// Check if AWS_PROFILE is set and use it
	if profile := os.Getenv("AWS_PROFILE"); profile != "" {
		sessOpts.Profile = profile
	}

	awsSession, err := session.NewSessionWithOptions(sessOpts)

	if err != nil {
		return nil, fmt.Errorf("error starting a new AWS session: %v", err)
	}

	return awsSession, nil
}
