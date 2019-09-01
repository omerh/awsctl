package cmd

import (
	"github.com/spf13/cobra"
)

var (
	createExample = `
	# Create EBS in region us-east-1
	create ebs --size 1 --count 1  --region us-east-1
	`
	createShort = ("Create AWS Resource")
)

var createCmd = &cobra.Command{
	Use:     "create ebs [--region]",
	Short:   createShort,
	Example: createExample,
	Run: func(cmd *cobra.Command, args []string) {
		if len(args) == 0 {
			cmd.Usage()
		}
	},
}

func init() {
	rootCmd.AddCommand(createCmd)
	createCmd.AddCommand(cmdCreateEbs)
}
