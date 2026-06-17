package cmd

import (
	"fmt"
	"os"

	"github.com/VT0x00/tonvault/internal/wallet"
	"github.com/olekukonko/tablewriter"
	"github.com/olekukonko/tablewriter/tw"
	"github.com/spf13/cobra"
)

var walletListCmd = &cobra.Command{
	Use:   "list",
	Short: "List all wallets",
	Long:  `Display all locally stored wallets with their addresses and aliases.`,
	RunE: func(cmd *cobra.Command, args []string) error {
		store, err := wallet.NewStore()
		if err != nil {
			return fmt.Errorf("failed to open wallet store: %w", err)
		}

		wallets := store.List()
		if len(wallets) == 0 {
			fmt.Println("No wallets found. Use 'tonvault wallet create' to create one.")
			return nil
		}

		table := tablewriter.NewTable(os.Stdout,
			tablewriter.WithRendition(tw.Rendition{
				Settings: tw.Settings{
					Lines: tw.Lines{
						ShowHeaderLine: tw.On,
					},
				},
			}),
		)

		table.Header("ID", "Name", "Address", "Version", "Network", "Default")

		for _, w := range wallets {
			defaultMark := ""
			if w.IsDefault {
				defaultMark = "*"
			}
			addr := w.Address
			if len(addr) > 20 {
				addr = addr[:10] + "..." + addr[len(addr)-10:]
			}
			if err := table.Append(
				w.ID,
				w.Name,
				addr,
				w.Version,
				w.Network,
				defaultMark,
			); err != nil {
				return err
			}
		}

		return table.Render()
	},
}

func init() {
	walletCmd.AddCommand(walletListCmd)
}
