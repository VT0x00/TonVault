package cmd

import (
	"fmt"

	"github.com/VT0x00/tonvault/internal/wallet"
	"github.com/spf13/cobra"
)

var walletExportCmd = &cobra.Command{
	Use:   "export <wallet_id>",
	Short: "Export wallet seed phrase",
	Long:  `Decrypt and display the wallet's seed phrase. Requires password.`,
	Args:  cobra.ExactArgs(1),
	RunE: func(cmd *cobra.Command, args []string) error {
		store, err := wallet.NewStore()
		if err != nil {
			return fmt.Errorf("failed to open wallet store: %w", err)
		}

		fmt.Printf("Exporting seed phrase for wallet '%s'\n", args[0])
		fmt.Println("⚠ Make sure no one is looking at your screen.")
		fmt.Println()

		words, err := wallet.RecoverSeedPhrase(store, args[0])
		if err != nil {
			return err
		}

		fmt.Println("Seed phrase:")
		fmt.Println("============")
		for i, word := range words {
			fmt.Printf("%3d. %s\n", i+1, word)
		}
		fmt.Println("============")

		return nil
	},
}

func init() {
	walletCmd.AddCommand(walletExportCmd)
}
