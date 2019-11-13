package cmd

import (
	"fmt"
	"log"

	"github.com/aws/aws-sdk-go/service/ecr"
	"github.com/omerh/awsctl/pkg/helper"
	"github.com/spf13/cobra"
)

var setEcrRegistryPolicyCmd = &cobra.Command{
	Use:   "ecrregistrypolicy",
	Short: "Set ECR Registry policy to untagged images",
	Run: func(cmd *cobra.Command, args []string) {
		apply, _ := cmd.Flags().GetBool("yes")
		retention, _ := cmd.Flags().GetInt("retention")
		region, _ := cmd.Flags().GetString("region")

		if region == "all" {
			awsRegions, _ := helper.GetAllAwsRegions()
			for _, r := range awsRegions {
				repos := helper.GetECRRepositories(r)
				checkAndSetEcrSliceLifecyclePolicy(repos, retention, r, apply)
			}
			return
		}

		// No region arg passed
		if region == "" {
			region = helper.GetDefaultAwsRegion()
		}

		repos := helper.GetECRRepositories(region)
		checkAndSetEcrSliceLifecyclePolicy(repos, retention, region, apply)
	},
}

func checkAndSetEcrSliceLifecyclePolicy(repos []*ecr.Repository, retention int, region string, apply bool) {
	for _, repo := range repos {
		policySet := helper.CheckECRRepositoryLifecyclePolicy(*repo.RepositoryName, region)
		if policySet == false {
			// Need to set lifecycle policy
			if apply == true {
				helper.SetEcrRepositoryLifecyclePolicy(*repo.RepositoryName, retention, region)
			} else {
				log.Printf("Will set retention to %v for %v days, pass --yes to execute the command", *repo.RepositoryName, retention)
			}
		}
	}
	fmt.Println("=====================================================================================================")
}
