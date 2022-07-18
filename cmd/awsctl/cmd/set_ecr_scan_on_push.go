package cmd

import (
	"fmt"
	"log"

	"github.com/aws/aws-sdk-go/service/ecr"
	"github.com/omerh/awsctl/pkg/helpers"
	"github.com/spf13/cobra"
)

var setEcrScanOnPushcmd = &cobra.Command{
	Use:   "ecrscanonpush",
	Short: "Set ecr repository configuration to scan on push",
	Run: func(cmd *cobra.Command, args []string) {
		apply, _ := cmd.Flags().GetBool("yes")
		region, _ := cmd.Flags().GetString("region")
		scanOnPush, _ := cmd.Flags().GetBool("scan")
		fmt.Println(scanOnPush)
		if region == "all" {
			awsRegions, _ := helpers.GetAllAwsRegions()
			for _, r := range awsRegions {
				repos := helpers.GetECRRepositories(r)
				setEcrRepositoryConfigurationScanOnPush(repos, scanOnPush, r, apply)
			}
			return
		}
		// No region arg passed
		if region == "" {
			region = helpers.GetDefaultAwsRegion()
		}
		repos := helpers.GetECRRepositories(region)
		setEcrRepositoryConfigurationScanOnPush(repos, scanOnPush, region, apply)
	},
}

func setEcrRepositoryConfigurationScanOnPush(repos []*ecr.Repository, scanOnPush bool, region string, apply bool) {
	for _, repo := range repos {
		if apply {
			helpers.SetEcrRepoImageScanOnPush(*repo.RepositoryName, region, scanOnPush)
		} else {
			log.Printf("will set scanOnPush to %v for repository %v, pass --yes to execute the command", scanOnPush, *repo.RepositoryName)
		}

	}
	fmt.Println("=====================================================================================================")
}

func init() {
	setEcrScanOnPushcmd.Flags().String("scan", "true", "Set repository configuration for ScanOnPush")
}
