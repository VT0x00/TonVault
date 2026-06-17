package ton

import (
	"context"
	"fmt"
	"math/big"
	"sync"
	"time"

	"github.com/xssnick/tonutils-go/address"
	"github.com/xssnick/tonutils-go/ton/jetton"
)

type KnownJetton struct {
	MasterAddress string
	Symbol        string
	Decimals      int
}

type JettonBalance struct {
	MasterAddress string
	Symbol        string
	WalletAddress string
	Balance       *big.Int
	Decimals      int
}

func GetKnownJettons(network string) []KnownJetton {
	if network == "testnet" {
		return testnetJettons
	}
	return mainnetJettons
}

var mainnetJettons = []KnownJetton{
	{
		MasterAddress: "EQCxE8D7E6kFzYBz1T5K2YXKsf1hCvjPBVj8e0Jw9KQkR0Vd",
		Symbol:        "USDT",
		Decimals:      6,
	},
	{
		MasterAddress: "EQAvlWFDxGF2lXm67y4yzC17wYKD9A5K3HxVhCxqLmsptcNn",
		Symbol:        "NOT",
		Decimals:      9,
	},
	{
		MasterAddress: "EQCvx3c6oGjGz6v_t4zpGJXXmRCB2sYx5yhN3Y-xQ1F7nKvc",
		Symbol:        "DOGS",
		Decimals:      9,
	},
	{
		MasterAddress: "EQBwKLLBK9sM3xhk6tC94BPKkM7AY7GDRQijYfIhYjEZM-LA",
		Symbol:        "STON",
		Decimals:      9,
	},
	{
		MasterAddress: "EQBlWgKnh_qbFYTXfKgGAQPxkxFsArDOSr9nlARSzydpNPwA",
		Symbol:        "GOVNO",
		Decimals:      9,
	},
	// {
	// 	MasterAddress: "",
	// 	Symbol:        "",
	// 	Decimals:      9,
	// },
}

var testnetJettons = []KnownJetton{
	{
		MasterAddress: "kQDRUsCqkRgJZFLnjyQ3P9jhphTGQfzhLE_9e6bYhwKElEwG",
		Symbol:        "USDT",
		Decimals:      6,
	},
}

func (c *Client) GetJettonBalances(ctx context.Context, ownerAddr *address.Address, jettons []KnownJetton) ([]JettonBalance, error) {
	ctx, cancel := context.WithTimeout(ctx, 30*time.Second)
	defer cancel()

	type result struct {
		jb  JettonBalance
		err error
	}

	ch := make(chan result, len(jettons))
	var wg sync.WaitGroup

	for _, j := range jettons {
		wg.Add(1)
		go func(j KnownJetton) {
			defer wg.Done()

			masterAddr, err := address.ParseAddr(j.MasterAddress)
			if err != nil {
				ch <- result{err: fmt.Errorf("invalid master address for %s: %w", j.Symbol, err)}
				return
			}

			master := jetton.NewJettonMasterClient(c.api, masterAddr)
			tokenWallet, err := master.GetJettonWallet(ctx, ownerAddr)
			if err != nil {
				ch <- result{err: fmt.Errorf("failed to get jetton wallet for %s: %w", j.Symbol, err)}
				return
			}

			balance, err := tokenWallet.GetBalance(ctx)
			if err != nil {
				ch <- result{err: fmt.Errorf("failed to get balance for %s: %w", j.Symbol, err)}
				return
			}

			ch <- result{jb: JettonBalance{
				MasterAddress: j.MasterAddress,
				Symbol:        j.Symbol,
				WalletAddress: tokenWallet.Address().String(),
				Balance:       balance,
				Decimals:      j.Decimals,
			}}
		}(j)
	}

	wg.Wait()
	close(ch)

	var balances []JettonBalance
	for r := range ch {
		if r.err != nil {
			continue
		}
		if r.jb.Balance != nil && r.jb.Balance.Sign() > 0 {
			balances = append(balances, r.jb)
		}
	}

	return balances, nil
}

func FormatJettonBalance(balance *big.Int, decimals int) string {
	if decimals == 0 {
		return balance.String()
	}
	div := new(big.Float).SetFloat64(pow10(decimals))
	val := new(big.Float).Quo(
		new(big.Float).SetInt(balance),
		div,
	)
	return val.Text('f', int(decimals))
}

func pow10(n int) float64 {
	r := 1.0
	for i := 0; i < n; i++ {
		r *= 10
	}
	return r
}
