package cmd

import (
	"fmt"
	"log"

	"github.com/aws/aws-sdk-go/aws"
	"github.com/aws/aws-sdk-go/service/ec2"
	"github.com/omerh/awsctl/pkg/helpers"
	"github.com/omerh/awsctl/pkg/outputs"
	"github.com/spf13/cobra"
)

var getEbsEncryptionCmd = &cobra.Command{
	Use:     "encryption",
	Short:   "get all ebs encryption types all,encrypted,unencrypted",
	Example: "encryption --type all --region all",
	Run: func(cmd *cobra.Command, args []string) {
		region, _ := cmd.Flags().GetString("region")
		encryptionFilter, _ := cmd.Flags().GetString("type")
		out, _ := cmd.Flags().GetString("out")

		if region == "all" {
			awsRegions, _ := helpers.GetAllAwsRegions()
			for _, r := range awsRegions {
				getEbsEncryption(r, encryptionFilter, out)
			}
			return
		}

		if region == "" {
			region = helpers.GetDefaultAwsRegion()
		}

		getEbsEncryption(region, encryptionFilter, out)
	},
}

func init() {
	getEbsEncryptionCmd.Flags().StringP("type", "t", "all", "all/encrypted/unencrypted")
}

func getEbsEncryption(region string, encryptionFilter string, out string) {
	awsSession, _ := helpers.InitAwsSession(region)
	svc := ec2.New(awsSession)

	encryptionFilterString := string(encryptionFilter)

	if encryptionFilter == "encrypted" {
		encryptionFilterString = "true"
	} else if encryptionFilter == "unencrypted" {
		encryptionFilterString = "false"
	}

	input := &ec2.DescribeVolumesInput{}

	if encryptionFilterString != "all" {
		awsVolumeFilters := []*ec2.Filter{
			{
				Name: aws.String("encrypted"),
				Values: []*string{
					aws.String(encryptionFilterString),
				},
			},
		}
		input = &ec2.DescribeVolumesInput{
			Filters: awsVolumeFilters,
		}
	}

	result, _ := svc.DescribeVolumes(input)

	switch out {
	case "json":
		outputs.PrintGenericJSONOutput(result.Volumes, region)
		// v, _ := json.Marshal(result.Volumes)
		// fmt.Println(string(v))
	default:
		log.Printf("Running on region: %v", region)
		for _, ebs := range result.Volumes {
			log.Println("-----------------------------------------------------------------------------------")
			if ebs.Attachments != nil {
				if *ebs.Encrypted {
					log.Printf("Found Ebs %v with encryption set to %v", *ebs.VolumeId, *ebs.Encrypted)
					log.Printf("with the key: %v", *ebs.KmsKeyId)
					log.Printf("Its attached to %v", *ebs.Attachments[0].InstanceId)
				} else {
					log.Printf("Found Ebs %v with encryption set to %v", *ebs.VolumeId, *ebs.Encrypted)
					log.Printf("Its attached to %v", *ebs.Attachments[0].InstanceId)
					log.Println("You have a running instance with no disk encryption")
				}
			} else {
				if *ebs.Encrypted {
					log.Printf("Found Ebs %v with encryption set to %v", *ebs.VolumeId, *ebs.Encrypted)
					log.Printf("with the key: %v", *ebs.KmsKeyId)
				} else {
					log.Printf("Found Ebs %v with encryption set to %v", *ebs.VolumeId, *ebs.Encrypted)
				}
			}
		}
		fmt.Println("=====================================================================================================")
	}
}
