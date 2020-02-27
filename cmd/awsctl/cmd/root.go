package cmd

import (
	"fmt"
	"os"

	"github.com/spf13/cobra"
)

var (
	rootCommand = `
		# List all regions
		awsctl list regions regions
		# List all availability zones in a specific region
		awsctl list azs --region eu-west-3
		# Delete all unused EBS volumes on all aws regions
		awsctl delete ebs --region all --yes
		# Create multiple EBS volumes
		awsctl create ebs --count 1 --size 10 --region us-east-1 -z us-east-1a --yes
		# Set cloudwatch logs with retention policy never to retention policy to 14 days on all regions
		awsctl set cloudwatch --region all --retention 14 --yes
	`
)

var rootCmd = &cobra.Command{
	Use:     "awsctl",
	Short:   "awsctl for managing aws infrastructure",
	Example: rootCommand,
	Run: func(cmd *cobra.Command, args []string) {
		if len(args) == 0 {
			cmd.Help()
			os.Exit(0)
		}
	},
}

// Persistent flags goes here
func init() {
	// Commands
	rootCmd.AddCommand(setCmd)
	rootCmd.AddCommand(getCmd)

	// Flags
	rootCmd.PersistentFlags().StringP("region", "r", "", "aws region/all")
	rootCmd.PersistentFlags().StringP("out", "o", "text", "Output text/json")
	rootCmd.PersistentFlags().Bool("slack", false, "send custom webhook slack message for monitor")
	// rootCmd.PersistentFlags().StringP("type", "t", "instance", "instance/cluster")
}

// Execute using cobra command
func Execute() {
	if err := rootCmd.Execute(); err != nil {
		fmt.Println(err)
		os.Exit(1)
	}
}
