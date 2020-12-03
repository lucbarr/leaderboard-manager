package cmd

import (
	"fmt"
	"os"

	"github.com/spf13/cobra"
)

var rootCmd = &cobra.Command{
	Use:   "leaderboard-manager",
	Short: "leaderboard-manager is a leaderboard manager (:",
	Run: func(cmd *cobra.Command, args []string) {
	},
}

// Execute runs the command
func Execute() {
	if err := rootCmd.Execute(); err != nil {
		fmt.Println(err)
		os.Exit(1)
	}
}
