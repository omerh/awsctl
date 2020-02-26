package cmd

import (
	"github.com/omerh/awsctl/pkg/helpers"
	"github.com/spf13/cobra"
)

var deleteAcmCertCmd = &cobra.Command{
	Use:   "certificates",
	Short: "Delete All unused certificated",
	Run: func(cmd *cobra.Command, Args []string) {
		region, _ := cmd.Flags().GetString("region")
		apply, _ := cmd.Flags().GetBool("yes")

		if region == "all" {
			awsRegions, _ := helpers.GetAllAwsRegions()
			for _, r := range awsRegions {
				helpers.DeleteUnusedAcmCertificates(r, apply)
			}
			return
		}

		// No region arg passed
		if region == "" {
			region = helpers.GetDefaultAwsRegion()
		}
		helpers.DeleteUnusedAcmCertificates(region, apply)
	},
}
