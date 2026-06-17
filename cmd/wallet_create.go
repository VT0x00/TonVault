package cmd

import (
	"fmt"

	"github.com/VT0x00/tonvault/internal/wallet"
	"github.com/spf13/cobra"
)

var walletCreateCmd = &cobra.Command{
	Use:   "create",
	Short: "Create a new wallet",
	Long:  `Generate a new wallet with a 24-word BIP39 mnemonic seed phrase.`,
	RunE: func(cmd *cobra.Command, args []string) error {
		store, err := wallet.NewStore()
		if err != nil {
			return fmt.Errorf("failed to open wallet store: %w", err)
		}

		w, err := wallet.CreateNewWallet(store, cfg.GetNetwork())
		if err != nil {
			return fmt.Errorf("failed to create wallet: %w", err)
		}

		fmt.Println()
		fmt.Printf("✓ Wallet created successfully!\n")
		fmt.Printf("  Name:    %s\n", w.Name)
		fmt.Printf("  Address: %s\n", w.Address)
		fmt.Printf("  Version: %s\n", w.Version)

		if store.Count() == 1 {
			store.SetDefault(w.ID)
			fmt.Println("  (set as default wallet)")
		}

		return nil
	},
}

func init() {
	walletCmd.AddCommand(walletCreateCmd)
}
