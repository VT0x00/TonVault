package models

import "time"

type Wallet struct {
	ID            string    `json:"id"`
	Name          string    `json:"name"`
	Address       string    `json:"address"`
	PublicKey     string    `json:"public_key"`
	EncryptedSeed string    `json:"encrypted_seed"`
	Version       string    `json:"version"`
	Network       string    `json:"network"`
	SubwalletID   uint32    `json:"subwallet_id"`
	CreatedAt     time.Time `json:"created_at"`
	IsDefault     bool      `json:"is_default"`
}

type Transaction struct {
	Hash        string    `json:"hash"`
	Type        string    `json:"type"`
	AssetType   string    `json:"asset_type"`
	AssetSymbol string    `json:"asset_symbol"`
	Amount      string    `json:"amount"`
	From        string    `json:"from"`
	To          string    `json:"to"`
	Comment     string    `json:"comment"`
	Fee         string    `json:"fee"`
	Status      string    `json:"status"`
	Timestamp   time.Time `json:"timestamp"`
	BlockHeight uint64    `json:"block_height"`
}

type Config struct {
	Network         string   `json:"network"`
	DefaultWalletID string   `json:"default_wallet_id"`
	ExplorerURL     string   `json:"explorer_url"`
	LogLevel        string   `json:"log_level"`
	LiteServers     []string `json:"lite_servers"`
}
