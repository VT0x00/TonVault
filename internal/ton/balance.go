package ton

import (
	"context"
	"fmt"
	"math/big"
	"time"

	"github.com/xssnick/tonutils-go/address"
	"github.com/xssnick/tonutils-go/tlb"
)

type BalanceInfo struct {
	TON *big.Int
}

func (c *Client) GetBalance(ctx context.Context, addr *address.Address) (*BalanceInfo, error) {
	ctx, cancel := context.WithTimeout(ctx, 10*time.Second)
	defer cancel()

	block, err := c.api.CurrentMasterchainInfo(ctx)
	if err != nil {
		return nil, fmt.Errorf("failed to get masterchain info: %w", err)
	}

	account, err := c.api.GetAccount(ctx, block, addr)
	if err != nil {
		return nil, fmt.Errorf("failed to get account: %w", err)
	}

	balance := big.NewInt(0)
	if account.State != nil {
		balance = account.State.Balance.Nano()
	}

	return &BalanceInfo{
		TON: balance,
	}, nil
}

func FormatBalance(nano *big.Int) string {
	ton := new(big.Float).Quo(
		new(big.Float).SetInt(nano),
		new(big.Float).SetFloat64(1e9),
	)
	return fmt.Sprintf("%.9f", ton)
}

func FormatBalanceShort(nano *big.Int) string {
	ton := new(big.Float).Quo(
		new(big.Float).SetInt(nano),
		new(big.Float).SetFloat64(1e9),
	)
	return fmt.Sprintf("%.4f", ton)
}

func NanoToTON(nano *big.Int) tlb.Coins {
	return tlb.FromNanoTON(nano)
}
