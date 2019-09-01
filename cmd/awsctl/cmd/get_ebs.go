package cmd

import "github.com/spf13/cobra"

var (
	getEbsExample = `
	# Get amazon resources
	get ebs get ebs encryption --type unencrypted --region us-east-1 --out json
	`
	getEbsShort = ("Get AWS Resources")
)

var getEbsCmd = &cobra.Command{
	Use:     "ebs",
	Short:   getEbsShort,
	Example: getEbsExample,
	Run: func(cmd *cobra.Command, args []string) {
		if len(args) == 0 {
			cmd.Usage()
		}
	},
}

func init() {
	getEbsCmd.AddCommand(getEbsEncryptionCmd)
}
