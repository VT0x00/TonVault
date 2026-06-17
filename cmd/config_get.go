package cmd

import (
	"fmt"

	"github.com/spf13/cobra"
)

var configGetCmd = &cobra.Command{
	Use:   "get [key]",
	Short: "Get config value(s)",
	Long:  `Get a specific config value or all values.`,
	Args:  cobra.MaximumNArgs(1),
	RunE: func(cmd *cobra.Command, args []string) error {
		if len(args) > 0 {
			fmt.Printf("%s = %s\n", args[0], cfg.Get(args[0]))
		} else {
			fmt.Printf("network           = %s\n", cfg.GetNetwork())
			fmt.Printf("default_wallet_id = %s\n", cfg.GetDefaultWalletID())
			fmt.Printf("explorer_url      = %s\n", cfg.GetExplorerURL())
			fmt.Printf("log_level         = %s\n", cfg.GetLogLevel())
		}
		return nil
	},
}

func init() {
	configCmd.AddCommand(configGetCmd)
}
