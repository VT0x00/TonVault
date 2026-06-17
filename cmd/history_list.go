package cmd

import (
	"fmt"
	"os"
	"time"

	"github.com/VT0x00/tonvault/internal/models"
	"github.com/VT0x00/tonvault/internal/ton"
	"github.com/VT0x00/tonvault/internal/wallet"
	"github.com/olekukonko/tablewriter"
	"github.com/olekukonko/tablewriter/tw"
	"github.com/spf13/cobra"
	"github.com/xssnick/tonutils-go/address"
)

var historyListFlags struct {
	limit   uint32
	address string
	network string
}

var historyListCmd = &cobra.Command{
	Use:   "list [wallet_id]",
	Short: "List recent transactions",
	Long: `Display recent transaction history.

Accepts a wallet ID (from local store), a raw TON address via --address,
or defaults to the configured default wallet.`,
	Args: cobra.MaximumNArgs(1),
	RunE: func(cmd *cobra.Command, args []string) error {
		var addr *address.Address
		network := historyListFlags.network

		if historyListFlags.address != "" {
			var err error
			addr, err = address.ParseAddr(historyListFlags.address)
			if err != nil {
				return fmt.Errorf("invalid address: %w", err)
			}
			if network == "" {
				network = "mainnet"
			}
		} else {
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

			addr, err = address.ParseAddr(w.Address)
			if err != nil {
				return fmt.Errorf("invalid wallet address: %w", err)
			}
			network = w.Network
		}

		ctx := cmd.Context()
		client, err := ton.NewClientForNetwork(ctx, network)
		if err != nil {
			return fmt.Errorf("failed to connect to TON network: %w", err)
		}
		defer client.Close()

		if key := cfg.GetTonCenterAPIKey(); key != "" {
			client.SetTonCenterAPIKey(key, network)
		}

		records, err := client.GetTransactionHistory(ctx, addr, historyListFlags.limit)
		if err != nil {
			return fmt.Errorf("failed to get transaction history: %w", err)
		}

		if len(records) == 0 {
			fmt.Println("No transactions found.")
			return nil
		}

		if uint32(len(records)) < historyListFlags.limit {
			fmt.Fprintf(os.Stderr, "Warning: fetched %d of %d requested transactions (lite server data retention limit)\n",
				len(records), historyListFlags.limit)
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

		table.Header("Date", "Type", "Amount (TON)", "From / To", "Fee", "Comment")

		for _, r := range records {
			ts := r.Timestamp.Format(time.RFC3339)
			amount := r.Amount
			if amount == "" {
				amount = "0"
			}

			addrStr := r.From
			if r.Type == "outgoing" {
				addrStr = r.To
			}
			if len(addrStr) > 16 {
				addrStr = addrStr[:8] + "..." + addrStr[len(addrStr)-8:]
			}

			fee := ton.FormatNanoTON(r.Fee)

			comment := r.Comment
			if len(comment) > 30 {
				comment = comment[:27] + "..."
			}

			if err := table.Append(ts, r.Type, amount, addrStr, fee, comment); err != nil {
				return err
			}
		}

		return table.Render()
	},
}

func init() {
	historyListCmd.Flags().Uint32Var(&historyListFlags.limit, "limit", 10, "number of transactions to show")
	historyListCmd.Flags().StringVar(&historyListFlags.address, "address", "", "TON address to query (bypasses wallet store)")
	historyListCmd.Flags().StringVar(&historyListFlags.network, "network", "", "network (mainnet/testnet, defaults to wallet network or mainnet)")
	historyCmd.AddCommand(historyListCmd)
}
