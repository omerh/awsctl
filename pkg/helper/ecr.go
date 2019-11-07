package helper

import (
	"fmt"
	"log"

	"github.com/aws/aws-sdk-go/aws"
	"github.com/aws/aws-sdk-go/service/ecr"
)

// GetECRRepositories from aws region
func GetECRRepositories(region string) []*ecr.Repository {
	log.Printf("Running on region: %v", region)
	awsSession, _ := InitAwsSession(region)
	svc := ecr.New(awsSession)
	input := &ecr.DescribeRepositoriesInput{}
	result, _ := svc.DescribeRepositories(input)

	registeriesSlice := result.Repositories

	// iterate over NextToken to retrive all repositories from Ecr in the region
	for result.NextToken != nil {
		input := &ecr.DescribeRepositoriesInput{
			NextToken: result.NextToken,
		}
		result, _ := svc.DescribeRepositories(input)
		for _, registry := range result.Repositories {
			registeriesSlice = append(registeriesSlice, registry)
		}
	}
	return registeriesSlice
}

// CheckECRRepositoryLifecyclePolicy for a repository in a region
func CheckECRRepositoryLifecyclePolicy(repositoryName string, region string) bool {
	awsSession, _ := InitAwsSession(region)
	svc := ecr.New(awsSession)
	input := &ecr.GetLifecyclePolicyInput{
		RepositoryName: aws.String(repositoryName),
	}
	_, err := svc.GetLifecyclePolicy(input)
	if err != nil {
		// Ecr lifecyclePolicy is not set
		// log.Printf("error: %v", err)
		return false
	}

	return true
}

// SetEcrRepositoryLifecyclePolicy set the life time policy
func SetEcrRepositoryLifecyclePolicy(repositoryName string, days int, region string) {
	log.Printf("Setting lifecycle policy to %v for %v days", repositoryName, days)
	awsSession, _ := InitAwsSession(region)
	svc := ecr.New(awsSession)
	
	// tagStatus is hardcoded to prevent mistakes like deleting latest that is not maintained
	lifecyclePolicyText := fmt.Sprintf(`{
		"rules": [
			{
				"rulePriority": 1,
				"description": "Expire untagged images by awsctl",
				"selection": {
					"tagStatus": "untagged",
					"countType": "sinceImagePushed",
					"countUnit": "days",
					"countNumber": %v
				},
				"action": {
					"type": "expire"
				}
			}
		]
	  }`, days)

	input := &ecr.PutLifecyclePolicyInput{
		RepositoryName:      aws.String(repositoryName),
		LifecyclePolicyText: aws.String(lifecyclePolicyText),
	}

	_, err := svc.PutLifecyclePolicy(input)

	if err != nil {
		log.Printf("There was a problem setting lifecycle policy to %v", repositoryName)
		log.Println(err)
	}
}