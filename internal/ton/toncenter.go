package ton

import (
	"context"
	"encoding/base64"
	"encoding/hex"
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"net/url"
	"strconv"
	"strings"
	"time"

	"github.com/xssnick/tonutils-go/address"
)

type toncenterTx struct {
	Hash  string `json:"hash"`
	UTime int64  `json:"utime"`
	Fee   string `json:"fee"`
	TransactionID struct {
		LT string `json:"lt"`
	} `json:"transaction_id"`
	InMsg   *toncenterMsg   `json:"in_msg"`
	OutMsgs []toncenterMsg  `json:"out_msgs"`
}

type toncenterMsg struct {
	Source      string `json:"source"`
	Destination string `json:"destination"`
	Value       string `json:"value"`
	MsgData     struct {
		Type string `json:"@type"`
		Text string `json:"text"`
		Body string `json:"body"`
	} `json:"msg_data"`
}

type toncenterResponse struct {
	OK     bool           `json:"ok"`
	Result []toncenterTx  `json:"result"`
	Error  string         `json:"error"`
}

type TonCenterClient struct {
	apiKey  string
	baseURL string
	http    *http.Client
}

func NewTonCenterClient(apiKey, network string) *TonCenterClient {
	base := "https://toncenter.com/api/v2"
	if network == "testnet" || network == "testnet2" {
		base = "https://testnet.toncenter.com/api/v2"
	}
	return &TonCenterClient{
		apiKey:  apiKey,
		baseURL: base,
		http:    &http.Client{Timeout: 30 * time.Second},
	}
}

func (t *TonCenterClient) GetTransactions(ctx context.Context, addr *address.Address, limit uint32) ([]*TxRecord, error) {
	u, err := url.Parse(t.baseURL + "/getTransactions")
	if err != nil {
		return nil, err
	}
	q := u.Query()
	q.Set("address", addr.Bounce(false).String())
	q.Set("limit", strconv.Itoa(int(limit)))
	q.Set("archival", "true")
	if t.apiKey != "" {
		q.Set("api_key", t.apiKey)
	}
	u.RawQuery = q.Encode()

	req, err := http.NewRequestWithContext(ctx, "GET", u.String(), nil)
	if err != nil {
		return nil, fmt.Errorf("failed to create request: %w", err)
	}

	resp, err := t.http.Do(req)
	if err != nil {
		return nil, fmt.Errorf("toncenter request failed: %w", err)
	}
	defer resp.Body.Close()

	body, err := io.ReadAll(resp.Body)
	if err != nil {
		return nil, fmt.Errorf("failed to read toncenter response: %w", err)
	}

	var apiResp toncenterResponse
	if err := json.Unmarshal(body, &apiResp); err != nil {
		return nil, fmt.Errorf("failed to parse toncenter response: %w", err)
	}

	if !apiResp.OK {
		return nil, fmt.Errorf("toncenter API error: %s", apiResp.Error)
	}

	ourAddr := addr.Bounce(false).String()
	seen := make(map[string]bool)
	records := make([]*TxRecord, 0, len(apiResp.Result))

	for _, tx := range apiResp.Result {
		record := parseTonCenterTx(tx, ourAddr)
		if record == nil {
			continue
		}
		if seen[record.Hash] {
			continue
		}
		seen[record.Hash] = true
		records = append(records, record)
	}
	return records, nil
}

func parseTonCenterTx(tx toncenterTx, ourAddr string) *TxRecord {
	hashBytes, err := base64.RawURLEncoding.DecodeString(tx.Hash)
	if err != nil {
		hashBytes, err = base64.StdEncoding.DecodeString(tx.Hash)
		if err != nil {
			return nil
		}
	}

	r := &TxRecord{
		Hash:      hex.EncodeToString(hashBytes),
		Timestamp: time.Unix(tx.UTime, 0),
		Fee:       tx.Fee,
	}

	if tx.InMsg == nil {
		return r
	}

	src := strings.TrimSpace(tx.InMsg.Source)
	dst := strings.TrimSpace(tx.InMsg.Destination)

	if dst == ourAddr {
		r.Type = "incoming"
		r.From = src
		r.Amount = tx.InMsg.Value
		r.Comment = decodeMsgData(tx.InMsg.MsgData)
	}

	if src == ourAddr {
		r.Type = "outgoing"
		r.Amount = tx.InMsg.Value
		r.Comment = decodeMsgData(tx.InMsg.MsgData)
		if len(tx.OutMsgs) > 0 {
			r.From = src
			r.To = tx.OutMsgs[0].Destination
			if comment := decodeMsgData(tx.OutMsgs[0].MsgData); comment != "" {
				r.Comment = comment
			}
		}
	}

	return r
}

func decodeMsgData(d struct {
	Type string `json:"@type"`
	Text string `json:"text"`
	Body string `json:"body"`
}) string {
	raw := d.Text
	if raw == "" {
		raw = d.Body
	}
	if raw == "" {
		return ""
	}
	decoded, err := base64.StdEncoding.DecodeString(raw)
	if err != nil {
		return raw
	}
	return strings.TrimSpace(string(decoded))
}
