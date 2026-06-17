package cmd

import "github.com/spf13/cobra"

var historyCmd = &cobra.Command{
	Use:   "history",
	Short: "Transaction history",
	Long:  `View and manage transaction history.`,
}

func init() {
	rootCmd.AddCommand(historyCmd)
}
