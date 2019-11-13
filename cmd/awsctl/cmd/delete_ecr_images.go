package cmd

import (
	"github.com/aws/aws-sdk-go/service/ecr"
	"github.com/omerh/awsctl/pkg/helper"
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
			awsRegions, _ := helper.GetAllAwsRegions()
			for _, r := range awsRegions {
				repos := helper.GetECRRepositories(r)
				getEcrImagesAndDeleteOld(repos, region, keep, apply)
			}
			return
		}

		// No region arg passed
		if region == "" {
			region = helper.GetDefaultAwsRegion()
		}

		repos := helper.GetECRRepositories(region)
		getEcrImagesAndDeleteOld(repos, region, keep, apply)
	},
}

func init() {
	cmdDeleteEcrImages.Flags().Int("keep", 999, "How many images to keep on the ECR repo")
	cmdDeleteEcrImages.MarkFlagRequired("keep")
}

func getEcrImagesAndDeleteOld(repos []*ecr.Repository, region string, keep int, apply bool) {
	for _, repo := range repos {
		helper.EcrDeleteOldImageBuildsAndKeep(*repo.RepositoryName, region, keep)
	}
}
