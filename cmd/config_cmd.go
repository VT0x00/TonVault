package cmd

import "github.com/spf13/cobra"

var configCmd = &cobra.Command{
	Use:   "config",
	Short: "Manage configuration",
	Long:  `View, set, and reset configuration values.`,
}

func init() {
	rootCmd.AddCommand(configCmd)
}
