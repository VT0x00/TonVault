package cmd

import "github.com/spf13/cobra"

var walletCmd = &cobra.Command{
	Use:   "wallet",
	Short: "Manage wallets",
	Long:  `Create, import, list, show info, delete, and export wallets.`,
}

func init() {
	rootCmd.AddCommand(walletCmd)
}
