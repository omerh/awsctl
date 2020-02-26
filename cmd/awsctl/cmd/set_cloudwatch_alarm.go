package cmd

import (
	"github.com/omerh/awsctl/pkg/helpers"
	"github.com/spf13/cobra"
)

var cmdCloudwatchAlarm = &cobra.Command{
	Use:   "cloudwatchalarm [region]",
	Short: "Find and Set cloudwatch alarms to aws resources",
	Run: func(cmd *cobra.Command, args []string) {
		// Command params
		// apply, _ := cmd.Flags().GetBool("yes")

		region, _ := cmd.Flags().GetString("region")

		// if region == "all" {
		// 	awsRegions, _ := helpers.GetAllAwsRegions()
		// 	for _, r := range awsRegions {

		// 	}
		// 	return
		// }

		// if region == "" {
		// 	region = helpers.GetDefaultAwsRegion()
		// }

		helpers.GetAllLmbdasInRegion(region)

	},
}

func init() {
	// Flags
	cmdCloudwatchAlarm.Flags().String("resource", "", "Type of resource to set the alarm")
	cmdCloudwatchAlarm.MarkFlagRequired("resource")
	cmdCloudwatchAlarm.Flags().String("type", "", "Type for alarm to set")
	cmdCloudwatchAlarm.MarkFlagRequired("type")
}
