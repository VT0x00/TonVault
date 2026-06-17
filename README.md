# TonVault

A feature-rich CLI wallet for the TON blockchain. Manage wallets, check balances, send tokens, and inspect transaction history — all from the terminal.

## Features

- **Wallet management**: create, import, list, export, delete wallets (V3, V4, V5R1)
- **Balance checks**: TON balance and Jetton token balances
- **Send tokens**: transfer TON and Jettons with optional comments
- **Transaction history**: inspect any TON address's transaction history (even wallets not created in TonVault)
- **Multi-network**: mainnet and testnet support
- **Encrypted storage**: seed phrases encrypted with AES-256-GCM

## Installation

### From source

```bash
git clone https://github.com/VT0x00/tonvault.git
cd tonvault
go build -o tonvault .
sudo mv tonvault /usr/local/bin/
```

### Download a pre-built binary

Download the latest release for your platform from the [releases page](https://github.com/VT0x00/tonvault/releases), extract it, and place the binary in your `PATH`.

## Quick start

```bash
# Create your first wallet
tonvault wallet create

# Set it as default and check your balance
tonvault balance show

# View transaction history for your wallet
tonvault history list

# Check any TON address (even from another wallet app)
tonvault history list --address EQD...
tonvault history list --address EQD... --limit 20 --network testnet

# Send TON
tonvault send ton EQD... --amount 1.5 --comment "thanks!"
```

## Command reference

### `tonvault wallet`

| Subcommand  | Description                              |
|-------------|------------------------------------------|
| `create`    | Create a new wallet with a seed phrase   |
| `import`    | Import wallet from a 24-word seed phrase |
| `list`      | List all locally stored wallets          |
| `info`      | Show wallet details (address, version)   |
| `export`    | Export wallet seed phrase                |
| `delete`    | Delete a wallet from local storage       |

**Wallet create / import options:**

| Flag           | Default    | Description                                     |
|----------------|------------|-------------------------------------------------|
| `--network`    | `""`       | Network: `mainnet` or `testnet`. Overrides the global `network` config for this wallet. |

```bash
# Create wallet for testnet
tonvault wallet create --network testnet

# Import wallet with network specified
tonvault wallet import --network mainnet
```

### `tonvault balance`

| Subcommand  | Description                       |
|-------------|-----------------------------------|
| `show`      | Show TON balance for a wallet     |
| `jettons`   | List Jetton token balances        |

### `tonvault send`

| Subcommand  | Description                       |
|-------------|-----------------------------------|
| `ton`       | Send TON coins to any address     |
| `jetton`    | Send Jetton tokens                |

### `tonvault history`

| Subcommand  | Description                                    |
|-------------|------------------------------------------------|
| `list`      | List recent transactions for an address        |
| `show`      | Show detailed information about a transaction  |

#### `history list` options

| Flag           | Default    | Description                                     |
|----------------|------------|-------------------------------------------------|
| `--address`    | `""`       | TON address to query (bypasses local wallet store). Query any address from any wallet app. |
| `--network`    | `""`       | Network: `mainnet` or `testnet`. Defaults to the wallet's configured network, or `mainnet`. |
| `--limit`      | `10`       | Number of transactions to fetch.                |

**Pagination:** Transactions are fetched in batches of up to 16 per request. If the lite server does not retain the full transaction history (common for older blocks), fewer transactions than requested may be returned — a warning is shown when this occurs.

#### `history show` options

| Flag           | Default    | Description                                     |
|----------------|------------|-------------------------------------------------|
| `--address`    | `""`       | TON address to query (bypasses local wallet store). |
| `--network`    | `""`       | Network: `mainnet` or `testnet`. Defaults to the wallet's configured network, or `mainnet`. |

Search the 100 most recent transactions of the wallet/address for a specific hash and display its full details (amount, counterparty, fee, timestamp, comment).

Examples:

```bash
# Transactions for the default wallet
tonvault history list

# Transactions for a specific stored wallet
tonvault history list my-wallet

# Transactions for any TON address (e.g. from Tonkeeper)
tonvault history list --address EQD4wM6J1S9sS5G8rV7jU7Xk2L3pQ5rN0mP1oI2uY3tR4eW5

# Transactions on testnet
tonvault history list --address EQD... --network testnet --limit 50

# Show details for a specific transaction
tonvault history show <tx_hash>

# Show details for a transaction on a specific address
tonvault history show --address EQD... <tx_hash>
```

### `tonvault network`

| Subcommand  | Description                       |
|-------------|-----------------------------------|
| `set`       | Switch between mainnet and testnet|
| `status`    | Show current network              |

### `tonvault config`

| Subcommand  | Description                       |
|-------------|-----------------------------------|
| `get`       | Get a configuration value         |
| `set`       | Set a configuration value         |
| `reset`     | Reset to default configuration    |

## Configuration

Config is stored at `~/.config/tonvault/config.json`.

| Key                  | Default      | Description                  |
|----------------------|--------------|------------------------------|
| `network`            | `mainnet`    | `mainnet` or `testnet`       |
| `default_wallet_id`  | `""`         | Default wallet for commands  |
| `explorer_url`       | (varies)     | Blockchain explorer URL      |
| `log_level`          | `info`       | Logging verbosity            |

## Transaction history for external wallets

TonVault can retrieve history for **any** TON address — not just wallets you created
inside TonVault. Use the `--address` flag to query addresses from Tonkeeper,
TonKeeper, Wallet, or any other TON wallet app:

```bash
tonvault history list --address EQD4wM6J1S9sS5G8rV7jU7Xk2L3pQ5rN0mP1oI2uY3tR4eW5
```

The output shows each transaction's date, type (incoming/outgoing), amount in TON,
counterparty address, and any attached comment/memo.

## Data storage

| Data            | Location                                    |
|-----------------|---------------------------------------------|
| Config          | `~/.config/tonvault/config.json`            |
| Wallets         | `~/.config/tonvault/wallets/wallets.json`   |
| Seed phrases    | Encrypted with AES-256-GCM in wallets.json  |

## Development

```bash
git clone https://github.com/VT0x00/tonvault.git
cd tonvault
go mod download
go build -o tonvault .
./tonvault --help
```

### Project structure

```
tonvault/
├── main.go              # Entry point
├── cmd/                 # CLI command definitions (cobra)
│   ├── root.go
│   ├── balance*.go
│   ├── config_*.go
│   ├── history*.go
│   ├── network*.go
│   ├── send*.go
│   └── wallet*.go
├── internal/
│   ├── config/          # Viper-based configuration manager
│   ├── models/          # Data types (Wallet, Transaction, Config)
│   ├── ton/             # TON blockchain interaction layer
│   │   ├── client.go    # Lite server connection
│   │   ├── balance.go   # Balance queries
│   │   ├── history.go   # Transaction history
│   │   ├── jettons.go   # Jetton support
│   │   └── transfer.go  # Send TON/Jettons
│   └── wallet/          # Wallet creation, storage, encryption
└── tonvault.md          # Technical specification
```

## License

MIT
