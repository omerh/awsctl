package cmd

import (
	"log"

	"github.com/omerh/awsctl/pkg/helper"
	"github.com/spf13/cobra"

	"github.com/aws/aws-sdk-go/aws"
	"github.com/aws/aws-sdk-go/service/ec2"
)

var cmdEbs = &cobra.Command{
	Use:   "ebs [region]",
	Short: "Locate available ebs and delete them",
	Run: func(cmd *cobra.Command, args []string) {
		delete, _ := cmd.Flags().GetBool("yes")
		filter, _ := cmd.Flags().GetString("filter")
		region, _ := cmd.Flags().GetString("region")

		if region == "all" {
			awsRegions, _ := helper.GetAllAwsRegions()
			for _, r := range awsRegions {
				seekAvailableVolumes(r, delete, filter)
			}
			return
		}

		if region == "" {
			region = helper.GetDefaultAwsRegion()
		}

		seekAvailableVolumes(region, delete, filter)
	},
}

func init() {
	cmdEbs.Flags().StringP("filter", "f", "available", "EBS status (creating|available|in-use|deleting|deleted|error)")
}

func seekAvailableVolumes(region string, delete bool, filter string) {
	log.Printf("Running on region: %v", region)
	awsSession, _ := helper.InitAwsSession(region)
	svc := ec2.New(awsSession)

	awsVolumeFilters := []*ec2.Filter{
		{
			Name: aws.String("status"),
			Values: []*string{
				aws.String(filter),
			},
		},
	}

	input := &ec2.DescribeVolumesInput{
		Filters: awsVolumeFilters,
	}

	result, _ := svc.DescribeVolumes(input)

	availableVolumeList := result.Volumes

	for result.NextToken != nil {
		input = &ec2.DescribeVolumesInput{
			NextToken: result.NextToken,
			Filters:   awsVolumeFilters,
		}
		result, _ = svc.DescribeVolumes(input)
		for _, volume := range result.Volumes {
			availableVolumeList = append(availableVolumeList, volume)
		}
	}

	// Filtering out volumes
	var filteredVolumeList []string

	// Filtered volume size
	var filteredVolumeSize int64

	for _, volume := range availableVolumeList {
		if volume.Tags != nil {
			keep := checkVolumeKeepTag(volume)
			if keep == true {
				log.Printf("Volume %s tag was tagged to keep", *volume.VolumeId)
			} else {
				filteredVolumeList = append(filteredVolumeList, *volume.VolumeId)
				filteredVolumeSize = filteredVolumeSize + *volume.Size
			}
		} else {
			filteredVolumeList = append(filteredVolumeList, *volume.VolumeId)
			filteredVolumeSize = filteredVolumeSize + *volume.Size
		}
	}

	if len(filteredVolumeList) > 0 {
		log.Printf("Found %v %s EBS volumes to delete, Total volumes size is: %d GiB", len(filteredVolumeList), filter, filteredVolumeSize)
	} else {
		log.Printf("No machine EBS volumes in state %s found", filter)
	}

	// Delete ebs volumes
	if len(filteredVolumeList) > 0 {
		for _, volume := range filteredVolumeList {
			if delete == true {
				log.Printf("Deleting volume %v", volume)
				input := &ec2.DeleteVolumeInput{
					VolumeId: aws.String(volume),
				}
				_, err := svc.DeleteVolume(input)
				if err != nil {
					log.Printf("Failed to delete volume %v", volume)
					return
				}
			} else {
				log.Printf("Would delete volume %v, add --yes to apply", volume)
			}
		}
	}
	log.Println("==================================================================")
}

func checkVolumeKeepTag(volume *ec2.Volume) bool {
	for _, tag := range volume.Tags {
		if *tag.Key == "keep" && *tag.Value == "true" {
			return true
		}
	}
	return false
}
