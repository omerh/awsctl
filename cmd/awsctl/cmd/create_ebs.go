package cmd

import (
	"fmt"
	"log"
	"strings"

	"github.com/aws/aws-sdk-go/aws"
	"github.com/aws/aws-sdk-go/service/ec2"
	"github.com/omerh/awsctl/pkg/helper"
	"github.com/spf13/cobra"
)

var cmdCreateEbs = &cobra.Command{
	Use:   "ebs --size 1 --count 1 --region region --availability-zone az --yes",
	Short: "Create a bunch of EBS volume",
	Run: func(cmd *cobra.Command, Args []string) {
		create, _ := cmd.Flags().GetBool("yes")
		region, _ := cmd.Flags().GetString("region")
		size, _ := cmd.Flags().GetInt64("size")
		count, _ := cmd.Flags().GetInt("count")
		az, _ := cmd.Flags().GetString("availability-zone")
		volumeType, _ := cmd.Flags().GetString("volume-type")

		if count < 1 {
			fmt.Println("EBS count must be greater than 0")
			return
		}

		if size < 1 {
			fmt.Println("EBS volume size must be greater than 0")
			return
		}

		if az == "" || !strings.Contains(az, region) {
			fmt.Println("EBS Avialability Zone must be declared")
			azs, _ := helper.GetAllAwsAzs(region)
			fmt.Printf("Possible azs are %v", azs)
			return
		}

		if region == "all" {
			awsRegions, _ := helper.GetAllAwsRegions()
			for _, r := range awsRegions {
				createEbsVolumes(r, size, count, create, az, volumeType)
			}
			return
		}

		if region == "" {
			region = helper.GetDefaultAwsRegion()
		}

		createEbsVolumes(region, size, count, create, az, volumeType)
	},
}

func init() {
	cmdCreateEbs.Flags().Int64P("size", "s", 0, "EBS Volume size")
	cmdCreateEbs.MarkFlagRequired("size")
	cmdCreateEbs.Flags().IntP("count", "c", 0, "EBS count")
	cmdCreateEbs.MarkFlagRequired("count")
	cmdCreateEbs.Flags().StringP("availability-zone", "z", "", "Availability zone")
	cmdCreateEbs.MarkFlagRequired("availability-zone")
	cmdCreateEbs.Flags().StringP("volume-type", "t", "gp2", "EBS Volume type (gp2,io1,st1,sc1,standart)")
}

func createEbsVolumes(region string, size int64, count int, create bool, az string, volumeType string) {
	log.Printf("Running on region: %v", region)
	awsSession, _ := helper.InitAwsSession(region)
	svc := ec2.New(awsSession)

	input := &ec2.CreateVolumeInput{
		AvailabilityZone: aws.String(az),
		Size:             aws.Int64(size),
		VolumeType:       aws.String(volumeType),
	}

	for i := 0; i < count; i++ {
		if create == true {
			result, err := svc.CreateVolume(input)
			if err != nil {
				log.Println(err)
				return
			}
			log.Printf("Creating EBS Volume %v, type: %v, size: %v GiB, in az: %v ", *result.VolumeId, *result.VolumeType, *result.Size, *result.AvailabilityZone)
		} else {
			log.Printf("Would create volume %v, pass command with --yes", i+1)
		}
	}
}
