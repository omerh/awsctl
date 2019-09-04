package helper

import (
	"fmt"

	"github.com/aws/aws-sdk-go/service/ecr"
)

// GetECRRepositories from aws region
func GetECRRepositories(region string, out string) {
	awsSession, _ := InitAwsSession(region)
	svc := ecr.New(awsSession)
	input := &ecr.DescribeRepositoriesInput{}
	response, _ := svc.DescribeRepositories(input)

	fmt.Println(response)
}
