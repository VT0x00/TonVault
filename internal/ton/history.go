package ton

import (
	"context"
	"encoding/hex"
	"errors"
	"fmt"
	"time"

	"github.com/xssnick/tonutils-go/address"
	"github.com/xssnick/tonutils-go/tlb"
	"github.com/xssnick/tonutils-go/ton"
)

type TxRecord struct {
	Hash        string
	Type        string
	Amount      string
	From        string
	To          string
	Comment     string
	Fee         string
	Timestamp   time.Time
	BlockHeight uint64
}

func (c *Client) GetTransactionHistory(
	ctx context.Context,
	addr *address.Address,
	limit uint32,
) ([]*TxRecord, error) {
	ctx, cancel := context.WithTimeout(ctx, 30*time.Second)
	defer cancel()

	block, err := c.api.CurrentMasterchainInfo(ctx)
	if err != nil {
		return nil, fmt.Errorf("failed to get masterchain info: %w", err)
	}

	account, err := c.api.WaitForBlock(block.SeqNo).GetAccount(ctx, block, addr)
	if err != nil {
		return nil, fmt.Errorf("failed to get account: %w", err)
	}

	if !account.IsActive || account.LastTxLT == 0 {
		return nil, nil
	}

	batchSize := limit
	if batchSize > 16 {
		batchSize = 16
	}

	var allTx []*tlb.Transaction
	nextLT := account.LastTxLT
	nextHash := account.LastTxHash

	var lastBatch []*tlb.Transaction
	slideFrom := 0

	for uint32(len(allTx)) < limit && nextLT != 0 {
		txs, err := c.api.ListTransactions(ctx, addr, batchSize, nextLT, nextHash)
		if err != nil {
			if errors.Is(err, ton.ErrNoTransactionsWereFound) {
				break
			}
			var lsErr ton.LSError
			if errors.As(err, &lsErr) && lsErr.Code == -400 {
				if len(lastBatch) == 0 || slideFrom >= len(lastBatch)-1 {
					break
				}
				// Slide forward through the last successful batch: try the PrevTxLT of
				// progressively newer transactions (slideFrom+1, slideFrom+2, ..., len-1).
				// Each step gives a more recent (less likely to be GC'd) anchor.
				found := false
				for i := slideFrom + 1; i < len(lastBatch); i++ {
					if lastBatch[i].PrevTxLT == 0 {
						continue
					}
					nextLT = lastBatch[i].PrevTxLT
					nextHash = lastBatch[i].PrevTxHash
					slideFrom = i
					found = true
					break
				}
				if !found {
					break
				}
				continue
			}
			return nil, fmt.Errorf("failed to list transactions: %w", err)
		}

		lastBatch = txs
		slideFrom = 0

		for i := len(txs) - 1; i >= 0; i-- {
			allTx = append(allTx, txs[i])
			if uint32(len(allTx)) >= limit {
				break
			}
		}

		if uint32(len(allTx)) >= limit {
			break
		}

		nextLT = txs[0].PrevTxLT
		nextHash = txs[0].PrevTxHash
	}

	records := make([]*TxRecord, 0, len(allTx))
	for _, tx := range allTx {
		record := &TxRecord{
			Hash:      hex.EncodeToString(tx.Hash),
			Timestamp: time.Unix(int64(tx.Now), 0),
		}

		record.Fee = tx.TotalFees.Coins.Nano().String()

		if tx.IO.In != nil && tx.IO.In.MsgType == tlb.MsgTypeInternal {
			if internal := tx.IO.In.AsInternal(); internal != nil {
				record.Type = "incoming"
				record.Amount = internal.Amount.String()
				if internal.SrcAddr != nil {
					record.From = internal.SrcAddr.Bounce(false).String()
				}
				if c := internal.Comment(); c != "" {
					record.Comment = c
				}
			}
		}

		if tx.IO.Out != nil {
			msgs, err := tx.IO.Out.ToSlice()
			if err == nil {
				for _, outMsg := range msgs {
					if outMsg.MsgType != tlb.MsgTypeInternal {
						continue
					}
					if internal := outMsg.AsInternal(); internal != nil {
						record.Type = "outgoing"
						record.Amount = internal.Amount.String()
						if internal.DstAddr != nil {
							record.To = internal.DstAddr.Bounce(false).String()
						}
						if c := internal.Comment(); c != "" {
							record.Comment = c
						}
					}
				}
			}
		}

		records = append(records, record)
	}

	if uint32(len(records)) < limit && c.toncenter != nil {
		tcRecords, err := c.toncenter.GetTransactions(ctx, addr, limit)
		if err == nil {
			records = mergeTxRecords(records, tcRecords, int(limit))
		}
	}

	return records, nil
}

func mergeTxRecords(a, b []*TxRecord, limit int) []*TxRecord {
	seen := make(map[string]bool, len(a)+len(b))
	merged := make([]*TxRecord, 0, min(len(a)+len(b), limit))

	for _, r := range a {
		if seen[r.Hash] {
			continue
		}
		seen[r.Hash] = true
		merged = append(merged, r)
		if len(merged) >= limit {
			return merged
		}
	}

	for _, r := range b {
		if seen[r.Hash] {
			continue
		}
		seen[r.Hash] = true
		merged = append(merged, r)
		if len(merged) >= limit {
			return merged
		}
	}

	return merged
}

func (c *Client) IsContractDeployed(ctx context.Context, addr *address.Address) (bool, error) {
	ctx, cancel := context.WithTimeout(ctx, 10*time.Second)
	defer cancel()

	block, err := c.api.CurrentMasterchainInfo(ctx)
	if err != nil {
		return false, err
	}

	account, err := c.api.GetAccount(ctx, block, addr)
	if err != nil {
		return false, err
	}

	return account.IsActive, nil
}
