package cmd

import (
	"os"

	"github.com/spf13/cobra"
)

var rootCmd = &cobra.Command{
	Use:   "chainrisk",
	Short: "Software supply chain risk analyzer",
}

func Execute() {
	err := rootCmd.Execute()

	if err != nil {
		println(err.Error())
		os.Exit(1)
	}
}
