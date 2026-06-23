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
	jettonSendAmount  string
	jettonSendComment string
)

var sendJettonCmd = &cobra.Command{
	Use:   "jetton <jetton_master_address> <to_address>",
	Short: "Send Jetton tokens",
	Long:  `Transfer Jetton tokens to any TON address.`,
	Args:  cobra.ExactArgs(2),
	RunE: func(cmd *cobra.Command, args []string) error {
		store, err := wallet.NewStore()
		if err != nil {
			return fmt.Errorf("failed to open wallet store: %w", err)
		}

		w := store.GetDefault()
		if w == nil {
			return fmt.Errorf("no wallets found. Create or import a wallet first")
		}

		jettonMasterAddr, err := address.ParseAddr(args[0])
		if err != nil {
			return fmt.Errorf("invalid jetton master address: %w", err)
		}

		toAddr, err := address.ParseAddr(args[1])
		if err != nil {
			return fmt.Errorf("invalid recipient address: %w", err)
		}

		amount := new(big.Int)
		amount, ok := amount.SetString(jettonSendAmount, 10)
		if !ok {
			return fmt.Errorf("invalid amount: %s", jettonSendAmount)
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

		fmt.Printf("Sending %s Jetton tokens (%s) to %s\n", jettonSendAmount, args[0], args[1])
		if jettonSendComment != "" {
			fmt.Printf("Comment: %s\n", jettonSendComment)
		}
		fmt.Print("Confirm? (y/N): ")

		var confirm string
		fmt.Scanln(&confirm)
		if confirm != "y" && confirm != "Y" {
			fmt.Println("Cancelled.")
			return nil
		}

		result, err := client.SendJetton(ctx, senderWallet, jettonMasterAddr, toAddr, amount, jettonSendComment)
		if err != nil {
			return fmt.Errorf("failed to send Jetton: %w", err)
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
	sendCmd.AddCommand(sendJettonCmd)
	sendJettonCmd.Flags().StringVarP(&jettonSendAmount, "amount", "a", "", "Amount in raw Jetton units (smallest unit)")
	sendJettonCmd.Flags().StringVarP(&jettonSendComment, "comment", "c", "", "Optional comment/memo")
	sendJettonCmd.MarkFlagRequired("amount")
}
