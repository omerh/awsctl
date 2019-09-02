package cmd

import (
	"github.com/omerh/awsctl/pkg/helper"
	"github.com/spf13/cobra"
)

var getRdsCmdExmaple = `
  awsctl get rds -t instance -r us-east-1
  awsctl get rds -t cluster -r us-east-1
`

var getRdsCmd = &cobra.Command{
	Use:     "rds",
	Short:   "Get RDS Instaces or clusters",
	Example: getRdsCmdExmaple,
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
