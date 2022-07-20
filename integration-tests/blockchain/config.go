package blockchain

import (
	"fmt"
	"strings"
	"time"

	"github.com/kelseyhightower/envconfig"
	"github.com/rs/zerolog/log"
)

var (
	// SimulatedEVMNetwork ensures that the test will use a default simulated geth instance
	SimulatedEVMNetwork = &EVMNetwork{
		Name:      "Simulated Geth",
		Simulated: true,
		ChainID:   1337,
		PrivateKeys: []string{
			"ac0974bec39a17e36ba4a6b4d238ff944bacb478cbed5efcae784d7bf4f2ff80",
			"59c6995e998f97a5a0044966f0945389dc9e86dae88c7a8412f4603b6b78690d",
			"5de4111afa1a4b94908f83103eb1f1706367c2e68ca870fc3fb9a804cdab365a",
		},
		ChainlinkTransactionLimit: 500000,
		Timeout:                   2 * time.Minute,
		MinimumConfirmations:      1,
		GasEstimationBuffer:       10000,
	}
)

// EVMNetwork configures all the data the test needs to connect and operate on an EVM compatible network
type EVMNetwork struct {
	// Human-readable name of the network:
	Name string `envconfig:"evm_name" default:"Unnamed EVM Network"`
	// Chain ID for the blockchain
	ChainID int64 `envconfig:"evm_chain_id" default:"1337"`
	// List of websocket URLs you want to connect to
	URLs []string `envconfig:"evm_urls" default:"ws://example.url"`
	// True if the network is simulated like a geth instance in dev mode. False if the network is a real test or mainnet
	Simulated bool `envconfig:"evm_simulated" default:"false"`
	// List of private keys to fund the tests
	PrivateKeys []string `envconfig:"evm_private_keys" default:"examplePrivateKey"`
	// Default gas limit to assume that Chainlink nodes will use. Used to try to estimate the funds that Chainlink
	// nodes require to run the tests.
	ChainlinkTransactionLimit uint64 `envconfig:"evm_chainlink_transaction_limit" default:"500000"`
	// How long to wait for on-chain operations before timing out an on-chain operation
	Timeout time.Duration `envconfig:"evm_transaction_timeout" default:"2m"`
	// How many block confirmations to wait to confirm on-chain events
	MinimumConfirmations int `envconfig:"evm_minimum_confirmations" default:"1"`
	// How much WEI to add to gas estimations for sending transactions
	GasEstimationBuffer uint64 `envconfig:"evm_gas_estimation_buffer" default:"1000"`

	// Only used internally, do not set
	URL string `ignored:"true"`
}

// LoadNetworkFromEnvironment loads an EVM network from default environment variables. Helpful in soak tests
func LoadNetworkFromEnvironment() *EVMNetwork {
	var network EVMNetwork
	if err := envconfig.Process("", &network); err != nil {
		log.Fatal().Err(err).Msg("Error loading network settings from environment variables")
	}
	log.Debug().Str("Name", network.Name).Int64("Chain ID", network.ChainID).Msg("Loaded Network")
	return &network
}

// ToMap marshalls the network's values to a generic map, useful for setting env vars on instances like the remote runner
// Map Structure
// "envconfig_key": stringValue
func (e *EVMNetwork) ToMap() map[string]interface{} {
	return map[string]interface{}{
		"evm_name":                        e.Name,
		"evm_chain_id":                    fmt.Sprint(e.ChainID),
		"evm_urls":                        strings.Join(e.URLs, ","),
		"evm_simulated":                   fmt.Sprint(e.Simulated),
		"evm_private_keys":                strings.Join(e.PrivateKeys, ","),
		"evm_chainlink_transaction_limit": fmt.Sprint(e.ChainlinkTransactionLimit),
		"evm_transaction_timeout":         fmt.Sprint(e.Timeout),
		"evm_minimum_confirmations":       fmt.Sprint(e.MinimumConfirmations),
		"evm_gas_estimation_buffer":       fmt.Sprint(e.GasEstimationBuffer),
	}
}

// ChainlinkValuesMap is a convenience function that marshalls the Chain ID and Chain URL into Chainlink Env var
// viable map
func (e *EVMNetwork) ChainlinkValuesMap() map[string]interface{} {
	valueMap := map[string]interface{}{}
	if !e.Simulated {
		valueMap["eth_url"] = e.URLs[0]
		valueMap["eth_chain_id"] = fmt.Sprint(e.ChainID)
	}
	return valueMap
}
