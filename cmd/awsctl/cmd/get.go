package cmd

import "github.com/spf13/cobra"

var (
	getExample = `
	# Get amazon resources
	get ec2 events --region us-east-1
	`
	getShort = ("Get AWS Resource")
)

var getCmd = &cobra.Command{
	Use:     "get",
	Short:   getShort,
	Example: getExample,
	Run: func(cmd *cobra.Command, args []string) {
		if len(args) == 0 {
			cmd.Usage()
		}
	},
}

func init() {
	getCmd.AddCommand(getEc2Cmd)
	getCmd.AddCommand(getEbsCmd)
	getCmd.AddCommand(getRegionsCmd)
	getCmd.AddCommand(getAzsCmd)
	getCmd.AddCommand(getRdsCmd)
	rootCmd.AddCommand(getCmd)
}
