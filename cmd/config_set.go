package cmd

import (
	"fmt"

	"github.com/spf13/cobra"
)

var configSetCmd = &cobra.Command{
	Use:   "set <key> <value>",
	Short: "Set a config value",
	Long:  `Set a configuration value (e.g., "log_level debug", "explorer_url ...").`,
	Args:  cobra.ExactArgs(2),
	RunE: func(cmd *cobra.Command, args []string) error {
		if err := cfg.Set(args[0], args[1]); err != nil {
			return fmt.Errorf("failed to set config: %w", err)
		}
		fmt.Printf("✓ %s = %s\n", args[0], args[1])
		return nil
	},
}

func init() {
	configCmd.AddCommand(configSetCmd)
}
