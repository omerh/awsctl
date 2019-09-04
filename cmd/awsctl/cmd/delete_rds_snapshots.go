package cmd

import (
	"github.com/omerh/awsctl/pkg/helper"
	"github.com/spf13/cobra"
)

var deleteRdsSnapshots = &cobra.Command{
	Use:     "rdssnapshots",
	Short:   "Delete RDS snapshot by date and identifier",
	Example: "Example",
	Run: func(cmd *cobra.Command, args []string) {
		region, _ := cmd.Flags().GetString("region")
		out, _ := cmd.Flags().GetString("out")
		rdsType, _ := cmd.Flags().GetString("type")
		dbName, _ := cmd.Flags().GetString("name")
		older, _ := cmd.Flags().GetInt("older")
		apply, _ := cmd.Flags().GetBool("yes")

		var rdsSnapshot []helper.RdsSnapshotInfo
		if region == "all" {
			awsRegions, _ := helper.GetAllAwsRegions()
			for _, r := range awsRegions {
				rdsSnapshot = helper.GetRDSSnapshots(dbName, rdsType, r, out)
				helper.DeleteRdsSnapshots(rdsSnapshot, older, r, apply, rdsType, out)
			}
			return
		}

		if region == "" {
			region = helper.GetDefaultAwsRegion()
		}

		rdsSnapshot = helper.GetRDSSnapshots(dbName, rdsType, region, out)
		helper.DeleteRdsSnapshots(rdsSnapshot, older, region, apply, rdsType, out)
	},
}

func init() {
	deleteRdsSnapshots.Flags().Int("older", 365, "snapshot older than n days")
	deleteRdsSnapshots.MarkFlagRequired("older")
	deleteRdsSnapshots.Flags().StringP("name", "n", "", "resourceId name")
}
