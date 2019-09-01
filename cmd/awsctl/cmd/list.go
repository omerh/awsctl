package cmd

import "github.com/spf13/cobra"

var (
	listExample = `
	# List amazon resources
	list aws regions
	list aws azs --region us-east-1
	`
	listShort = ("List AWS Resources")
)

var listCmd = &cobra.Command{
	Use:     "list",
	Short:   listShort,
	Example: listExample,
	Run: func(cmd *cobra.Command, args []string) {
		if len(args) == 0 {
			cmd.Usage()
		}
	},
}

func init() {
	rootCmd.AddCommand(listCmd)
	listCmd.AddCommand(listRegionsCmd)
	listCmd.AddCommand(listAzsCmd)
}
