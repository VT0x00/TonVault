package cmd

import (
	"fmt"

	"github.com/VT0x00/tonvault/internal/wallet"
	"github.com/spf13/cobra"
)

var walletImportCmd = &cobra.Command{
	Use:   "import",
	Short: "Import wallet from seed phrase",
	Long:  `Restore an existing wallet by entering its 24-word seed phrase.`,
	RunE: func(cmd *cobra.Command, args []string) error {
		store, err := wallet.NewStore()
		if err != nil {
			return fmt.Errorf("failed to open wallet store: %w", err)
		}

		w, err := wallet.ImportWalletFromSeed(store, cfg.GetNetwork())
		if err != nil {
			return fmt.Errorf("failed to import wallet: %w", err)
		}

		fmt.Println()
		fmt.Printf("✓ Wallet imported successfully!\n")
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
	walletCmd.AddCommand(walletImportCmd)
}
