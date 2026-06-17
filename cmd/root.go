package cmd

import (
	"fmt"
	"os"

	"github.com/VT0x00/tonvault/internal/config"
	"github.com/spf13/cobra"
)

var (
	cfg            *config.Manager
	cfgFile        string
)

var rootCmd = &cobra.Command{
	Use:   "tonvault",
	Short: "TonVault - A feature-rich CLI wallet for the TON blockchain",
	Long: `TonVault is a command-line wallet for The Open Network (TON) blockchain.
It supports wallet management, balance checks, TON and Jetton transfers,
and transaction history — all from your terminal.`,
	PersistentPreRunE: func(cmd *cobra.Command, args []string) error {
		var err error
		if cfg == nil {
			cfg, err = config.NewManager()
			if err != nil {
				return fmt.Errorf("failed to initialize config: %w", err)
			}
		}
		return nil
	},
}

func Execute() {
	if err := rootCmd.Execute(); err != nil {
		fmt.Fprintln(os.Stderr, err)
		os.Exit(1)
	}
}

func init() {
	rootCmd.PersistentFlags().StringVar(&cfgFile, "config", "", "config file path")
}

func GetConfig() *config.Manager {
	return cfg
}
