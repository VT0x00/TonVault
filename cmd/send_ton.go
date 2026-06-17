package cmd

import (
	"fmt"
	"math/big"

	"github.com/VT0x00/tonvault/internal/ton"
	"github.com/VT0x00/tonvault/internal/wallet"
	"github.com/spf13/cobra"
	"github.com/xssnick/tonutils-go/address"
	tonwallet "github.com/xssnick/tonutils-go/ton/wallet"
)

var (
	sendAmount   string
	sendComment  string
	sendNoBounce bool
)

var sendTonCmd = &cobra.Command{
	Use:   "ton <to_address>",
	Short: "Send TON coins",
	Long:  `Transfer TON coins to any TON address.`,
	Args:  cobra.ExactArgs(1),
	RunE: func(cmd *cobra.Command, args []string) error {
		store, err := wallet.NewStore()
		if err != nil {
			return fmt.Errorf("failed to open wallet store: %w", err)
		}

		w := store.GetDefault()
		if w == nil {
			return fmt.Errorf("no wallets found. Create or import a wallet first")
		}

		toAddr, err := address.ParseAddr(args[0])
		if err != nil {
			return fmt.Errorf("invalid recipient address: %w", err)
		}

		amountNano := new(big.Int)
		amountNano, ok := amountNano.SetString(sendAmount, 10)
		if !ok {
			return fmt.Errorf("invalid amount: %s", sendAmount)
		}

		ctx := cmd.Context()
		client, err := ton.NewClientForNetwork(ctx, w.Network)
		if err != nil {
			return fmt.Errorf("failed to connect to TON network: %w", err)
		}
		defer client.Close()

		words, err := wallet.RecoverSeedPhrase(store, w.ID)
		if err != nil {
			return err
		}

		senderWallet, err := tonwallet.FromSeed(client.GetAPI(), words, tonwallet.V4R2)
		if err != nil {
			return fmt.Errorf("failed to initialize wallet: %w", err)
		}

		fmt.Printf("Sending %s nanoTON to %s\n", sendAmount, args[0])
		if sendComment != "" {
			fmt.Printf("Comment: %s\n", sendComment)
		}
		fmt.Print("Confirm? (y/N): ")

		var confirm string
		fmt.Scanln(&confirm)
		if confirm != "y" && confirm != "Y" {
			fmt.Println("Cancelled.")
			return nil
		}

		result, err := client.SendTON(ctx, senderWallet, toAddr, amountNano, sendComment)
		if err != nil {
			return fmt.Errorf("failed to send: %w", err)
		}

		fmt.Printf("✓ Transaction sent!\n")
		if result.TxHash != "" {
			fmt.Printf("Hash: %s\n", result.TxHash)
			fmt.Printf("Explorer: %s/transaction/%s\n", cfg.GetExplorerURL(), result.TxHash)
		}

		return nil
	},
}

func init() {
	sendCmd.AddCommand(sendTonCmd)
	sendTonCmd.Flags().StringVarP(&sendAmount, "amount", "a", "", "Amount in nanoTON")
	sendTonCmd.Flags().StringVarP(&sendComment, "comment", "c", "", "Optional comment/memo")
	sendTonCmd.Flags().BoolVarP(&sendNoBounce, "no-bounce", "", false, "Disable bounce")
	sendTonCmd.MarkFlagRequired("amount")
}
