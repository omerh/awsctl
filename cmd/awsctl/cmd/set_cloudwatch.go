package cmd

import (
	"github.com/omerh/awsctl/pkg/helpers"
	"github.com/spf13/cobra"
)

var cmdCloudwatch = &cobra.Command{
	Use:   "cloudwatch",
	Short: "Find and Set cloudwatch logs with no retention policy set to never and set them to 14 days",
	Run: func(cmd *cobra.Command, args []string) {
		// Command params
		apply, _ := cmd.Parent().PersistentFlags().GetBool("yes")
		retention, _ := cmd.Parent().PersistentFlags().GetInt64("retention")
		region, _ := cmd.Root().PersistentFlags().GetString("region")
		overwrite, _ := cmd.Parent().PersistentFlags().GetBool("overwrite")

		if region == "all" {
			awsRegions, _ := helpers.GetAllAwsRegions()
			for _, r := range awsRegions {
				// locateNeverExpireCloudwatchlogs(r, retention, apply)
				cloudwatchGroups := helpers.GetCloudwatchGroups(r)
				helpers.SetCloudwatchGroupsExpiry(r, retention, cloudwatchGroups, apply, overwrite)
			}
			return
		}

		if region == "" {
			region = helpers.GetDefaultAwsRegion()
		}
		cloudwatchGroups := helpers.GetCloudwatchGroups(region)
		helpers.SetCloudwatchGroupsExpiry(region, retention, cloudwatchGroups, apply, overwrite)
	},
}
