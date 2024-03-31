package cmd

import (
	"github.com/aws/aws-sdk-go/service/ecr"
	"github.com/omerh/awsctl/pkg/helpers"
	"github.com/spf13/cobra"
)

var cmdDeleteEcrImages = &cobra.Command{
	Use:   "ecr",
	Short: "Delete ecr docker images with how much images to keep",
	Run: func(cmd *cobra.Command, args []string) {
		apply, _ := cmd.Flags().GetBool("yes")
		keep, _ := cmd.Flags().GetInt("keep")
		region, _ := cmd.Flags().GetString("region")

		if region == "all" {
			awsRegions, _ := helpers.GetAllAwsRegions()
			for _, r := range awsRegions {
				repos := helpers.GetECRRepositories(r)
				getEcrImagesAndDeleteOld(repos, region, keep, apply)
			}
			return
		}

		// No region arg passed
		if region == "" {
			region = helpers.GetDefaultAwsRegion()
		}

		repos := helpers.GetECRRepositories(region)
		getEcrImagesAndDeleteOld(repos, region, keep, apply)
	},
}

func init() {
	cmdDeleteEcrImages.Flags().Int("keep", 999, "How many images to keep on the ECR repo")
	cmdDeleteEcrImages.MarkFlagRequired("keep")
}

func getEcrImagesAndDeleteOld(repos []*ecr.Repository, region string, keep int, apply bool) {
	for _, repo := range repos {
		imageDetails, deleteImages := helpers.EcrDescribeImages(*repo.RepositoryName, region, keep)
		if deleteImages > 0 {
			// needs to be deleted
			var digest []string
			for _, imageDetail := range imageDetails {
				digest = append(digest, *imageDetail.ImageDigest)
			}
			// log.Printf("%v", digest)
			helpers.DeleteEcrImages(*repo.RepositoryName, digest, region, apply)
		}
	}
}
