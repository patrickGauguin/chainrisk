package cmd

import (
	"fmt"

	"github.com/spf13/cobra"
)

var scanCmd = &cobra.Command{
	Use:   "scan [repo-url]",
	Short: "Scan a GitHub repository",
	Run: func(cmd *cobra.Command, args []string) {
		fmt.Println("scanning: " + args[0])
	},
}

func init() {
	rootCmd.AddCommand(scanCmd)
}
