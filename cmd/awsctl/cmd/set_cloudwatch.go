package cmd

import (
	"fmt"
	"log"

	"github.com/aws/aws-sdk-go/aws"
	"github.com/aws/aws-sdk-go/service/cloudwatchlogs"
	"github.com/omerh/awsctl/pkg/helpers"
	"github.com/spf13/cobra"
)

var cmdCloudwatch = &cobra.Command{
	Use:   "cloudwatch [region]",
	Short: "Find and Set cloudwatch logs with no retention policy set to never and set them to 14 days",
	Run: func(cmd *cobra.Command, args []string) {
		// Command params
		apply, _ := cmd.Flags().GetBool("yes")
		retention, _ := cmd.Flags().GetInt64("retention")
		region, _ := cmd.Flags().GetString("region")

		if region == "all" {
			awsRegions, _ := helpers.GetAllAwsRegions()
			for _, r := range awsRegions {
				locateNeverExpireCloudwatchlogs(r, retention, apply)
			}
			return
		}

		if region == "" {
			region = helpers.GetDefaultAwsRegion()
		}
		locateNeverExpireCloudwatchlogs(region, retention, apply)
	},
}

func locateNeverExpireCloudwatchlogs(region string, retention int64, apply bool) {
	log.Printf("Running on region: %v", region)
	awsSession, _ := helpers.InitAwsSession(region)
	svc := cloudwatchlogs.New(awsSession)

	// default empty input beofre retriving next tokens
	input := &cloudwatchlogs.DescribeLogGroupsInput{}
	// get all cloudwatch log groups
	result, err := svc.DescribeLogGroups(input)

	if err != nil {
		log.Println(err)
		return
	}

	cloudwatchGroups := result.LogGroups

	for result.NextToken != nil {
		input = &cloudwatchlogs.DescribeLogGroupsInput{
			NextToken: result.NextToken,
		}
		result, err = svc.DescribeLogGroups(input)
		if err != nil {
			log.Println(err)
			return
		}
		for _, group := range result.LogGroups {
			cloudwatchGroups = append(cloudwatchGroups, group)
		}
	}

	// Calculate the region total size for cost
	var totalLogByteSize int64
	noRetentionSet := false

	for _, group := range cloudwatchGroups {
		if group.RetentionInDays == nil {
			totalLogByteSize = totalLogByteSize + *group.StoredBytes
			noRetentionSet = true
			if apply == true {
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

	// helpers.GetAwsServiceCost()

	if noRetentionSet == true {
		log.Printf("Region %s total log size:", region)
		log.Printf("Total log size in with no retention policy is: %v bytes", totalLogByteSize)
	}

	fmt.Println("=====================================================================================================")
}
