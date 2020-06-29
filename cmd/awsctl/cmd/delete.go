package cmd

import (
	"github.com/spf13/cobra"
)

var (
	deleteExample = `
	# Delete unused ebs in region us-east-1
	delete ebs --region us-east-1
	# Delete unused eip in all regions
	awsctl delete eip -r all
	`
	deleteShort = ("Delete AWS Resource")
)

var deleteCmd = &cobra.Command{
	Use:     "delete",
	Short:   deleteShort,
	Example: deleteExample,
	Run: func(cmd *cobra.Command, args []string) {
		if len(args) == 0 {
			cmd.Usage()
		}
	},
}

func init() {
	rootCmd.AddCommand(deleteCmd)
	deleteCmd.AddCommand(cmdEbs)
	deleteCmd.AddCommand(cmdEip)
	deleteCmd.AddCommand(deleteRdsSnapshots)
	deleteCmd.AddCommand(cmdDeleteEcrImages)
	deleteCmd.AddCommand(deleteAcmCertCmd)
	deleteCmd.AddCommand(cmdDeleteNI)
	deleteCmd.PersistentFlags().BoolP("yes", "y", false, "Specify --yes to execute")
}
