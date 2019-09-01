package cmd

import (
	"fmt"

	"github.com/omerh/awsctl/pkg/helper"
	"github.com/omerh/awsctl/pkg/hooks"
	"github.com/omerh/awsctl/pkg/outputs"
	"github.com/spf13/cobra"
)

var listAzsCmd = &cobra.Command{
	Use:   "azs",
	Short: "List all azs in a regions",
	Run: func(cmd *cobra.Command, Args []string) {
		region, _ := cmd.Flags().GetString("region")
		azs, _ := helper.GetAllAwsAzs(region)
		out, _ := cmd.Flags().GetString("out")
		slack, _ := cmd.Flags().GetBool("slack")

		switch out {
		case "json":
			outputs.PrintGenericJSONOutput(azs, region)
			// fmt.Println(azs)
		default:
			fmt.Printf("Available AWS azs in region %v: \n", region)
			for _, az := range azs {
				fmt.Println(az)
			}
		}
		if slack == true {
			hooks.SendSlackWebhook("test")
		}
	},
}

func init() {
	listAzsCmd.Flags().StringP("region", "r", "us-east-1", "Aws region")
	listAzsCmd.MarkFlagRequired("region")
}
