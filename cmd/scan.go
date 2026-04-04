package cmd

import (
	"fmt"
	"os"

	"github.com/patrickGauguin/chainrisk/internal/report"
	"github.com/patrickGauguin/chainrisk/internal/scanner"
	"github.com/spf13/cobra"
)

var scanCmd = &cobra.Command{
	Use:   "scan [repo-url]",
	Short: "Scan a GitHub repository",
	Run: func(cmd *cobra.Command, args []string) {
		token := os.Getenv("GITHUB_TOKEN")
		s := scanner.New(token)

		result, err := s.Scan(args[0])
		if err != nil {
			fmt.Println("Error:", err)
			return
		}

		report.PrintTerminal(result)
	},
}

func init() {
	rootCmd.AddCommand(scanCmd)
}
