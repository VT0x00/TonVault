package cmd

import (
	"fmt"

	"github.com/VT0x00/tonvault/internal/wallet"
	"github.com/spf13/cobra"
)

var walletCreateFlags struct {
	network string
}

var walletCreateCmd = &cobra.Command{
	Use:   "create",
	Short: "Create a new wallet",
	Long:  `Generate a new wallet with a 24-word BIP39 mnemonic seed phrase.`,
	RunE: func(cmd *cobra.Command, args []string) error {
		store, err := wallet.NewStore()
		if err != nil {
			return fmt.Errorf("failed to open wallet store: %w", err)
		}

		network := walletCreateFlags.network
		if network == "" {
			network = cfg.GetNetwork()
		}

		w, err := wallet.CreateNewWallet(store, network)
		if err != nil {
			return fmt.Errorf("failed to create wallet: %w", err)
		}

		fmt.Println()
		fmt.Printf("✓ Wallet created successfully!\n")
		fmt.Printf("  Name:    %s\n", w.Name)
		fmt.Printf("  Address: %s\n", w.Address)
		fmt.Printf("  Version: %s\n", w.Version)
		fmt.Printf("  Network: %s\n", w.Network)

		if store.Count() == 1 {
			store.SetDefault(w.ID)
			fmt.Println("  (set as default wallet)")
		}

		return nil
	},
}

func init() {
	walletCreateCmd.Flags().StringVar(&walletCreateFlags.network, "network", "", "network (mainnet/testnet, defaults to global config)")
	walletCmd.AddCommand(walletCreateCmd)
}
