package cmd

import (
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
		override, _ := cmd.Flags().GetBool("override")

		if region == "all" {
			awsRegions, _ := helpers.GetAllAwsRegions()
			for _, r := range awsRegions {
				// locateNeverExpireCloudwatchlogs(r, retention, apply)
				cloudwatchGroups := helpers.GetCloudwatchGroups(r)
				helpers.SetCloudwatchGroupsExpiry(r, retention, cloudwatchGroups, apply, override)
			}
			return
		}

		if region == "" {
			region = helpers.GetDefaultAwsRegion()
		}
		cloudwatchGroups := helpers.GetCloudwatchGroups(region)
		helpers.SetCloudwatchGroupsExpiry(region, retention, cloudwatchGroups, apply, override)
	},
}
