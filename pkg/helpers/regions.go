package helpers

import (
	"fmt"
	"log"
	"os"

	"github.com/aws/aws-sdk-go/aws"
	"github.com/aws/aws-sdk-go/aws/session"
	"github.com/aws/aws-sdk-go/service/ec2"
)

var allRegions []string

// GetAllAwsRegions will retrive all aws regions
//
func GetAllAwsRegions() ([]string, error) {
	// log.Println("Getting all regions...")
	request := &ec2.DescribeRegionsInput{}
	awsRegion := os.Getenv("AWS_REGION")
	if awsRegion == "" {
		awsRegion = "us-east-1"
	}
	config := aws.NewConfig().WithRegion(awsRegion)
	sess, err := session.NewSession(config)
	if err != nil {
		return allRegions, fmt.Errorf("Error starting a new AWS session: %v", err)
	}

	client := ec2.New(sess, config)

	response, err := client.DescribeRegions(request)
	if err != nil {
		return allRegions, fmt.Errorf("Got an error while querying for valid regions (verify your AWS credentials?): %v", err)
	}

	for _, region := range response.Regions {
		allRegions = append(allRegions, *region.RegionName)
	}

	return allRegions, nil
}

// GetDefaultAwsRegion resolve deafult region
//
func GetDefaultAwsRegion() (region string) {
	region = os.Getenv("AWS_REGION")
	if region == "" {
		log.Println("No region is set in environment, please set AWS_REGION environment variable or pass --region")
		os.Exit(1)
	} else {
		log.Printf("Using default region %s from environment variable AWS_REGION", region)
	}
	return region
}
