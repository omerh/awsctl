package cmd

import (
	"flag"
	"fmt"
	"log"

	"github.com/omerh/awsctl/pkg/helper"

	"github.com/spf13/cobra"
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

		t := helper.GetAwsServiceCost()

		log.Println(&t)

		return nil
	},
}

func init() {
	cmdCheck.Flags().StringP("name", "n", "", "Whats your name")
}
