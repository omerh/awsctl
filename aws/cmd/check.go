package cmd

import (
	"flag"
	"fmt"
	"log"

	"github.com/spf13/cobra"
	"gitlab.com/omerh/awsctl/aws/helper"
)

var cmdCheck = &cobra.Command{
	Use: "check",
	RunE: func(cmd *cobra.Command, args []string) error {
		fmt.Println("In check")
		flag.CommandLine.Parse([]string{})

		name, _ := cmd.Flags().GetString("name")

		if name == "" {
			name = "Missing"
		}

		region, _ := cmd.Flags().GetString("region")
		if region == "" {
			region = "missing"
		}

		// fmt.Println(name)
		// fmt.Println(region)

		// cmd.PersistentFlags().StringP("region", "r", "AWS Region", "Please use the region")

		t := helper.GetAwsServiceCost()

		log.Println(&t)

		// fmt.Println(val)

		return nil
	},
}

func init() {
	// rootCmd.AddCommand(cmdCheck)
	cmdCheck.Flags().StringP("name", "n", "", "Whats your name")
	// cmdCheck.MarkFlagRequired("name")
	// rootCmd.PersistentFlags().arra
}
