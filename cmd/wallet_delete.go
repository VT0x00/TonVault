package cmd

import (
	"fmt"
	"os"

	"github.com/VT0x00/tonvault/internal/wallet"
	"github.com/olekukonko/tablewriter"
	"github.com/olekukonko/tablewriter/tw"
	"github.com/spf13/cobra"
)

var walletDeleteCmd = &cobra.Command{
	Use:   "delete <wallet_id>",
	Short: "Delete a wallet",
	Long:  `Remove a wallet from local storage.`,
	Args:  cobra.ExactArgs(1),
	RunE: func(cmd *cobra.Command, args []string) error {
		store, err := wallet.NewStore()
		if err != nil {
			return fmt.Errorf("failed to open wallet store: %w", err)
		}

		fmt.Printf("Are you sure you want to delete wallet '%s'? (y/N): ", args[0])
		var confirm string
		fmt.Scanln(&confirm)
		if confirm != "y" && confirm != "Y" {
			fmt.Println("Deletion cancelled.")
			return nil
		}

		if err := store.Delete(args[0]); err != nil {
			fmt.Fprintf(os.Stderr, "Error: wallet '%s' not found.\n\n", args[0])
			fmt.Fprintf(os.Stderr, "Available wallets:\n")

			wallets := store.List()
			if len(wallets) == 0 {
				fmt.Fprintln(os.Stderr, "  (no wallets)")
			} else {
				t := tablewriter.NewTable(os.Stderr,
					tablewriter.WithRendition(tw.Rendition{
						Settings: tw.Settings{
							Lines: tw.Lines{ShowHeaderLine: tw.On},
						},
					}),
				)
				t.Header("ID", "Name")
				for _, w := range wallets {
					t.Append(w.ID, w.Name)
				}
				t.Render()
			}
			return fmt.Errorf("deletion failed")
		}

		fmt.Printf("✓ Wallet '%s' deleted.\n", args[0])
		return nil
	},
}

func init() {
	walletCmd.AddCommand(walletDeleteCmd)
}
