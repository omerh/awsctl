package cmd

import (
	"fmt"
	"log"

	"github.com/aws/aws-sdk-go/aws"
	"github.com/aws/aws-sdk-go/service/ec2"
	"github.com/omerh/awsctl/pkg/helpers"
	"github.com/spf13/cobra"
)

var cmdEip = &cobra.Command{
	Use:   "eip [region]",
	Short: "Relese unused elastic ip",
	Run: func(cmd *cobra.Command, args []string) {
		// Commnadline arguments
		release, _ := cmd.Flags().GetBool("yes")
		region, err := cmd.Flags().GetString("region")

		if err != nil {
			log.Println(err)
			return
		}
		if region == "all" {
			awsRegions, _ := helpers.GetAllAwsRegions()
			for _, r := range awsRegions {
				manageElasticIPAddresses(r, release)
			}
			return
		}

		if region == "" {
			region = helpers.GetDefaultAwsRegion()
		}

		manageElasticIPAddresses(region, release)
	},
}

func manageElasticIPAddresses(region string, release bool) {
	log.Printf("Running on region: %v", region)
	awsSession, _ := helpers.InitAwsSession(region)
	svc := ec2.New(awsSession)

	input := &ec2.DescribeAddressesInput{}

	result, _ := svc.DescribeAddresses(input)

	for _, address := range result.Addresses {
		if address.AssociationId == nil {
			keep := checkKeepAddressTag(address)
			if keep == true {
				log.Printf("IP Address %s tag was set to keep", *address.PublicIp)
			} else {
				if release == true {
					input := &ec2.ReleaseAddressInput{
						AllocationId: aws.String(*address.AllocationId),
					}
					_, err := svc.ReleaseAddress(input)
					if err != nil {
						log.Println(err)
						return
					}
					log.Printf("IP address %s was released", *address.PublicIp)
				} else {
					log.Printf("IP Address %s, would be released, add --yes to apply", *address.PublicIp)
				}
			}
		}
	}
	fmt.Println("=====================================================================================================")
}

func checkKeepAddressTag(address *ec2.Address) bool {
	for _, tag := range address.Tags {
		if *tag.Key == "keep" && *tag.Value == "true" {
			return true
		}
	}
	return false
}
