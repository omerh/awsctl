package cmd

import (
	"github.com/omerh/awsctl/pkg/helper"
	"github.com/spf13/cobra"
)

var getRdsCmd = &cobra.Command{
	Use:   "rds",
	Short: "Get All RDS Instaces",
	Run: func(cmd *cobra.Command, Args []string) {
		region, _ := cmd.Flags().GetString("region")
		out, _ := cmd.Flags().GetString("out")

		if region == "all" {
			awsRegions, _ := helper.GetAllAwsRegions()
			for _, r := range awsRegions {
				helper.GetAllRdsInstances(r, out)
			}
			return
		}

		if region == "" {
			region = helper.GetDefaultAwsRegion()
		}

		helper.GetAllRdsInstances(region, out)
	},
}

func init() {
	getRdsCmd.Flags().StringP("region", "r", "us-east-1", "Aws region")
	getRdsCmd.MarkFlagRequired("region")
}
