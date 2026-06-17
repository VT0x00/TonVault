package cmd

import (
	"fmt"
	"os"

	"github.com/VT0x00/tonvault/internal/models"
	"github.com/VT0x00/tonvault/internal/ton"
	"github.com/VT0x00/tonvault/internal/wallet"
	"github.com/olekukonko/tablewriter"
	"github.com/olekukonko/tablewriter/tw"
	"github.com/xssnick/tonutils-go/address"
	"github.com/spf13/cobra"
)

var balanceJettonsCmd = &cobra.Command{
	Use:   "jettons [wallet_id]",
	Short: "Show Jetton balances",
	Long:  `List all Jetton tokens held by the wallet.`,
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
		}
		if err != nil || w == nil {
			return fmt.Errorf("wallet not found")
		}

		addr, err := address.ParseAddr(w.Address)
		if err != nil {
			return fmt.Errorf("invalid wallet address: %w", err)
		}

		ctx := cmd.Context()
		client, err := ton.NewClientForNetwork(ctx, w.Network)
		if err != nil {
			return fmt.Errorf("failed to connect to TON network: %w", err)
		}
		defer client.Close()

		balances, err := client.GetJettonBalances(ctx, addr, ton.GetKnownJettons(w.Network))
		if err != nil {
			return fmt.Errorf("failed to query jetton balances: %w", err)
		}

		if len(balances) == 0 {
			fmt.Println("No Jetton balances found.")
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

		table.Header("Symbol", "Balance", "Wallet Address")

		for _, b := range balances {
			formatted := ton.FormatJettonBalance(b.Balance, b.Decimals)
			walletAddr := b.WalletAddress
			if len(walletAddr) > 16 {
				walletAddr = walletAddr[:8] + "..." + walletAddr[len(walletAddr)-8:]
			}
			if err := table.Append(b.Symbol, formatted, walletAddr); err != nil {
				return err
			}
		}

		return table.Render()
	},
}

func init() {
	balanceCmd.AddCommand(balanceJettonsCmd)
}
