package cmd

import (
	"fmt"
	"os"

	"github.com/patrickGauguin/chainrisk/internal/github"
	"github.com/spf13/cobra"
)

var scanCmd = &cobra.Command{
	Use:   "scan [repo-url]",
	Short: "Scan a GitHub repository",
	Run: func(cmd *cobra.Command, args []string) {
		token := os.Getenv("GITHUB_TOKEN")
		client := github.New(token)

		owner, repo, err := github.ParseOwnerRepo(args[0])

		if err != nil {
			fmt.Println("Error:", err)
			return
		}

		repoInfo, err := client.GetRepo(owner, repo)

		if err != nil {
			fmt.Println("Error:", err)
			return
		}

		fmt.Printf("Name: 		%s\n", repoInfo.FullName)
		fmt.Printf("Language: 	%s\n", repoInfo.Language)
		fmt.Printf("Stars: 		%d\n", repoInfo.Stars)
		fmt.Printf("Archived: 	%v\n", repoInfo.Archived)
		fmt.Printf("Pushed: 	%s\n", repoInfo.LastPushed.Format("2006-01-02"))
	},
}

func init() {
	rootCmd.AddCommand(scanCmd)
}
