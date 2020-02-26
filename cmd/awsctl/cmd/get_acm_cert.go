package cmd

import (
	"github.com/omerh/awsctl/pkg/helpers"
	"github.com/spf13/cobra"
)

var getAcmCertCmd = &cobra.Command{
	Use:   "certificates",
	Short: "Get All expiring certificated that are being used",
	Run: func(cmd *cobra.Command, Args []string) {
		region, _ := cmd.Flags().GetString("region")

		if region == "all" {
			awsRegions, _ := helpers.GetAllAwsRegions()
			for _, r := range awsRegions {
				getCertifcateInformation(r)
			}
			return
		}

		// No region arg passed
		if region == "" {
			region = helpers.GetDefaultAwsRegion()
		}
		getCertifcateInformation(region)

	},
}

func getCertifcateInformation(region string) {
	// Get all certificates
	certificateList := helpers.GetAcmCertificates(region)

	// Get each certificate information
	for _, certificate := range certificateList {
		certificateInfo := helpers.DescribeAcmCertificate(region, *certificate.CertificateArn)
		helpers.CheckCertificateStatus(certificateInfo, region)
	}
}
