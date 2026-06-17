package cmd

import (
	"fmt"
	"time"

	"github.com/VT0x00/tonvault/internal/ton"
	"github.com/VT0x00/tonvault/internal/wallet"
	"github.com/spf13/cobra"
	"github.com/xssnick/tonutils-go/address"
)

var historyShowFlags struct {
	address string
	network string
}

var historyShowCmd = &cobra.Command{
	Use:   "show <tx_hash>",
	Short: "Show transaction details",
	Long: `Display detailed information about a specific transaction.

Accepts a wallet ID (from local store), a raw TON address via --address,
or defaults to the configured default wallet.`,
	Args: cobra.ExactArgs(1),
	RunE: func(cmd *cobra.Command, args []string) error {
		txHash := args[0]

		var addr *address.Address
		network := historyShowFlags.network

		if historyShowFlags.address != "" {
			var err error
			addr, err = address.ParseAddr(historyShowFlags.address)
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

			w := store.GetDefault()
			if w == nil {
				return fmt.Errorf("no wallets found")
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

		records, err := client.GetTransactionHistory(ctx, addr, 100)
		if err != nil {
			return fmt.Errorf("failed to get transaction history: %w", err)
		}

		var tx *ton.TxRecord
		for _, r := range records {
			if r.Hash == txHash {
				tx = r
				break
			}
		}

		if tx == nil {
			return fmt.Errorf("transaction %s not found in recent history", txHash)
		}

		ts := tx.Timestamp.Format(time.RFC3339)

		fmt.Printf("Hash:      %s\n", tx.Hash)
		fmt.Printf("Type:      %s\n", tx.Type)
		fmt.Printf("Amount:    %s TON\n", tx.Amount)
		fmt.Printf("From:      %s\n", tx.From)
		fmt.Printf("To:        %s\n", tx.To)
		fmt.Printf("Fee:       %s\n", ton.FormatNanoTON(tx.Fee))
		fmt.Printf("Date:      %s\n", ts)
		if tx.Comment != "" {
			fmt.Printf("Comment:   %s\n", tx.Comment)
		}

		return nil
	},
}

func init() {
	historyShowCmd.Flags().StringVar(&historyShowFlags.address, "address", "", "TON address (bypasses wallet store)")
	historyShowCmd.Flags().StringVar(&historyShowFlags.network, "network", "", "network (mainnet/testnet, defaults to wallet network or mainnet)")
	historyCmd.AddCommand(historyShowCmd)
}
