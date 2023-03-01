package cmd

import (
	"github.com/spf13/cobra"
)

var rootCmd = &cobra.Command{
	Use:   "quoter",
	Short: "Quoter CLI app",
	Long:  `Quoter CLI app`,
}

func Execute() error {
	return rootCmd.Execute()
}
