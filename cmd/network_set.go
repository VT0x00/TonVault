package cmd

import (
	"fmt"

	"github.com/spf13/cobra"
)

var networkSetCmd = &cobra.Command{
	Use:   "set <mainnet|testnet>",
	Short: "Set network",
	Long:  `Switch between Mainnet and Testnet.`,
	Args:  cobra.ExactArgs(1),
	RunE: func(cmd *cobra.Command, args []string) error {
		net := args[0]
		if net != "mainnet" && net != "testnet" {
			return fmt.Errorf("invalid network: %s (use 'mainnet' or 'testnet')", net)
		}

		if err := cfg.SetNetwork(net); err != nil {
			return fmt.Errorf("failed to set network: %w", err)
		}

		fmt.Printf("✓ Network set to %s\n", net)
		return nil
	},
}

func init() {
	networkCmd.AddCommand(networkSetCmd)
}
