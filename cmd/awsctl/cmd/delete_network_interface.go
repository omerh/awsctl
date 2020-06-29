package cmd

import (
	"log"

	"github.com/aws/aws-sdk-go/service/ec2"
	"github.com/omerh/awsctl/pkg/helpers"
	"github.com/spf13/cobra"
)

var cmdDeleteNI = &cobra.Command{
	Use:   "ni",
	Short: "Delete network interface",
	Run: func(cmd *cobra.Command, args []string) {
		region, _ := cmd.Flags().GetString("region")
		filter, _ := cmd.Flags().GetString("filter")
		apply, _ := cmd.Flags().GetBool("yes")

		if filter != "available" && filter != "in-use" {
			log.Println("filter can be abailable or in-use")
			return
		}

		if region == "all" {
			awsRegions, _ := helpers.GetAllAwsRegions()
			for _, r := range awsRegions {
				networkInterfacesSlice := helpers.GetAllElasticNetworkInterfaces(r, filter)
				deleteNetworkInterface(networkInterfacesSlice, r, apply)
			}
			return
		}
		// No region arg passed
		if region == "" {
			region = helpers.GetDefaultAwsRegion()
		}
		networkInterfacesSlice := helpers.GetAllElasticNetworkInterfaces(region, filter)
		deleteNetworkInterface(networkInterfacesSlice, region, apply)
	},
}

func deleteNetworkInterface(networkInterfacesSlice []*ec2.NetworkInterface, region string, apply bool) {
	for _, n := range networkInterfacesSlice {
		ok := helpers.DeleteNetworkInterface(region, *n.NetworkInterfaceId, apply)
		if !ok {
			log.Println("Failed to delete")
		}
	}
}

func init() {
	cmdDeleteNI.Flags().String("filter", "available", "status: available, in-use")
}
