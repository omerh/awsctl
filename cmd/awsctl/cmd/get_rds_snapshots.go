package cmd

import (
	"github.com/omerh/awsctl/pkg/helpers"
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
			awsRegions, _ := helpers.GetAllAwsRegions()
			for _, r := range awsRegions {
				rdsSnapshotInfoSlice := helpers.GetRDSSnapshots(dbName, rdsType, r, out)
				if out != "json" {
					helpers.PrintRdsSnapshotInformation(rdsSnapshotInfoSlice, region, out)
				}
			}
			return
		}

		if region == "" {
			region = helpers.GetDefaultAwsRegion()
		}

		rdsSnapshotInfoSlice := helpers.GetRDSSnapshots(dbName, rdsType, region, out)
		if out != "json" {
			helpers.PrintRdsSnapshotInformation(rdsSnapshotInfoSlice, region, out)
		}
	},
}

func init() {
	getRdsSnapshots.Flags().StringP("name", "n", "", "resourceId name")
	getRdsSnapshots.Flags().StringP("type", "t", "instance", "instance/cluster")
}
