package cmd

import (
	"fmt"

	"github.com/spf13/cobra"
)

const version = "v0.0.2"

func init() {
	rootCmd.AddCommand(versionCmd)
}

var versionCmd = &cobra.Command{
	Use:   "version",
	Short: "Print version of awsctl",
	Run: func(cmd *cobra.Command, args []string) {
		fmt.Printf("awsctl version %v", version)
	},
}
