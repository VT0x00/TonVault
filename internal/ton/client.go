package ton

import (
	"context"
	"fmt"
	"time"

	"github.com/xssnick/tonutils-go/liteclient"
	"github.com/xssnick/tonutils-go/ton"
)

type Client struct {
	api            ton.APIClientWrapped
	connectionPool *liteclient.ConnectionPool
	configURL      string
}

var NetworkConfigURLs = map[string]string{
	"mainnet": "https://ton-blockchain.github.io/global.config.json",
	"testnet": "https://ton-blockchain.github.io/testnet-global.config.json",
}

func NewClientForNetwork(ctx context.Context, network string) (*Client, error) {
	if network == "" {
		network = "mainnet"
	}
	url, ok := NetworkConfigURLs[network]
	if !ok {
		return nil, fmt.Errorf("unknown network: %s", network)
	}
	return NewClient(ctx, url)
}

func NewClient(ctx context.Context, configURL string) (*Client, error) {
	pool := liteclient.NewConnectionPool()

	if err := pool.AddConnectionsFromConfigUrl(ctx, configURL); err != nil {
		return nil, fmt.Errorf("failed to connect to lite servers: %w", err)
	}

	api := ton.NewAPIClient(pool, ton.ProofCheckPolicyFast)

	cfg, err := liteclient.GetConfigFromUrl(ctx, configURL)
	if err != nil {
		pool.Stop()
		return nil, fmt.Errorf("failed to load network config: %w", err)
	}
	api.SetTrustedBlockFromConfig(cfg)

	return &Client{
		api:            api,
		connectionPool: pool,
		configURL:      configURL,
	}, nil
}

func (c *Client) GetAPI() ton.APIClientWrapped {
	return c.api
}

func (c *Client) Close() {
	c.connectionPool.Stop()
}

func (c *Client) GetMasterchainInfo(ctx context.Context) (*ton.BlockIDExt, error) {
	ctx, cancel := context.WithTimeout(ctx, 10*time.Second)
	defer cancel()
	return c.api.CurrentMasterchainInfo(ctx)
}
