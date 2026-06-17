package cmd

import "github.com/spf13/cobra"

var balanceCmd = &cobra.Command{
	Use:   "balance",
	Short: "Check balances",
	Long:  `View TON and Jetton balances for your wallets.`,
}

func init() {
	rootCmd.AddCommand(balanceCmd)
}
