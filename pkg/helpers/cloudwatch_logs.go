package helpers

import (
	"fmt"
	"log"
	"os"

	"github.com/aws/aws-sdk-go/aws"
	"github.com/aws/aws-sdk-go/service/cloudwatchlogs"
)

// GetCloudwatchGroups getting all cloudwatch groups in a region
func GetCloudwatchGroups(region string) []*cloudwatchlogs.LogGroup {
	log.Printf("Running on region: %v", region)
	awsSession, _ := InitAwsSession(region)
	svc := cloudwatchlogs.New(awsSession)

	// default empty input before retrieving next tokens
	input := &cloudwatchlogs.DescribeLogGroupsInput{}
	// get all cloudwatch log groups
	result, err := svc.DescribeLogGroups(input)

	if err != nil {
		log.Println(err)
		os.Exit(1)
	}

	cloudwatchGroups := result.LogGroups

	for result.NextToken != nil {
		input = &cloudwatchlogs.DescribeLogGroupsInput{
			NextToken: result.NextToken,
		}
		result, err = svc.DescribeLogGroups(input)
		if err != nil {
			log.Println(err)
			os.Exit(1)
		}
		// for _, group := range result.LogGroups {
		// 	cloudwatchGroups = append(cloudwatchGroups, group)
		// }
		cloudwatchGroups = append(cloudwatchGroups, result.LogGroups...)
	}
	return cloudwatchGroups
}

func bytesToMB(bytes uint64) float64 {
	const bytesInMB = 1024 * 1024
	return float64(bytes) / float64(bytesInMB)
}

// SetCloudwatchGroupsExpiry Set expiry on a cloudwatch group
func SetCloudwatchGroupsExpiry(region string, retention int64, cloudwatchGroups []*cloudwatchlogs.LogGroup, apply bool, override bool) {
	awsSession, _ := InitAwsSession(region)
	svc := cloudwatchlogs.New(awsSession)

	var totalLogByteSize int64
	noRetentionSet := false
	for _, group := range cloudwatchGroups {
		if group.RetentionInDays == nil || override {
			totalLogByteSize = totalLogByteSize + *group.StoredBytes
			noRetentionSet = true
			if apply && override {
				// set input filter
				input := &cloudwatchlogs.PutRetentionPolicyInput{
					LogGroupName:    aws.String(*group.LogGroupName),
					RetentionInDays: aws.Int64(retention),
				}
				// put retention policy
				_, err := svc.PutRetentionPolicy(input)
				if err != nil {
					log.Println(err)
					return
				}
				log.Printf("Retention policy for %s was set to %v", *group.LogGroupName, retention)
			} else {
				log.Printf("Group %s retention policy would be set to %d (size is %v Bytes), --yes to apply", *group.LogGroupName, retention, *group.StoredBytes)
			}
		}
	}
	if noRetentionSet {
		log.Printf("Region %s total log size:", region)
		totalLogByteSizeInMB := bytesToMB(uint64(totalLogByteSize))
		// log.Println(totalLogByteSize)
		log.Printf("Total log size in with no retention policy is: %.2f MB", totalLogByteSizeInMB)
	}

	fmt.Println("=====================================================================================================")
}
