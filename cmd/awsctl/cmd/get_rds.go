package cmd

import (
	"github.com/omerh/awsctl/pkg/helpers"
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
			awsRegions, _ := helpers.GetAllAwsRegions()
			for _, r := range awsRegions {
				helpers.GetAllRds(r, rdsType, out)
			}
			return
		}

		if region == "" {
			region = helpers.GetDefaultAwsRegion()
		}

		helpers.GetAllRds(region, rdsType, out)
	},
}

func init() {
	getRdsCmd.Flags().StringP("type", "t", "instance", "instance/cluster")
}
