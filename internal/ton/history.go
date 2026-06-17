package ton

import (
	"context"
	"encoding/hex"
	"fmt"
	"time"

	"github.com/xssnick/tonutils-go/address"
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

	txs, err := c.api.ListTransactions(ctx, addr, limit, 0, nil)
	if err != nil {
		return nil, fmt.Errorf("failed to list transactions: %w", err)
	}

	var records []*TxRecord
	for _, tx := range txs {
		record := &TxRecord{
			Hash:      hex.EncodeToString(tx.Hash),
			Timestamp: time.Unix(int64(tx.Now), 0),
		}

		record.Fee = tx.TotalFees.Coins.Nano().String()

		if tx.IO.In != nil {
			if internal := tx.IO.In.AsInternal(); internal != nil {
				record.Type = "incoming"
				record.Amount = internal.Amount.String()
				if internal.SrcAddr != nil {
					record.From = internal.SrcAddr.String()
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
					if internal := outMsg.AsInternal(); internal != nil {
						record.Type = "outgoing"
						record.Amount = internal.Amount.String()
						if internal.DstAddr != nil {
							record.To = internal.DstAddr.String()
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

	return records, nil
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
