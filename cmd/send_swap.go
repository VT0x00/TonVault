package cmd

import (
	"fmt"
	"math/big"
	"strings"

	"github.com/VT0x00/tonvault/internal/ton"
	"github.com/VT0x00/tonvault/internal/wallet"
	"github.com/spf13/cobra"
	"github.com/xssnick/tonutils-go/address"
	tonwallet "github.com/xssnick/tonutils-go/ton/wallet"
)

var (
	swapFromToken string
	swapToToken   string
	swapAmount    string
	swapMinAmount string
	swapDex       string
)

const tonSymbol = "TON"

var sendSwapCmd = &cobra.Command{
	Use:   "swap",
	Short: "Swap tokens via DEX",
	Long: `Swap tokens via a decentralized exchange (STON.fi).

Supported tokens: "TON" for native coins, or a Jetton master address.

Examples:
  tonvault send swap --from TON --to EQD... --amount 1000000000 --min-amount 500000
  tonvault send swap --from EQD... --to TON --amount 1000000 --min-amount 500000000
  tonvault send swap --from EQD... --to EQD... --amount 1000000 --min-amount 500000`,
	Args: cobra.NoArgs,
	RunE: func(cmd *cobra.Command, args []string) error {
		store, err := wallet.NewStore()
		if err != nil {
			return fmt.Errorf("failed to open wallet store: %w", err)
		}

		w := store.GetDefault()
		if w == nil {
			return fmt.Errorf("no wallets found. Create or import a wallet first")
		}

		var tokenIn *address.Address
		var tokenOut *address.Address
		tokenInStr := strings.TrimSpace(swapFromToken)
		tokenOutStr := strings.TrimSpace(swapToToken)

		if tokenInStr == "" || tokenOutStr == "" {
			return fmt.Errorf("both --from and --to are required")
		}

		if tokenInStr != tonSymbol {
			tokenIn, err = address.ParseAddr(tokenInStr)
			if err != nil {
				return fmt.Errorf("invalid --from address: %w", err)
			}
		}

		if tokenOutStr != tonSymbol {
			tokenOut, err = address.ParseAddr(tokenOutStr)
			if err != nil {
				return fmt.Errorf("invalid --to address: %w", err)
			}
		}

		amountIn := new(big.Int)
		amountIn, ok := amountIn.SetString(swapAmount, 10)
		if !ok {
			return fmt.Errorf("invalid --amount: %s", swapAmount)
		}

		minAmountOut := new(big.Int)
		minAmountOut, ok = minAmountOut.SetString(swapMinAmount, 10)
		if !ok {
			return fmt.Errorf("invalid --min-amount: %s", swapMinAmount)
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

		dex := ton.DexType(swapDex)

		fmt.Printf("Swapping %s %s → %s via %s\n", swapAmount, tokenInStr, tokenOutStr, strings.ToUpper(string(dex)))
		fmt.Printf("Min amount out: %s\n", swapMinAmount)
		fmt.Print("Confirm? (y/N): ")

		var confirm string
		fmt.Scanln(&confirm)
		if confirm != "y" && confirm != "Y" {
			fmt.Println("Cancelled.")
			return nil
		}

		result, err := client.Swap(ctx, senderWallet, dex, w.Network, tokenIn, tokenOut, amountIn, minAmountOut)
		if err != nil {
			return fmt.Errorf("swap failed: %w", err)
		}

		fmt.Printf("✓ Swap transaction sent!\n")
		if result.TxHash != "" {
			fmt.Printf("Hash: %s\n", result.TxHash)
			fmt.Printf("Explorer: %s/transaction/%s\n", cfg.GetExplorerURL(), result.TxHash)
		}

		return nil
	},
}

func init() {
	sendCmd.AddCommand(sendSwapCmd)
	sendSwapCmd.Flags().StringVar(&swapFromToken, "from", "", "Source token address or \"TON\"")
	sendSwapCmd.Flags().StringVar(&swapToToken, "to", "", "Destination token address or \"TON\"")
	sendSwapCmd.Flags().StringVarP(&swapAmount, "amount", "a", "", "Amount to swap (raw units for jettons, nanoTON for TON)")
	sendSwapCmd.Flags().StringVar(&swapMinAmount, "min-amount", "", "Minimum amount to receive (raw units)")
	sendSwapCmd.Flags().StringVar(&swapDex, "dex", string(ton.DexStonFi), "DEX to use (stonfi)")
	sendSwapCmd.MarkFlagRequired("from")
	sendSwapCmd.MarkFlagRequired("to")
	sendSwapCmd.MarkFlagRequired("amount")
	sendSwapCmd.MarkFlagRequired("min-amount")
}
