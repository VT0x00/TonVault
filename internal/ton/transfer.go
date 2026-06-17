package ton

import (
	"context"
	"fmt"
	"math/big"
	"time"

	"github.com/xssnick/tonutils-go/address"
	"github.com/xssnick/tonutils-go/tlb"
	"github.com/xssnick/tonutils-go/ton/jetton"
	"github.com/xssnick/tonutils-go/ton/wallet"
	"github.com/xssnick/tonutils-go/tvm/cell"
)

type TransferResult struct {
	TxHash string
}

func (c *Client) SendTON(
	ctx context.Context,
	w *wallet.Wallet,
	toAddr *address.Address,
	amountNano *big.Int,
	comment string,
) (*TransferResult, error) {
	ctx, cancel := context.WithTimeout(ctx, 60*time.Second)
	defer cancel()

	coins := tlb.FromNanoTON(amountNano)

	var err error
	if comment != "" {
		err = w.Transfer(ctx, toAddr, coins, comment)
	} else {
		err = w.Transfer(ctx, toAddr, coins, "")
	}

	if err != nil {
		return nil, fmt.Errorf("failed to send TON: %w", err)
	}

	return &TransferResult{}, nil
}

func (c *Client) SendJetton(
	ctx context.Context,
	w *wallet.Wallet,
	jettonMasterAddr *address.Address,
	toAddr *address.Address,
	amount *big.Int,
	comment string,
) (*TransferResult, error) {
	ctx, cancel := context.WithTimeout(ctx, 60*time.Second)
	defer cancel()

	token := jetton.NewJettonMasterClient(c.api, jettonMasterAddr)
	tokenWallet, err := token.GetJettonWallet(ctx, w.WalletAddress())
	if err != nil {
		return nil, fmt.Errorf("failed to get jetton wallet: %w", err)
	}

	amt := tlb.FromNanoTON(amount)

	var forwardPayload *cell.Cell
	if comment != "" {
		forwardPayload, err = wallet.CreateCommentCell(comment)
		if err != nil {
			return nil, fmt.Errorf("failed to create comment cell: %w", err)
		}
	}

	transferPayload, err := jetton.BuildTransferPayload(
		toAddr,
		toAddr,
		amt,
		tlb.ZeroCoins,
		forwardPayload,
		nil,
	)
	if err != nil {
		return nil, fmt.Errorf("failed to build transfer payload: %w", err)
	}

	msg := wallet.SimpleMessage(tokenWallet.Address(), tlb.MustFromTON("0.05"), transferPayload)
	_, _, err = w.SendWaitTransaction(ctx, msg)
	if err != nil {
		return nil, fmt.Errorf("failed to send Jetton: %w", err)
	}

	return &TransferResult{}, nil
}
