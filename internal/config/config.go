package config

import (
	"errors"
	"os"
	"path/filepath"

	"github.com/spf13/viper"
)

const (
	DefaultNetwork     = "mainnet"
	DefaultLogLevel    = "info"
	DefaultExplorerURL = "https://tonviewer.com"
	AppName            = "tonvault"
)

var DefaultLiteServers = map[string][]string{
	"mainnet": {
		"https://ton-blockchain.github.io/global.config.json",
	},
	"testnet": {
		"https://ton-blockchain.github.io/testnet-global.config.json",
	},
}

type Manager struct {
	viper *viper.Viper
	path  string
}

func NewManager() (*Manager, error) {
	home, err := os.UserHomeDir()
	if err != nil {
		return nil, errors.New("cannot determine home directory")
	}

	configDir := filepath.Join(home, ".config", AppName)
	if err := os.MkdirAll(configDir, 0700); err != nil {
		return nil, err
	}

	v := viper.New()
	v.SetConfigName("config")
	v.SetConfigType("json")
	v.AddConfigPath(configDir)

	v.SetDefault("network", DefaultNetwork)
	v.SetDefault("log_level", DefaultLogLevel)
	v.SetDefault("explorer_url", DefaultExplorerURL)
	v.SetDefault("lite_servers", DefaultLiteServers[DefaultNetwork])
	v.SetDefault("default_wallet_id", "")
	v.SetDefault("toncenter_api_key", "")

	m := &Manager{
		viper: v,
		path:  filepath.Join(configDir, "config.json"),
	}

	if err := v.ReadInConfig(); err != nil {
		var notFound viper.ConfigFileNotFoundError
		if !errors.As(err, &notFound) {
			return nil, err
		}
		if err := v.SafeWriteConfig(); err != nil {
			return nil, err
		}
	}

	return m, nil
}

func (m *Manager) GetNetwork() string {
	return m.viper.GetString("network")
}

func (m *Manager) SetNetwork(net string) error {
	m.viper.Set("network", net)
	m.viper.Set("lite_servers", DefaultLiteServers[net])
	if net == "testnet" {
		m.viper.Set("explorer_url", "https://testnet.tonviewer.com")
	} else {
		m.viper.Set("explorer_url", DefaultExplorerURL)
	}
	return m.viper.WriteConfig()
}

func (m *Manager) GetDefaultWalletID() string {
	return m.viper.GetString("default_wallet_id")
}

func (m *Manager) SetDefaultWalletID(id string) error {
	m.viper.Set("default_wallet_id", id)
	return m.viper.WriteConfig()
}

func (m *Manager) GetExplorerURL() string {
	return m.viper.GetString("explorer_url")
}

func (m *Manager) GetLogLevel() string {
	return m.viper.GetString("log_level")
}

func (m *Manager) GetLiteServers() []string {
	return m.viper.GetStringSlice("lite_servers")
}

func (m *Manager) GetTonCenterAPIKey() string {
	return m.viper.GetString("toncenter_api_key")
}

func (m *Manager) Set(key, value string) error {
	m.viper.Set(key, value)
	return m.viper.WriteConfig()
}

func (m *Manager) Get(key string) string {
	return m.viper.GetString(key)
}

func (m *Manager) Reset() error {
	cfg := map[string]interface{}{
		"network":           DefaultNetwork,
		"default_wallet_id": "",
		"explorer_url":      DefaultExplorerURL,
		"log_level":         DefaultLogLevel,
		"lite_servers":      DefaultLiteServers[DefaultNetwork],
		"toncenter_api_key": "",
	}
	for k, v := range cfg {
		m.viper.Set(k, v)
	}
	return m.viper.WriteConfig()
}

func (m *Manager) Path() string {
	return m.path
}
