package cmd

import (
	"github.com/omerh/awsctl/pkg/helper"
	"github.com/spf13/cobra"
)

var getEc2Events = &cobra.Command{
	Use:     "events",
	Short:   "List all ec2 instances events in a region",
	Example: "events -r us-east-1",
	Run: func(cmd *cobra.Command, Args []string) {
		region, _ := cmd.Flags().GetString("region")
		out, _ := cmd.Flags().GetString("out")

		if region == "all" {
			awsRegions, _ := helper.GetAllAwsRegions()
			for _, r := range awsRegions {
				helper.GetAllEc2Events(r, out)
			}
			return
		}

		if region == "" {
			region = helper.GetDefaultAwsRegion()
		}

		helper.GetAllEc2Events(region, out)
	},
}

func init() {
	getEc2Events.Flags().StringP("region", "r", "us-east-1", "Aws region")
	getEc2Events.MarkFlagRequired("region")
}
