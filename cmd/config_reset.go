package cmd

import (
	"fmt"

	"github.com/spf13/cobra"
)

var configResetCmd = &cobra.Command{
	Use:   "reset",
	Short: "Reset config to defaults",
	Long:  `Reset all configuration values to their defaults.`,
	RunE: func(cmd *cobra.Command, args []string) error {
		fmt.Print("Reset configuration to defaults? (y/N): ")
		var confirm string
		fmt.Scanln(&confirm)
		if confirm != "y" && confirm != "Y" {
			fmt.Println("Cancelled.")
			return nil
		}
		if err := cfg.Reset(); err != nil {
			return fmt.Errorf("failed to reset config: %w", err)
		}
		fmt.Println("✓ Configuration reset to defaults.")
		return nil
	},
}

func init() {
	configCmd.AddCommand(configResetCmd)
}
