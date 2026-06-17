package cmd

import (
	"fmt"

	"github.com/spf13/cobra"
)

var historyShowCmd = &cobra.Command{
	Use:   "show <tx_hash>",
	Short: "Show transaction details",
	Long:  `Display detailed information about a specific transaction.`,
	Args:  cobra.ExactArgs(1),
	RunE: func(cmd *cobra.Command, args []string) error {
		fmt.Println("Transaction details: coming soon")
		return nil
	},
}

func init() {
	historyCmd.AddCommand(historyShowCmd)
}
