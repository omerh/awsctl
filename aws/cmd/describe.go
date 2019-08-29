package cmd

import (
	"github.com/spf13/cobra"
)

var (
	describeExample = `
	# Describe ebs volume
	describe ebs --region us-east-1
	`
	describeShort = ("Describe AWS Resource")
)

var describeCmd = &cobra.Command{
	Use:     "describe",
	Short:   describeShort,
	Example: describeExample,
	Run: func(cmd *cobra.Command, args []string) {
		if len(args) == 0 {
			cmd.Usage()
		}
	},
}

func init() {
	rootCmd.AddCommand(describeCmd)
}
