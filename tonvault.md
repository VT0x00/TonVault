# Technical Specification: TonVault CLI Wallet

## 1. Project Overview

**Project Name:** TonVault

**Type:** Command-Line Interface (CLI) application for The Open Network (TON) blockchain

**Language:** Go (Golang)

**Objective:** Develop a feature-rich CLI wallet for the TON blockchain that provides functionality comparable to major mobile wallets like Tonkeeper, including TON and Jetton transfers, transaction history, balance checking, and wallet management, all accessible from the terminal.

---

## 2. Core Features

### 2.1 Wallet Management

| Feature | Description |
|---------|-------------|
| **Create Wallet** | Generate a new wallet with a 24-word BIP39 mnemonic seed phrase |
| **Import Wallet** | Restore an existing wallet from a seed phrase or private key |
| **List Wallets** | Display all locally stored wallets with their addresses and aliases |
| **Delete Wallet** | Remove a wallet from local storage |
| **Export Wallet** | Export wallet private key or seed phrase (with confirmation) |
| **Wallet Info** | Display detailed information about a specific wallet (address, balance, version, etc.) |

### 2.2 Balance & Account Information

| Feature | Description |
|---------|-------------|
| **View Balance** | Display TON balance in both nanoTON and TON units |
| **View Jetton Balances** | List all Jetton tokens held by the wallet with their balances |
| **Account Details** | Show wallet address, public key, wallet version, and deployment status |
| **Multi-wallet Support** | Manage and switch between multiple wallets |

### 2.3 TON Transfers

| Feature | Description |
|---------|-------------|
| **Send TON** | Transfer TON coins to any TON address |
| **Add Comment** | Attach an optional text comment (memo) to the transaction |
| **Specify Amount** | Support for both TON and nanoTON units |
| **Bounce Option** | Specify whether the transaction should bounce if the recipient contract doesn't exist |
| **Transaction Confirmation** | Display transaction hash and link to blockchain explorer after sending |
| **Gas Estimation** | Estimate and display gas fees before confirming the transaction |

### 2.4 Jetton Transfers

| Feature | Description |
|---------|-------------|
| **Send Jettons** | Transfer Jetton tokens (e.g., USDT, NOT, etc.) to any TON address |
| **Add Comment** | Attach an optional text comment (memo) to Jetton transfers |
| **Jetton Selection** | Interactive selection of which Jetton to send from available balances |
| **Gas Estimation** | Estimate and display gas fees (paid in TON) before confirming |
| **Transaction Confirmation** | Display transaction hash and link to blockchain explorer |

### 2.5 Transaction History

| Feature | Description |
|---------|-------------|
| **List Transactions** | Display recent transaction history with pagination |
| **Transaction Details** | Show detailed information for a specific transaction (hash, amount, sender, recipient, timestamp, status, comment) |
| **Filter Transactions** | Filter by type (incoming/outgoing), asset (TON/Jetton), or date range |
| **Explorer Link** | Generate and display a link to view the transaction on TON blockchain explorers |

### 2.6 Additional Features

| Feature | Description |
|---------|-------------|
| **Network Selection** | Switch between Mainnet and Testnet |
| **Address Validation** | Validate TON addresses before sending |
| **QR Code Generation** | Generate QR code for wallet address (optional) |
| **Configuration Management** | Persistent configuration storage (network, default wallet, etc.) |

---

## 3. Non-Functional Requirements

### 3.1 Performance
- Transaction signing and submission should complete within 5 seconds under normal network conditions
- Balance checks should complete within 2 seconds
- Support for concurrent operations using Go's goroutines

### 3.2 Security
- Seed phrases and private keys MUST be stored encrypted
- Encryption using AES-256-GCM or similar industry-standard algorithms
- No plaintext storage of sensitive data
- Interactive confirmation for all transaction submissions
- Support for passphrase-protected wallets

### 3.3 Usability
- Intuitive command structure with subcommands (e.g., `tonvault wallet create`, `tonvault send ton`)
- Interactive mode for guided transactions
- Help commands and usage examples for all commands
- Color-coded output for better readability
- Progress indicators for long-running operations

### 3.4 Reliability
- Automatic retry with exponential backoff for failed requests
- Connection pooling and automatic failover between multiple lite servers
- Graceful handling of network errors
- Transaction status verification after submission

### 3.5 Maintainability
- Clean, modular architecture following Go best practices
- Comprehensive unit and integration tests
- Well-documented code and API
- Structured logging with configurable log levels

---

## 4. Technology Stack

### 4.1 Core Libraries

| Component | Technology | Purpose |
|-----------|------------|---------|
| **Blockchain Interaction** | `tonutils-go` (github.com/xssnick/tonutils-go) | Primary Go library for TON blockchain interaction with native ADNL and lite protocol support |
| **Alternative** | `tongo` (github.com/tonkeeper/tongo) | Alternative Go library with comprehensive TON primitives |
| **High-level API** | `tonapi-go` (github.com/tonkeeper/tonapi-go) | Optional high-level API for simplified Jetton and NFT operations |

### 4.2 Wallet & Cryptography

| Component | Technology | Purpose |
|-----------|------------|---------|
| **Mnemonic Generation** | BIP39 standard | Generate and validate 24-word seed phrases |
| **Key Derivation** | Ed25519 | TON uses Ed25519 for wallet key derivation |
| **Encryption** | AES-256-GCM | Encrypt sensitive data (seed phrases, private keys) |

### 4.3 CLI Framework

| Component | Technology | Purpose |
|-----------|------------|---------|
| **CLI Framework** | `cobra` (github.com/spf13/cobra) | Industry-standard CLI framework for Go |
| **Configuration** | `viper` (github.com/spf13/viper) | Configuration management with support for multiple formats |
| **Output Formatting** | `tablewriter` (github.com/olekukonko/tablewriter) | Formatted table output for transaction lists and balances |
| **Color Output** | `color` (github.com/fatih/color) | Colored terminal output for better UX |

### 4.4 Data Storage

| Component | Technology | Purpose |
|-----------|------------|---------|
| **Local Storage** | Encrypted JSON files | Store wallet data, configuration, and transaction cache |
| **Database (Optional)** | SQLite (with GORM) | Optional relational storage for transaction history |

### 4.5 TON-Specific Components

| Component | Technology | Purpose |
|-----------|------------|---------|
| **Wallet Versions** | V3, V4, V5R1, Highload V3 | Support for different TON wallet contract versions |
| **Jetton Standard** | TEP-74 | Implement Jetton transfers according to TEP-74 standard |
| **Transaction Comment** | Text comment field | Support for optional transaction comments/memos |
| **Network Config** | Global config files | Mainnet: `https://ton-blockchain.github.io/global.config.json` |

### 4.6 Development Tools

| Component | Technology | Purpose |
|-----------|------------|---------|
| **Build Tool** | Go modules | Dependency management and builds |
| **Testing** | Go standard testing + testify | Unit and integration testing |
| **Linting** | golangci-lint | Code quality and style enforcement |
| **CI/CD** | GitHub Actions | Automated testing and releases |

---

## 5. CLI Command Structure

```
tonvault
├── wallet
│   ├── create          # Create a new wallet
│   ├── import          # Import wallet from seed phrase
│   ├── list            # List all wallets
│   ├── info            # Show wallet details
│   ├── delete          # Delete a wallet
│   └── export          # Export wallet (private key/seed)
├── balance
│   ├── show            # Show TON balance
│   └── jettons         # Show Jetton balances
├── send
│   ├── ton             # Send TON coins
│   └── jetton          # Send Jettons
├── history
│   ├── list            # List transaction history
│   └── show            # Show transaction details
├── network
│   ├── set             # Set network (mainnet/testnet)
│   └── status          # Show current network status
└── config
    ├── set             # Set configuration values
    ├── get             # Get configuration values
    └── reset           # Reset to default configuration
```

---

## 6. Data Models

### 6.1 Wallet
```go
type Wallet struct {
    ID          string    `json:"id"`
    Name        string    `json:"name"`
    Address     string    `json:"address"`
    PublicKey   string    `json:"public_key"`
    EncryptedSeed string `json:"encrypted_seed"`
    Version     string    `json:"version"` // "v3", "v4", "v5r1"
    SubwalletID uint32    `json:"subwallet_id"`
    CreatedAt   time.Time `json:"created_at"`
    IsDefault   bool      `json:"is_default"`
}
```

### 6.2 Transaction
```go
type Transaction struct {
    Hash        string    `json:"hash"`
    Type        string    `json:"type"` // "incoming", "outgoing"
    AssetType   string    `json:"asset_type"` // "ton", "jetton"
    AssetSymbol string    `json:"asset_symbol"`
    Amount      string    `json:"amount"` // in nano units
    From        string    `json:"from"`
    To          string    `json:"to"`
    Comment     string    `json:"comment"`
    Fee         string    `json:"fee"`
    Status      string    `json:"status"` // "pending", "confirmed", "failed"
    Timestamp   time.Time `json:"timestamp"`
    BlockHeight uint64    `json:"block_height"`
}
```

### 6.3 Configuration
```go
type Config struct {
    Network         string `json:"network"` // "mainnet", "testnet"
    DefaultWalletID string `json:"default_wallet_id"`
    ExplorerURL     string `json:"explorer_url"`
    LogLevel        string `json:"log_level"`
    LiteServers     []string `json:"lite_servers"`
}
```

---

## 7. API Integration

### 7.1 Blockchain Connection Methods

| Method | Description | Use Case |
|--------|-------------|----------|
| **Lite Client (ADNL)** | Direct connection to TON lite servers | Primary method for all operations |
| **TON API (HTTP)** | REST API via TonAPI.io or self-hosted | Optional for high-level operations |

### 7.2 Key Operations

1. **Wallet Initialization**: Derive address and public key from seed phrase
2. **Balance Query**: Get balance using lite client API
3. **Transaction Building**: Construct external messages for sending TON/Jettons
4. **Transaction Signing**: Sign messages with wallet private key
5. **Transaction Broadcasting**: Send signed messages to the network
6. **Transaction History**: Query transaction history via lite client or TonAPI

---

## 8. Development Phases

### Phase 1: Foundation (Weeks 1-2)
- Project setup and Go module initialization
- CLI framework integration (Cobra)
- Configuration management (Viper)
- Basic wallet creation and import functionality
- Connection to TON network (lite client)

### Phase 2: Core Features (Weeks 3-4)
- Balance checking (TON and Jettons)
- TON transfer with comments
- Transaction signing and broadcasting
- Transaction history retrieval

### Phase 3: Jetton Support (Weeks 5-6)
- Jetton balance detection
- Jetton transfer implementation (TEP-74)
- Jetton-specific error handling
- Gas estimation for Jetton transfers

### Phase 4: Polish & Testing (Weeks 7-8)
- Comprehensive test coverage
- Error handling and edge cases
- Documentation and usage examples
- Performance optimization
- Security audit

---

## 9. Testing Strategy

| Test Type | Description |
|-----------|-------------|
| **Unit Tests** | Test individual components (wallet generation, address validation, etc.) |
| **Integration Tests** | Test against TON testnet for real transactions |
| **Mock Tests** | Mock blockchain responses for deterministic testing |
| **Security Tests** | Test encryption, key storage, and sensitive data handling |
| **CLI Tests** | Test command parsing, flags, and output formatting |

---

## 10. Deliverables

1. **Source Code**: Complete Go source code with modular architecture
2. **Binary Releases**: Compiled binaries for Linux, macOS, and Windows
3. **Documentation**:
   - Installation guide
   - User manual with command reference
   - Developer documentation
4. **Test Suite**: Comprehensive unit and integration tests
5. **Example Scripts**: Example usage scripts and automation examples

---

## 11. Security Considerations

- **Seed Phrase Storage**: Never store seed phrases in plaintext. Use encryption with a user-provided password.
- **Memory Safety**: Clear sensitive data from memory after use.
- **Input Validation**: Validate all user inputs (addresses, amounts, etc.).
- **Transaction Confirmation**: Always require explicit user confirmation before broadcasting transactions.
- **Error Messages**: Avoid exposing sensitive information in error messages.
- **Rate Limiting**: Implement rate limiting for API calls to avoid being blocked.

---

## 12. References

- TON Documentation: https://docs.ton.org
- TON SDKs: https://old-docs.ton.org/v3/guidelines/dapps/apis-sdks/sdk
- TON Wallet Contracts: https://old-docs.ton.org/v3/documentation/smart-contracts/contracts-specs/wallet-contract
- Jetton Standard (TEP-74): https://github.com/ton-blockchain/TEPs/blob/master/text/0074-jettons-standard.md
- TonAPI Documentation: https://docs.tonconsole.com/tonapi
