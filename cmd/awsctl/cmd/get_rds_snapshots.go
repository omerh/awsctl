package cmd

import (
	"github.com/omerh/awsctl/pkg/helper"
	"github.com/spf13/cobra"
)

var getRdsSnapshotsExample = `
	awsctl get rdssnapshots -t cluster -r us-east-1 -n dbclustername
	awsctl get rdssnapshots -t instance -r us-east-1 -n dbinstancename -o json
`

var getRdsSnapshots = &cobra.Command{
	Use:     "rdssnapshots",
	Short:   "Get RDS snapshots",
	Example: getRdsSnapshotsExample,
	Run: func(cmd *cobra.Command, args []string) {
		region, _ := cmd.Flags().GetString("region")
		out, _ := cmd.Flags().GetString("out")
		rdsType, _ := cmd.Flags().GetString("type")
		dbName, _ := cmd.Flags().GetString("name")

		if region == "all" {
			awsRegions, _ := helper.GetAllAwsRegions()
			for _, r := range awsRegions {
				helper.GetRDSSnapshots(dbName, rdsType, r, out)
			}
			return
		}

		if region == "" {
			region = helper.GetDefaultAwsRegion()
		}

		helper.GetRDSSnapshots(dbName, rdsType, region, out)
	},
}

func init() {
	getRdsSnapshots.Flags().StringP("name", "n", "", "resourceId name")
}
