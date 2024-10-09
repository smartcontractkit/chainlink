package devenv

import (
	"encoding/json"
	"fmt"
	"os"

	"github.com/pelletier/go-toml/v2"
)

const (
	DefaultAFNPasphrase       = "my-not-so-secret-passphrase"
	DefaultRageProxy          = "127.0.0.1:8081"
	DefaultProxyListenAddress = "127.0.0.1:8080"
	DefaultDiscovererDbPath   = "/app/rageproxy-discoverer-db.json"
)

var (
	DefaultRageProxySharedConfig = ProxySharedConfig{
		Host: HostConfig{
			DurationBetweenDials: 1000000000,
		},
		Discoverer: DiscovererConfig{
			DeltaReconcile: 1000000000,
		},
	}
)

type Networking struct {
	RageProxy     string   `toml:"rageproxy"`
	Bootstrappers []string `toml:"bootstrappers"`
}

type HomeChain struct {
	Name                 string `toml:"name"`
	CapabilitiesRegistry string `toml:"capabilities_registry"`
	CCIPConfig           string `toml:"ccip_config"`
	RMNHome              string `toml:"rmn_home"`
}

type Stability struct {
	Type string `toml:"type"`
}

type RemoteChain struct {
	Name             string    `toml:"name"`
	Stability        Stability `toml:"stability"`
	StartBlockNumber int       `toml:"start_block_number"`
	OffRamp          string    `toml:"off_ramp"`
	RMNRemote        string    `toml:"rmn_remote"`
}

type SharedConfig struct {
	Networking   Networking    `toml:"networking"`
	HomeChain    HomeChain     `toml:"home_chain"`
	RemoteChains []RemoteChain `toml:"remote_chains"`
}

func (s SharedConfig) afn2ProxySharedConfigFile() (string, error) {
	data, err := toml.Marshal(s)
	if err != nil {
		return "", fmt.Errorf("failed to marshal afn2Proxy shared config: %w", err)
	}
	return CreateTempFile(data, "afn2proxy_shared")
}

type LocalConfig struct {
	Chains []Chain `toml:"chains"`
}

func (l LocalConfig) afn2ProxyLocalConfigFile() (string, error) {
	data, err := toml.Marshal(l)
	if err != nil {
		return "", fmt.Errorf("failed to marshal afn2Proxy local config: %w", err)
	}
	return CreateTempFile(data, "afn2proxy_local")
}

type Chain struct {
	Name string `toml:"name"`
	RPC  string `toml:"rpc"`
}

type ProxyLocalConfig struct {
	ListenAddresses   []string `json:"ListenAddresses"`
	AnnounceAddresses []string `json:"AnnounceAddresses"`
	ProxyAddress      string   `json:"ProxyAddress"`
	DiscovererDbPath  string   `json:"DiscovererDbPath"`
}

func (l ProxyLocalConfig) rageProxyLocal() (string, error) {
	data, err := json.Marshal(l)
	if err != nil {
		return "", fmt.Errorf("failed to marshal rageProxy local config: %w", err)
	}
	return CreateTempFile(data, "rageproxy_local")
}

type HostConfig struct {
	DurationBetweenDials int64 `json:"DurationBetweenDials"`
}

type DiscovererConfig struct {
	DeltaReconcile int64 `json:"DeltaReconcile"`
}

type ProxySharedConfig struct {
	Host       HostConfig       `json:"Host"`
	Discoverer DiscovererConfig `json:"Discoverer"`
}

func (s ProxySharedConfig) rageProxyShared() (string, error) {
	data, err := json.Marshal(s)
	if err != nil {
		return "", fmt.Errorf("failed to marshal rageProxy shared config: %w", err)
	}
	return CreateTempFile(data, "rageproxy_shared")
}

type RMNConfig struct {
	Shared      SharedConfig
	Local       LocalConfig
	ProxyLocal  ProxyLocalConfig
	ProxyShared ProxySharedConfig
}

func CreateTempFile(data []byte, pattern string) (string, error) {
	file, err := os.CreateTemp("", pattern)
	if err != nil {
		return "", fmt.Errorf("failed to create temp file for %s: %w", pattern, err)
	}
	_, err = file.Write(data)
	if err != nil {
		return "", fmt.Errorf("failed to write  %s: %w", pattern, err)
	}
	return file.Name(), nil
}
