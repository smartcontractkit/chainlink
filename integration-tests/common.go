package integration_tests

import (
	"github.com/smartcontractkit/chainlink-testing-framework/blockchain"
	"github.com/smartcontractkit/chainlink-testing-framework/config"
	"time"
)

// Just use os.GetEnv to get any secrets there, you don't need to fill URLs, that is handled by environment when you connect

var (
	DefaultGethSettings = &config.ETHNetwork{
		Name:    "Ethereum Geth dev",
		Type:    blockchain.SimulatedEthNetwork,
		ChainID: 1337,
		PrivateKeys: []string{
			"ac0974bec39a17e36ba4a6b4d238ff944bacb478cbed5efcae784d7bf4f2ff80",
			"59c6995e998f97a5a0044966f0945389dc9e86dae88c7a8412f4603b6b78690d",
			"5de4111afa1a4b94908f83103eb1f1706367c2e68ca870fc3fb9a804cdab365a",
		},
		ChainlinkTransactionLimit: 500000,
		Timeout:                   2 * time.Minute,
		MinimumConfirmations:      1,
		GasEstimationBuffer:       10000,
		BlockGasLimit:             40000000,
	}
)
