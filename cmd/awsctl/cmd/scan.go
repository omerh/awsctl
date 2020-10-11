package cmd

import "github.com/spf13/cobra"

var (
	scanExample = `
  # Scan amazon resources
  awsctl scan s3 --region us-east-1
	
  # Scan amzon security groups
  awsctl scan asg --region us-east-1
`
	// scanShort = ("Scan AWS Resource")
)

var scanCmd = &cobra.Command{
	Use:     "scan",
	Short:   "Scan AWS Resource",
	Example: scanExample,
	Run: func(cmd *cobra.Command, args []string) {
		if len(args) == 0 {
			cmd.Usage()
		}
	},
}

func init() {
	scanCmd.AddCommand(scanS3)
}
