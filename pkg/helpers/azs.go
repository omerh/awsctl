package helpers

import (
	"fmt"

	"github.com/aws/aws-sdk-go/aws"
	"github.com/aws/aws-sdk-go/aws/session"
	"github.com/aws/aws-sdk-go/service/ec2"
)

var allAzs []string

// GetAllAwsAzs retriva all available azs in a region
//
func GetAllAwsAzs(region string) (azs []string, err error) {
	input := &ec2.DescribeAvailabilityZonesInput{}
	config := aws.NewConfig().WithRegion(region)
	sess, err := session.NewSession(config)
	if err != nil {
		return allRegions, fmt.Errorf("Error starting a new AWS session: %v", err)
	}
	svc := ec2.New(sess, config)
	result, _ := svc.DescribeAvailabilityZones(input)

	for _, r := range result.AvailabilityZones {
		azs = append(azs, *r.ZoneName)
	}
	return azs, err
}
