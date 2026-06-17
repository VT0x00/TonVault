package cmd

import (
	"fmt"
	"os"

	"github.com/spf13/cobra"
)

var networkStatusCmd = &cobra.Command{
	Use:   "status",
	Short: "Show network status",
	Long:  `Display current network configuration and connection status.`,
	RunE: func(cmd *cobra.Command, args []string) error {
		fmt.Printf("Network:       %s\n", cfg.GetNetwork())
		fmt.Printf("Lite Servers:  %v\n", cfg.GetLiteServers())
		fmt.Printf("Explorer URL:  %s\n", cfg.GetExplorerURL())
		fmt.Printf("Config Path:   %s\n", cfg.Path())
		hostname, _ := os.Hostname()
		fmt.Printf("Host:          %s\n", hostname)
		return nil
	},
}

func init() {
	networkCmd.AddCommand(networkStatusCmd)
}
