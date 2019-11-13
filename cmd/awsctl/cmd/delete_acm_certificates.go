package cmd

import (
	"github.com/omerh/awsctl/pkg/helper"
	"github.com/spf13/cobra"
)

var deleteAcmCertCmd = &cobra.Command{
	Use:   "certificates",
	Short: "Delete All unused certificated",
	Run: func(cmd *cobra.Command, Args []string) {
		region, _ := cmd.Flags().GetString("region")
		apply, _ := cmd.Flags().GetBool("yes")

		if region == "all" {
			awsRegions, _ := helper.GetAllAwsRegions()
			for _, r := range awsRegions {
				helper.DeleteUnusedAcmCertificates(r, apply)
			}
			return
		}

		// No region arg passed
		if region == "" {
			region = helper.GetDefaultAwsRegion()
		}
		helper.DeleteUnusedAcmCertificates(region, apply)
	},
}
