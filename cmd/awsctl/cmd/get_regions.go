package cmd

import (
	"fmt"
	"log"

	"github.com/omerh/awsctl/pkg/helpers"
	"github.com/omerh/awsctl/pkg/outputs"
	"github.com/spf13/cobra"
)

var getRegionsCmd = &cobra.Command{
	Use:   "regions",
	Short: "Get all regions",
	Run: func(cmd *cobra.Command, Args []string) {
		regions, _ := helpers.GetAllAwsRegions()
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
