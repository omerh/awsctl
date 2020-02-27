package cmd

import (
	"fmt"

	"github.com/omerh/awsctl/pkg/helpers"
	"github.com/spf13/cobra"
)

var getCloudwatchAlarmCmd = &cobra.Command{
	Use:   "cloudwatchalarm",
	Short: "Get cloudwatch alarms",
	Run: func(cmd *cobra.Command, Args []string) {
		region, _ := cmd.Flags().GetString("region")
		// out, _ := cmd.Flags().GetString("out")

		if region == "all" {
			awsRegions, _ := helpers.GetAllAwsRegions()
			for _, r := range awsRegions {
				alarms := helpers.ListCloudwatchAlarms(r)
				fmt.Println(alarms)
			}
			return
		}

		if region == "" {
			region = helpers.GetDefaultAwsRegion()
		}

		alarms := helpers.ListCloudwatchAlarms(region)
		fmt.Println(alarms)
	},
}
