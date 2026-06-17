package cmd

import (
	"fmt"

	"github.com/VT0x00/tonvault/internal/models"
	"github.com/VT0x00/tonvault/internal/ton"
	"github.com/VT0x00/tonvault/internal/wallet"
	"github.com/spf13/cobra"
	"github.com/xssnick/tonutils-go/address"
)

var balanceShowCmd = &cobra.Command{
	Use:   "show [wallet_id]",
	Short: "Show TON balance",
	Long:  `Display TON balance for a wallet.`,
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

		balance, err := client.GetBalance(ctx, addr)
		if err != nil {
			return fmt.Errorf("failed to get balance: %w", err)
		}

		fmt.Printf("Wallet: %s\n", w.Name)
		fmt.Printf("Address: %s\n", w.Address)
		fmt.Printf("Balance: %s TON\n", ton.FormatBalance(balance.TON))
		fmt.Printf("         %s nanoTON\n", balance.TON.String())

		return nil
	},
}

func init() {
	balanceCmd.AddCommand(balanceShowCmd)
}
