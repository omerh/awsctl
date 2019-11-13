package cmd

import (
	"fmt"

	"github.com/omerh/awsctl/pkg/helper"
	"github.com/spf13/cobra"
)

var getAcmCertCmd = &cobra.Command{
	Use:   "certificates",
	Short: "Get All expiring certificated that are being used",
	Run: func(cmd *cobra.Command, Args []string) {
		region, _ := cmd.Flags().GetString("region")

		if region == "all" {
			awsRegions, _ := helper.GetAllAwsRegions()
			for _, r := range awsRegions {
				getCertifcateInformation(r)
			}
			return
		}

		// No region arg passed
		if region == "" {
			region = helper.GetDefaultAwsRegion()
		}
		getCertifcateInformation(region)

	},
}

func getCertifcateInformation(region string) {
	certificateList := helper.GetAcmCertificates(region)
	for _, certificate := range certificateList {
		t := helper.DescribeAcmCertificate(region, *certificate.CertificateArn)
		fmt.Println(*t.Certificate.DomainName)
	}
}
