package cmd

import (
	"github.com/omerh/awsctl/pkg/helper"
	"github.com/spf13/cobra"
)

var getRdsCmd = &cobra.Command{
	Use:   "rds",
	Short: "Get RDS Instaces or clusters",
	Run: func(cmd *cobra.Command, Args []string) {
		region, _ := cmd.Flags().GetString("region")
		out, _ := cmd.Flags().GetString("out")
		rdsType, _ := cmd.Flags().GetString("type")

		if region == "all" {
			awsRegions, _ := helper.GetAllAwsRegions()
			for _, r := range awsRegions {
				helper.GetAllRds(r, rdsType, out)
			}
			return
		}

		if region == "" {
			region = helper.GetDefaultAwsRegion()
		}

		helper.GetAllRds(region, rdsType, out)
	},
}

func init() {
	getRdsCmd.Flags().StringP("type", "t", "instance", "instance/cluster")
}
