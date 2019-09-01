package cmd

import (
	"log"

	"github.com/aws/aws-sdk-go/aws"
	"github.com/aws/aws-sdk-go/service/ec2"
	"github.com/omerh/awsctl/pkg/helper"
	"github.com/omerh/awsctl/pkg/outputs"
	"github.com/spf13/cobra"
)

var describeEbsCmd = &cobra.Command{
	Use:   "ebs",
	Short: "Describe an EBS volume",
	Run: func(cmd *cobra.Command, args []string) {
		volumeID, _ := cmd.Flags().GetString("id")
		region, _ := cmd.Flags().GetString("region")
		out, _ := cmd.Flags().GetString("out")

		if region == "" {
			region = helper.GetDefaultAwsRegion()
		}

		describeEbsVolumeByID(volumeID, region, out)
	},
}

func init() {
	describeCmd.AddCommand(describeEbsCmd)
	describeEbsCmd.Flags().String("id", "", "EBS Volume ID")
	describeEbsCmd.MarkFlagRequired("id")
	describeEbsCmd.Flags().StringP("region", "r", "", "aws region")
	describeEbsCmd.MarkFlagRequired("region")
}

func describeEbsVolumeByID(volumeID string, region string, out string) {
	awsSession, _ := helper.InitAwsSession(region)
	svc := ec2.New(awsSession)

	awsVolumeFilters := []*ec2.Filter{
		{
			Name: aws.String("volume-id"),
			Values: []*string{
				aws.String(volumeID),
			},
		},
	}

	input := &ec2.DescribeVolumesInput{
		Filters: awsVolumeFilters,
	}

	result, err := svc.DescribeVolumes(input)

	if err != nil {
		log.Println(err)
		return
	}

	switch out {
	case "json":
		outputs.PrintGenericJSONOutput(result.Volumes, region)
	default:
		// for _, r := range result.Volumes {

		// }
		outputs.PrintGenericTextOutput(result.Volumes, region, "Volume Information:")
	}
}
