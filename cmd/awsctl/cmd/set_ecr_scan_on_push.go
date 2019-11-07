package cmd

import (
	"fmt"
	"log"

	"github.com/aws/aws-sdk-go/service/ecr"
	"github.com/omerh/awsctl/pkg/helper"
	"github.com/spf13/cobra"
)

var setEcrScanOnPushcmd = &cobra.Command{
	Use:   "ecrscanonpush",
	Short: "Set ecr repository configuration to scan on push",
	Run: func(cmd *cobra.Command, args []string) {
		apply, _ := cmd.Flags().GetBool("yes")
		region, _ := cmd.Flags().GetString("region")
		scanOnPush, _ := cmd.Flags().GetBool("scan")

		if region == "all" {
			awsRegions, _ := helper.GetAllAwsRegions()
			for _, r := range awsRegions {
				repos := helper.GetECRRepositories(r)
				setEcrRepositoryConfigurationScanOnPush(repos, scanOnPush, r, apply)
			}
			return
		}
		// No region arg passed
		if region == "" {
			region = helper.GetDefaultAwsRegion()
		}
		repos := helper.GetECRRepositories(region)
		setEcrRepositoryConfigurationScanOnPush(repos, scanOnPush, region, apply)
	},
}

func setEcrRepositoryConfigurationScanOnPush(repos []*ecr.Repository, scanOnPush bool, region string, apply bool) {
	for _, repo := range repos {
		if apply == true {
			helper.SetEcrRepoImageScanOnPush(*repo.RepositoryName, region, scanOnPush)
		} else {
			log.Printf("will set scanOnPush to %v for repository %v, pass --yes to execute the command", scanOnPush, *repo.RepositoryName)
		}

	}
	fmt.Println("=====================================================================================================")
}

func init() {
	setEcrScanOnPushcmd.Flags().Bool("scan", true, "Set repository configuration for ScanOnPush")
}
