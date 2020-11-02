package cmd

import (
	"fmt"

	"github.com/omerh/awsctl/pkg/helpers"
	"github.com/spf13/cobra"
)

var (
	getRiExample = `
	# Get amazon resources
	get ri --status active --region us-east-1
	`
	getRihort = ("Get AWS EC2 Reservation Details")
)

var getRiCmd = &cobra.Command{
	Use:     "ri",
	Short:   getRihort,
	Example: getRiExample,
	Run: func(cmd *cobra.Command, args []string) {
		region, _ := cmd.Flags().GetString("region")

		if region == "all" {
			awsRegions, _ := helpers.GetAllAwsRegions()
			for _, r := range awsRegions {
				riSummary, ec2Summary := generateRIreport(r)
				utilizationReport(riSummary, ec2Summary, r)
			}
			return
		}

		if region == "" {
			region = helpers.GetDefaultAwsRegion()
		}

		riSummary, ec2Summary := generateRIreport(region)
		utilizationReport(riSummary, ec2Summary, region)
	},
}

func generateRIreport(region string) (map[string]int64, map[string]int64) {
	riSummary := helpers.GetAllReservations(region, "active")
	allEC2 := helpers.GetAllEC2Instances(region, "ondemand", "running")
	ec2Summary := helpers.SummeriesEC2Instances(allEC2)
	return riSummary, ec2Summary
}

func utilizationReport(ri map[string]int64, ec2 map[string]int64, region string) {
	fmt.Printf("Reservation status in %v (Only for running instances)\n", region)
	fmt.Println("=================================================================")

	for instanceType, instanceCount := range ri {
		util := utilization(ec2[instanceType], instanceCount)
		fmt.Printf("- You have %v %v instances with %v active resevation, utilizing %v %% of reserved instances\n", ec2[instanceType], instanceType, instanceCount, util)
	}
	fmt.Println()
}

func utilization(instanceCount int64, riCount int64) int {
	if instanceCount >= riCount {
		return 100
	}
	util := int(float64(instanceCount) / float64(riCount) * 100)
	return util
}
