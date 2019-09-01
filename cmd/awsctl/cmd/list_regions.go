package cmd

import (
	"fmt"
	"log"

	"github.com/omerh/awsctl/pkg/helper"
	"github.com/omerh/awsctl/pkg/outputs"
	"github.com/spf13/cobra"
)

var listRegionsCmd = &cobra.Command{
	Use:   "regions",
	Short: "List all regions",
	Run: func(cmd *cobra.Command, Args []string) {
		regions, _ := helper.GetAllAwsRegions()
		out, _ := cmd.Flags().GetString("out")

		switch out {
		case "json":
			outputs.PrintGenericJSONOutput(regions, "")
		default:
			log.Println("Listing all regions...")
			for _, region := range regions {
				fmt.Println(region)
			}
		}
	},
}
