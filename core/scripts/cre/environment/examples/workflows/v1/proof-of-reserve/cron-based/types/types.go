package types

import (
	"github.com/ethereum/go-ethereum/common"

	"github.com/smartcontractkit/chainlink-common/pkg/workflows/sdk"
)

type WorkflowConfig struct {
	// name of the secret that stores authentication key
	AuthKeySecretName string `yaml:"auth_key_secret_name"`
	ChainSelector     uint64 `yaml:"chain_selector,omitempty"`
	ChainFamily       string `yaml:"chain_family,omitempty"`
	ChainID           string `yaml:"chain_id,omitempty"`
	BalanceReaderConfig
	ComputeConfig
}

type BalanceReaderConfig struct {
	BalanceReaderAddress string           `yaml:"balance_reader_address"`
	AddressesToRead      []common.Address `yaml:"addresses_to_read,omitempty"`
}

type ComputeConfig struct {
	FeedID                string          `yaml:"feed_id"`
	URL                   string          `yaml:"url"`
	DataFeedsCacheAddress string          `yaml:"consumer_address"`
	WriteTargetName       string          `yaml:"write_target_name"`
	AuthKey               sdk.SecretValue `yaml:"auth_key_secret_name,omitempty"`
}
