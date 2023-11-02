package cmd

import (
	"github.com/spf13/cobra"
)

var (
	setExample = `
	# Set cloudwatch log groups in region us-east-1 to 7 days
	set cloudwatch --region us-east-1 --retention 7
	`
	setShort = ("Set AWS Resource")
)

var setCmd = &cobra.Command{
	Use:     "set cloudwatch [--region] [--retention]",
	Short:   setShort,
	Example: setExample,
	Run: func(cmd *cobra.Command, args []string) {
		if len(args) == 0 {
			cmd.Usage()
		}
	},
}

func init() {
	// Commands
	setCmd.AddCommand(cmdCloudwatch)
	setCmd.AddCommand(setEcrRegistryPolicyCmd)
	setCmd.AddCommand(setEcrScanOnPushcmd)
	setCmd.AddCommand((cmdCloudwatchAlarm))

	// Flags
	setCmd.PersistentFlags().BoolP("yes", "y", false, "Specify --yes to execute")
	setCmd.PersistentFlags().Bool("override", false, "Specify --override to for cloudwatch override")
	setCmd.PersistentFlags().Int64("retention", 14, "Retention in days")
}
