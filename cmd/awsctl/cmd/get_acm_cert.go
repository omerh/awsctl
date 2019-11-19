package cmd

import (
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
	// Get all certificates
	certificateList := helper.GetAcmCertificates(region)

	// Get each certificate information
	for _, certificate := range certificateList {
		certificateInfo := helper.DescribeAcmCertificate(region, *certificate.CertificateArn)
		helper.CheckCertificateStatus(certificateInfo, region)
	}
}
