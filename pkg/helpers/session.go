package helpers

import (
	"fmt"

	"github.com/aws/aws-sdk-go/aws"
	"github.com/aws/aws-sdk-go/aws/session"
)

var awsSession *session.Session

// InitAwsSession initialize aws session
func InitAwsSession(region string) (*session.Session, error) {

	config := aws.NewConfig().WithRegion(region)
	config = config.WithCredentialsChainVerboseErrors(true)

	awsSession, err := session.NewSession(config)

	if err != nil {
		return nil, fmt.Errorf("error starting a new AWS session: %v", err)
	}

	return awsSession, nil
}
