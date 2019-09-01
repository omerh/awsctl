package cmd

import "github.com/spf13/cobra"

var (
	getEc2Example = `
	# Get amazon resources
	get ec2 events --region us-east-1
	`
	getEc2Short = ("Get AWS EC2 Details")
)

var getEc2Cmd = &cobra.Command{
	Use:     "ec2",
	Short:   getEc2Short,
	Example: getEc2Example,
	Run: func(cmd *cobra.Command, args []string) {
		if len(args) == 0 {
			cmd.Usage()
		}
	},
}

func init() {
	getEc2Cmd.AddCommand(getEc2Events)
}
