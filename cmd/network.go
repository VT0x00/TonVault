package cmd

import "github.com/spf13/cobra"

var networkCmd = &cobra.Command{
	Use:   "network",
	Short: "Manage network settings",
	Long:  `Switch between Mainnet and Testnet, view network status.`,
}

func init() {
	rootCmd.AddCommand(networkCmd)
}
