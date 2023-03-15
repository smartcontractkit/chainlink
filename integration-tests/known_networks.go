package networks

import (
	"fmt"
	"os"
	"strings"
	"time"

	"github.com/rs/zerolog/log"

	"github.com/smartcontractkit/chainlink-testing-framework/blockchain"
	"github.com/smartcontractkit/chainlink-testing-framework/logging"
)

// Pre-configured test networks and their connections
// Some networks with public RPC endpoints are already filled out, but make use of environment variables to use info like
// private RPC endpoints and private keys.
var (
	// SelectedNetworks uses the SELECTED_NETWORKS env var to determine which network to run the test on.
	// For use in tests that utilize multiple chains. For tests on one chain, see SelectedNetwork
	// For CCIP use index 1 and 2 of SELECTED_NETWORKS to denote source and destination network respectively
	SelectedNetworks []blockchain.EVMNetwork = determineSelectedNetworks()
	// SelectedNetwork uses the first listed network in SELECTED_NETWORKS, for use in tests on only one chain
	SelectedNetwork blockchain.EVMNetwork = SelectedNetworks[0]

	// SimulatedEVM represents a simulated network
	SimulatedEVM blockchain.EVMNetwork = blockchain.SimulatedEVMNetwork
	// generalEVM is a customizable network through environment variables
	// This is getting little use, and causes some confusion. Can re-enable if people want it.
	// generalEVM blockchain.EVMNetwork = blockchain.LoadNetworkFromEnvironment()

	// SimulatedevmNonDev1 represents a simulated network which can be used to deploy a non-dev geth node
	SimulatedEVMNonDev1 = blockchain.EVMNetwork{
		Name:                 "source-chain",
		Simulated:            true,
		ClientImplementation: blockchain.EthereumClientImplementation,
		ChainID:              1337,
		PrivateKeys: []string{
			"ac0974bec39a17e36ba4a6b4d238ff944bacb478cbed5efcae784d7bf4f2ff80",
		},
		URLs:                      []string{"ws://source-chain-ethereum-geth:8546"},
		HTTPURLs:                  []string{"http://source-chain-ethereum-geth:8544"},
		ChainlinkTransactionLimit: 500000,
		Timeout:                   blockchain.JSONStrDuration{2 * time.Minute},
		MinimumConfirmations:      1,
		GasEstimationBuffer:       10000,
	}

	// SimulatedEVM_NON_DEV_2 represents a simulated network with chain id 2337 which can be used to deploy a non-dev geth node
	SimulatedEVMNonDev2 = blockchain.EVMNetwork{
		Name:                 "dest-chain",
		Simulated:            true,
		ClientImplementation: blockchain.EthereumClientImplementation,
		ChainID:              2337,
		PrivateKeys: []string{
			"ac0974bec39a17e36ba4a6b4d238ff944bacb478cbed5efcae784d7bf4f2ff80",
		},
		URLs:                      []string{"ws://dest-chain-ethereum-geth:8546"},
		HTTPURLs:                  []string{"http://dest-chain-ethereum-geth:8544"},
		ChainlinkTransactionLimit: 500000,
		Timeout:                   blockchain.JSONStrDuration{2 * time.Minute},
		MinimumConfirmations:      1,
		GasEstimationBuffer:       10000,
	}

	SimulatedEVMNonDev = blockchain.EVMNetwork{
		Name:                 "simulated",
		Simulated:            true,
		ClientImplementation: blockchain.EthereumClientImplementation,
		ChainID:              1337,
		PrivateKeys: []string{
			"ac0974bec39a17e36ba4a6b4d238ff944bacb478cbed5efcae784d7bf4f2ff80",
		},
		URLs:                      []string{"ws://simulated-ethereum-geth:8546"},
		HTTPURLs:                  []string{"http://simulated-ethereum-geth:8544"},
		ChainlinkTransactionLimit: 500000,
		Timeout:                   blockchain.JSONStrDuration{2 * time.Minute},
		MinimumConfirmations:      1,
		GasEstimationBuffer:       10000,
	}

	// sepoliaTestnet https://sepolia.dev/
	SepoliaTestnet blockchain.EVMNetwork = blockchain.EVMNetwork{
		Name:                      "Sepolia Testnet",
		ClientImplementation:      blockchain.EthereumClientImplementation,
		ChainID:                   11155111,
		Simulated:                 false,
		ChainlinkTransactionLimit: 5000,
		Timeout:                   blockchain.JSONStrDuration{time.Minute},
		MinimumConfirmations:      1,
		GasEstimationBuffer:       1000,
	}

	// goerliTestnet https://goerli.net/
	GoerliTestnet blockchain.EVMNetwork = blockchain.EVMNetwork{
		Name:                      "Goerli Testnet",
		ClientImplementation:      blockchain.EthereumClientImplementation,
		ChainID:                   5,
		Simulated:                 false,
		ChainlinkTransactionLimit: 5000,
		Timeout:                   blockchain.JSONStrDuration{5 * time.Minute},
		MinimumConfirmations:      1,
		GasEstimationBuffer:       1000,
	}

	// klaytnBaobab https://klaytn.foundation/
	KlaytnBaobab blockchain.EVMNetwork = blockchain.EVMNetwork{
		Name:                      "Klaytn Baobab",
		ClientImplementation:      blockchain.KlaytnClientImplementation,
		ChainID:                   1001,
		Simulated:                 false,
		ChainlinkTransactionLimit: 5000,
		Timeout:                   blockchain.JSONStrDuration{time.Minute},
		MinimumConfirmations:      1,
		GasEstimationBuffer:       0,
	}

	// metisStardust https://www.metis.io/
	MetisStardust blockchain.EVMNetwork = blockchain.EVMNetwork{
		Name:                      "Metis Stardust",
		ClientImplementation:      blockchain.MetisClientImplementation,
		ChainID:                   588,
		Simulated:                 false,
		ChainlinkTransactionLimit: 5000,
		Timeout:                   blockchain.JSONStrDuration{time.Minute},
		MinimumConfirmations:      1,
		GasEstimationBuffer:       1000,
	}

	// arbitrumGoerli https://developer.offchainlabs.com/docs/public_chains
	ArbitrumGoerli blockchain.EVMNetwork = blockchain.EVMNetwork{
		Name:                      "Arbitrum Goerli",
		ClientImplementation:      blockchain.ArbitrumClientImplementation,
		ChainID:                   421613,
		Simulated:                 false,
		ChainlinkTransactionLimit: 5000,
		Timeout:                   blockchain.JSONStrDuration{time.Minute},
		MinimumConfirmations:      0,
		GasEstimationBuffer:       0,
	}

	// optimismGoerli https://dev.optimism.io/kovan-to-goerli/
	OptimismGoerli blockchain.EVMNetwork = blockchain.EVMNetwork{
		Name:                      "Optimism Goerli",
		ClientImplementation:      blockchain.OptimismClientImplementation,
		ChainID:                   420,
		Simulated:                 false,
		ChainlinkTransactionLimit: 5000,
		Timeout:                   blockchain.JSONStrDuration{time.Minute},
		MinimumConfirmations:      1,
		GasEstimationBuffer:       0,
	}

	// rskTestnet https://www.rsk.co/
	RSKTestnet blockchain.EVMNetwork = blockchain.EVMNetwork{
		Name:                      "RSK Testnet",
		ClientImplementation:      blockchain.RSKClientImplementation,
		ChainID:                   31,
		Simulated:                 false,
		ChainlinkTransactionLimit: 5000,
		Timeout:                   blockchain.JSONStrDuration{time.Minute},
		MinimumConfirmations:      1,
		GasEstimationBuffer:       1000,
	}

	// PolygonMumbai https://mumbai.polygonscan.com/
	PolygonMumbai blockchain.EVMNetwork = blockchain.EVMNetwork{
		Name:                      "Polygon Mumbai",
		ClientImplementation:      blockchain.PolygonClientImplementation,
		ChainID:                   80001,
		Simulated:                 false,
		ChainlinkTransactionLimit: 5000,
		Timeout:                   blockchain.JSONStrDuration{time.Minute},
		MinimumConfirmations:      1,
		GasEstimationBuffer:       1000,
	}

	Avalanche = blockchain.EVMNetwork{
		Name:                      "Avalanche Fuji",
		ClientImplementation:      blockchain.EthereumClientImplementation,
		ChainID:                   43113,
		Simulated:                 false,
		ChainlinkTransactionLimit: 5000,
		Timeout:                   blockchain.JSONStrDuration{time.Minute},
		MinimumConfirmations:      1,
		GasEstimationBuffer:       1000,
	}

	mappedNetworks = map[string]blockchain.EVMNetwork{
		"SIMULATED":        SimulatedEVM,
		"SIMULATED_1":      SimulatedEVMNonDev1,
		"SIMULATED_2":      SimulatedEVMNonDev2,
		"SIMULATED_NONDEV": SimulatedEVMNonDev,
		// "GENERAL":         generalEVM, // See above
		"GOERLI":          GoerliTestnet,
		"SEPOLIA":         SepoliaTestnet,
		"KLAYTN_BAOBAB":   KlaytnBaobab,
		"METIS_STARDUST":  MetisStardust,
		"ARBITRUM_GOERLI": ArbitrumGoerli,
		"OPTIMISM_GOERLI": OptimismGoerli,
		"RSK":             RSKTestnet,
		"MUMBAI":          PolygonMumbai,
		"AVALANCHE_FUJI":  Avalanche,
	}
)

// determineSelectedNetworks uses `SELECTED_NETWORKS` to determine which network(s) to run the tests on
func determineSelectedNetworks() []blockchain.EVMNetwork {
	logging.Init()
	selectedNetworks := make([]blockchain.EVMNetwork, 0)
	rawSelectedNetworks := strings.ToUpper(os.Getenv("SELECTED_NETWORKS"))
	setNetworkNames := strings.Split(rawSelectedNetworks, ",")

	for _, setNetworkName := range setNetworkNames {
		if chosenNetwork, valid := mappedNetworks[setNetworkName]; valid {
			log.Info().
				Interface("SELECTED_NETWORKS", setNetworkNames).
				Str("Network Name", chosenNetwork.Name).
				Msg("Read network choice from 'SELECTED_NETWORKS'")
			setURLs(setNetworkName, &chosenNetwork)
			setKeys(setNetworkName, &chosenNetwork)
			selectedNetworks = append(selectedNetworks, chosenNetwork)
		} else {
			validNetworks := make([]string, 0)
			for validNetwork := range mappedNetworks {
				validNetworks = append(validNetworks, validNetwork)
			}
			log.Fatal().
				Interface("SELECTED_NETWORKS", setNetworkNames).
				Str("Valid Networks", strings.Join(validNetworks, ", ")).
				Msg("SELECTED_NETWORKS value is invalid. Use a valid network(s).")
		}
	}
	return selectedNetworks
}

// setURLs sets a network URL(s) based on env vars
func setURLs(prefix string, network *blockchain.EVMNetwork) {
	prefix = strings.Trim(prefix, "_")
	prefix = strings.ToUpper(prefix)

	if strings.Contains(prefix, "SIMULATED") { // Use defaults for SIMULATED
		return
	}

	wsEnvVar := fmt.Sprintf("%s_URLS", prefix)
	httpEnvVar := fmt.Sprintf("%s_HTTP_URLS", prefix)
	if os.Getenv(wsEnvVar) == "" {
		wsURLs := strings.Split(os.Getenv("EVM_URLS"), ",")
		httpURLs := strings.Split(os.Getenv("EVM_HTTP_URLS"), ",")
		log.Warn().
			Interface("EVM_URLS", wsURLs).
			Interface("EVM_HTTP_URLS", httpURLs).
			Msg(fmt.Sprintf("No '%s' env var defined, defaulting to 'EVM_URLS'", wsEnvVar))
		network.URLs = wsURLs
		network.HTTPURLs = httpURLs
		return
	}
	wsURLs := strings.Split(os.Getenv(wsEnvVar), ",")
	httpURLs := strings.Split(os.Getenv(httpEnvVar), ",")
	network.URLs = wsURLs
	network.HTTPURLs = httpURLs
	log.Info().Interface(wsEnvVar, wsURLs).Interface(httpEnvVar, httpURLs).Msg("Read network URLs")
}

// setKeys sets a network's private key(s) based on env vars
func setKeys(prefix string, network *blockchain.EVMNetwork) {
	prefix = strings.Trim(prefix, "_")
	prefix = strings.ToUpper(prefix)

	if strings.Contains(prefix, "SIMULATED") { // Use defaults for SIMULATED
		return
	}

	envVar := fmt.Sprintf("%s_KEYS", prefix)
	if os.Getenv(envVar) == "" {
		keys := strings.Split(os.Getenv("EVM_KEYS"), ",")
		log.Warn().
			Interface("EVM_KEYS", keys).
			Msg(fmt.Sprintf("No '%s' env var defined, defaulting to 'EVM_KEYS'", envVar))
		network.PrivateKeys = keys
		return
	}
	keys := strings.Split(os.Getenv(envVar), ",")
	network.PrivateKeys = keys
	log.Info().Interface(envVar, keys).Msg("Read network Keys")
}
