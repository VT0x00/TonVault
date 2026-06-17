package cmd

import (
	"fmt"

	"github.com/VT0x00/tonvault/internal/models"
	"github.com/VT0x00/tonvault/internal/wallet"
	"github.com/spf13/cobra"
)

var walletInfoCmd = &cobra.Command{
	Use:   "info [wallet_id]",
	Short: "Show wallet details",
	Long:  `Display detailed information about a specific wallet.`,
	Args:  cobra.MaximumNArgs(1),
	RunE: func(cmd *cobra.Command, args []string) error {
		store, err := wallet.NewStore()
		if err != nil {
			return fmt.Errorf("failed to open wallet store: %w", err)
		}

		var w *models.Wallet
		if len(args) > 0 {
			w, err = store.Get(args[0])
		} else {
			w = store.GetDefault()
			if w == nil {
				return fmt.Errorf("no wallets found. Use 'tonvault wallet create' to create one")
			}
			fmt.Printf("Showing default wallet\n\n")
		}

		if err != nil {
			return fmt.Errorf("wallet not found: %w", err)
		}

		fmt.Printf("Name:       %s\n", w.Name)
		fmt.Printf("ID:         %s\n", w.ID)
		fmt.Printf("Address:    %s\n", w.Address)
		fmt.Printf("Public Key: %s\n", w.PublicKey)
		fmt.Printf("Version:    %s\n", w.Version)
		fmt.Printf("Network:    %s\n", w.Network)
		fmt.Printf("Created:    %s\n", w.CreatedAt.Format("2006-01-02 15:04:05"))
		fmt.Printf("Default:    %v\n", w.IsDefault)

		return nil
	},
}

func init() {
	walletCmd.AddCommand(walletInfoCmd)
}
