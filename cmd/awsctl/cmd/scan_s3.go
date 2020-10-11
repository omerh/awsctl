package cmd

import (
	"log"

	"github.com/omerh/awsctl/pkg/helpers"
	"github.com/spf13/cobra"
)

type s3bucket struct {
	name          string
	region        string
	encryptionSet bool
}

var scanS3 = &cobra.Command{
	Use:     "s3",
	Short:   "Scan S3 bucket for misconfigurations and security vulnerabilities",
	Example: "scan s3 --region us-east-1",
	Run: func(cmd *cobra.Command, Args []string) {
		var s3buckets []s3bucket
		buckets := helpers.GetAllS3Buckets()
		for _, bucket := range buckets {
			// get region bucket for addional meta
			region := helpers.GetS3BucketLocation(*bucket.Name)
			// get s3 server side encryption
			encryptioSet := helpers.CheckBucketEncryption(*bucket.Name, region)
			// get s3 public access
			helpers.GetS3PuclicAccess(*bucket.Name, region)
			bucket := s3bucket{
				name:          *bucket.Name,
				region:        region,
				encryptionSet: encryptioSet,
			}
			s3buckets = append(s3buckets, bucket)
		}
		log.Println(s3buckets[0])
	},
}
