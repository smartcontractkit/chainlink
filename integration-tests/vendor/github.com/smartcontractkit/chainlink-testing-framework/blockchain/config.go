package blockchain

import (
	"encoding/json"
	"errors"
	"fmt"
	"strings"
	"time"

	"github.com/kelseyhightower/envconfig"
	"github.com/rs/zerolog/log"
)

var (
	// SimulatedEVMNetwork ensures that the test will use a default simulated geth instance
	SimulatedEVMNetwork = EVMNetwork{
		Name:                 "Simulated Geth",
		ClientImplementation: EthereumClientImplementation,
		Simulated:            true,
		ChainID:              1337,
		URLs:                 []string{"ws://geth:8546"},
		HTTPURLs:             []string{"http://geth:8544"},
		PrivateKeys: []string{
			"ac0974bec39a17e36ba4a6b4d238ff944bacb478cbed5efcae784d7bf4f2ff80",
			"59c6995e998f97a5a0044966f0945389dc9e86dae88c7a8412f4603b6b78690d",
			"5de4111afa1a4b94908f83103eb1f1706367c2e68ca870fc3fb9a804cdab365a",
			"8d351f5fc88484c65c15d44c7d3aa8779d12383373fb42d802e4576a50f765e5",
			"44fd8327d465031c71b20d7a5ba60bb01d33df8256fba406467bcb04e6f7262c",
			"809871f5c72d01a953f44f65d8b7bd0f3e39aee084d8cd0bc17ba3c386391814",
			"f29f5fda630ac9c0e39a8b05ec5b4b750a2e6ef098e612b177c6641bb5a675e1",
			"99b256477c424bb0102caab28c1792a210af906b901244fa67e2b704fac5a2bb",
			"bb74c3a9439ca83d09bcb4d3e5e65d8bc4977fc5b94be4db73772b22c3ff3d1a",
			"58845406a51d98fb2026887281b4e91b8843bbec5f16b89de06d5b9a62b231e8",
		},
		ChainlinkTransactionLimit: 500000,
		Timeout:                   JSONStrDuration{2 * time.Minute},
		MinimumConfirmations:      1,
		GasEstimationBuffer:       10000,
	}
)

// EVMNetwork configures all the data the test needs to connect and operate on an EVM compatible network
type EVMNetwork struct {
	// Human-readable name of the network:
	Name string `envconfig:"evm_name" default:"Unnamed EVM Network" toml:"evm_name" json:"evm_name"`
	// Chain ID for the blockchain
	ChainID int64 `envconfig:"evm_chain_id" default:"1337" toml:"evm_chain_id" json:"evm_chain_id"`
	// List of websocket URLs you want to connect to
	URLs []string `envconfig:"evm_urls" default:"ws://example.url" toml:"evm_urls" json:"evm_urls"`
	// List of websocket URLs you want to connect to
	HTTPURLs []string `envconfig:"evm_http_urls" default:"http://example.url" toml:"evm_http_urls" json:"evm_http_urls"`
	// True if the network is simulated like a geth instance in dev mode. False if the network is a real test or mainnet
	Simulated bool `envconfig:"evm_simulated" default:"false" toml:"evm_simulated" json:"evm_simulated"`
	// List of private keys to fund the tests
	PrivateKeys []string `envconfig:"evm_keys" default:"examplePrivateKey" toml:"evm_keys" json:"evm_keys"`
	// Default gas limit to assume that Chainlink nodes will use. Used to try to estimate the funds that Chainlink
	// nodes require to run the tests.
	ChainlinkTransactionLimit uint64 `envconfig:"evm_chainlink_transaction_limit" default:"500000" toml:"evm_chainlink_transaction_limit" json:"evm_chainlink_transaction_limit"`
	// How long to wait for on-chain operations before timing out an on-chain operation
	Timeout JSONStrDuration `envconfig:"evm_transaction_timeout" default:"2m" toml:"evm_transaction_timeout" json:"evm_transaction_timeout"`
	// How many block confirmations to wait to confirm on-chain events
	MinimumConfirmations int `envconfig:"evm_minimum_confirmations" default:"1" toml:"evm_minimum_confirmations" json:"evm_minimum_confirmations"`
	// How much WEI to add to gas estimations for sending transactions
	GasEstimationBuffer uint64 `envconfig:"evm_gas_estimation_buffer" default:"1000" toml:"evm_gas_estimation_buffer" json:"evm_gas_estimation_buffer"`
	// ClientImplementation is the blockchain client to use when interacting with the test chain
	ClientImplementation ClientImplementation `envconfig:"client_implementation" default:"Ethereum" toml:"client_implementation" json:"client_implementation"`

	// Only used internally, do not set
	URL string `ignored:"true"`
}

// LoadNetworkFromEnvironment loads an EVM network from default environment variables. Helpful in soak tests
func LoadNetworkFromEnvironment() EVMNetwork {
	var network EVMNetwork
	if err := envconfig.Process("", &network); err != nil {
		log.Fatal().Err(err).Msg("Error loading network settings from environment variables")
	}
	log.Debug().Str("Name", network.Name).Int64("Chain ID", network.ChainID).Msg("Loaded Network")
	return network
}

// ToMap marshalls the network's values to a generic map, useful for setting env vars on instances like the remote runner
// Map Structure
// "envconfig_key": stringValue
func (e *EVMNetwork) ToMap() map[string]interface{} {
	return map[string]interface{}{
		"evm_name":                        e.Name,
		"evm_chain_id":                    fmt.Sprint(e.ChainID),
		"evm_urls":                        strings.Join(e.URLs, ","),
		"evm_http_urls":                   strings.Join(e.HTTPURLs, ","),
		"evm_simulated":                   fmt.Sprint(e.Simulated),
		"evm_keys":                        strings.Join(e.PrivateKeys, ","),
		"evm_chainlink_transaction_limit": fmt.Sprint(e.ChainlinkTransactionLimit),
		"evm_transaction_timeout":         fmt.Sprint(e.Timeout),
		"evm_minimum_confirmations":       fmt.Sprint(e.MinimumConfirmations),
		"evm_gas_estimation_buffer":       fmt.Sprint(e.GasEstimationBuffer),
		"client_implementation":           fmt.Sprint(e.ClientImplementation),
	}
}

var (
	evmNetworkTOML = `[[EVM]]
ChainID = '%d'
MinContractPayment = '0'
%s`

	evmNodeTOML = `[[EVM.Nodes]]
Name = '%s'
WSURL = '%s'
HTTPURL = '%s'`
)

// MustChainlinkTOML marshals EVM network values into a TOML setting snippet. Will fail if error is encountered
// Can provide more detailed config for the network if non-default behaviors are desired.
func (e *EVMNetwork) MustChainlinkTOML(networkDetails string) string {
	if len(e.HTTPURLs) != len(e.URLs) || len(e.HTTPURLs) == 0 || len(e.URLs) == 0 {
		log.Fatal().
			Int("WS Count", len(e.URLs)).
			Int("HTTP Count", len(e.HTTPURLs)).
			Interface("WS URLs", e.URLs).
			Interface("HTTP URLs", e.HTTPURLs).
			Msg("Amount of HTTP and WS URLs should match, and not be empty")
		return ""
	}
	netString := fmt.Sprintf(evmNetworkTOML, e.ChainID, networkDetails)
	for index := range e.URLs {
		netString = fmt.Sprintf("%s\n\n%s", netString,
			fmt.Sprintf(evmNodeTOML, fmt.Sprintf("node-%d", index), e.URLs[index], e.HTTPURLs[index]))
	}

	return netString
}

// JSONStrDuration is JSON friendly duration that can be parsed from "1h2m0s" Go format
type JSONStrDuration struct {
	time.Duration
}

func (d *JSONStrDuration) MarshalJSON() ([]byte, error) {
	return json.Marshal(d.String())
}

func (d *JSONStrDuration) UnmarshalJSON(b []byte) error {
	var v interface{}
	if err := json.Unmarshal(b, &v); err != nil {
		return err
	}
	switch value := v.(type) {
	case string:
		var err error
		d.Duration, err = time.ParseDuration(value)
		if err != nil {
			return err
		}
		return nil
	default:
		return errors.New("invalid duration")
	}
}
