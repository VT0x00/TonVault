package ton

import (
	"context"
	"encoding/hex"
	"fmt"
	"math/big"
	"time"

	"github.com/xssnick/tonutils-go/address"
	"github.com/xssnick/tonutils-go/tlb"
	"github.com/xssnick/tonutils-go/ton/jetton"
	"github.com/xssnick/tonutils-go/ton/wallet"
	"github.com/xssnick/tonutils-go/tvm/cell"
)

type DexType string

const (
	DexStonFi DexType = "stonfi"
)

const (
	stonfiOpSwap = 0x25938561
)

var dexRouters = map[DexType]map[string]string{
	DexStonFi: {
		"mainnet": "EQB3ncyBUTjZUA5EnFKR5kEnBvb3cP0oOs3pPyFCzCjr5GQa",
		"testnet": "kQAHmpxT19l0Hso1Mw3pn_eej8tKQh7qH9kFBzBPWZ4A4LfQ",
	},
}

func (c *Client) GetDexRouterAddr(dex DexType, network string) (*address.Address, error) {
	routers, ok := dexRouters[dex]
	if !ok {
		return nil, fmt.Errorf("unknown DEX: %s", dex)
	}
	addrStr, ok := routers[network]
	if !ok {
		addrStr = routers["mainnet"]
	}
	return address.ParseAddr(addrStr)
}

func buildStonfiSwapPayload(minOut *big.Int, userAddr *address.Address) *cell.Cell {
	return cell.BeginCell().
		MustStoreUInt(stonfiOpSwap, 32).
		MustStoreUInt(0, 64).
		MustStoreAddr(userAddr).
		MustStoreBigCoins(minOut).
		MustStoreAddr(nil).
		MustStoreAddr(userAddr).
		MustStoreRef(cell.BeginCell().EndCell()).
		EndCell()
}

func (c *Client) SwapJetton(
	ctx context.Context,
	w *wallet.Wallet,
	routerAddr *address.Address,
	jettonMasterIn *address.Address,
	amountIn *big.Int,
	minAmountOut *big.Int,
) (*TransferResult, error) {
	ctx, cancel := context.WithTimeout(ctx, 120*time.Second)
	defer cancel()

	userAddr := w.WalletAddress()

	master := jetton.NewJettonMasterClient(c.api, jettonMasterIn)
	tokenWallet, err := master.GetJettonWallet(ctx, userAddr)
	if err != nil {
		return nil, fmt.Errorf("failed to get jetton wallet: %w", err)
	}

	forwardPayload := buildStonfiSwapPayload(minAmountOut, userAddr)

	transferPayload, err := jetton.BuildTransferPayload(
		routerAddr,
		userAddr,
		tlb.FromNanoTON(amountIn),
		tlb.MustFromTON("0.1"),
		forwardPayload,
		nil,
	)
	if err != nil {
		return nil, fmt.Errorf("failed to build transfer payload: %w", err)
	}

	msg := wallet.SimpleMessage(tokenWallet.Address(), tlb.MustFromTON("0.05"), transferPayload)
	tx, _, err := w.SendWaitTransaction(ctx, msg)
	if err != nil {
		return nil, fmt.Errorf("failed to execute swap: %w", err)
	}

	return &TransferResult{TxHash: hex.EncodeToString(tx.Hash)}, nil
}

func (c *Client) SwapTonToJetton(
	ctx context.Context,
	w *wallet.Wallet,
	routerAddr *address.Address,
	jettonMasterOut *address.Address,
	amountIn *big.Int,
	minAmountOut *big.Int,
) (*TransferResult, error) {
	ctx, cancel := context.WithTimeout(ctx, 120*time.Second)
	defer cancel()

	userAddr := w.WalletAddress()

	body := buildStonfiSwapPayload(minAmountOut, userAddr)

	msg := wallet.SimpleMessage(routerAddr, tlb.FromNanoTON(amountIn), body)
	tx, _, err := w.SendWaitTransaction(ctx, msg)
	if err != nil {
		return nil, fmt.Errorf("failed to execute swap: %w", err)
	}

	return &TransferResult{TxHash: hex.EncodeToString(tx.Hash)}, nil
}

func (c *Client) Swap(
	ctx context.Context,
	w *wallet.Wallet,
	dex DexType,
	network string,
	tokenIn *address.Address,
	tokenOut *address.Address,
	amountIn *big.Int,
	minAmountOut *big.Int,
) (*TransferResult, error) {
	routerAddr, err := c.GetDexRouterAddr(dex, network)
	if err != nil {
		return nil, fmt.Errorf("failed to get DEX router: %w", err)
	}

	if tokenIn == nil {
		return c.SwapTonToJetton(ctx, w, routerAddr, tokenOut, amountIn, minAmountOut)
	}
	return c.SwapJetton(ctx, w, routerAddr, tokenIn, amountIn, minAmountOut)
}
