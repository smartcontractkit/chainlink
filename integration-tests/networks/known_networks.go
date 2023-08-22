// Package networks holds all known network information for the tests
package networks

import (
	"fmt"
	"os"
	"strings"
	"time"

	"github.com/rs/zerolog/log"

	"github.com/smartcontractkit/chainlink-testing-framework/utils"

	"github.com/smartcontractkit/chainlink-testing-framework/blockchain"
	"github.com/smartcontractkit/chainlink-testing-framework/logging"
)

// Pre-configured test networks and their connections
// Some networks with public RPC endpoints are already filled out, but make use of environment variables to use info like
// private RPC endpoints and private keys.
var (
	// To create replica of simulated EVM network, with different chain ids
	AdditionalSimulatedChainIds = []int64{3337, 4337, 5337, 6337, 7337, 8337, 9337, 9338}
	AdditionalSimulatedPvtKeys  = []string{
		"5de4111afa1a4b94908f83103eb1f1706367c2e68ca870fc3fb9a804cdab365a",
		"7c852118294e51e653712a81e05800f419141751be58f605c371e15141b007a6",
		"47e179ec197488593b187f80a00eb0da91f1b9d0b13f8733639f19c30a34926a",
		"8b3a350cf5c34c9194ca85829a2df0ec3153be0318b5e2d3348e872092edffba",
		"92db14e403b83dfe3df233f83dfa3a0d7096f21ca9b0d6d6b8d88b2b4ec1564e",
		"4bbbf85ce3377467afe5d46f804f221813b2bb87f24d81f60f1fcdbf7cbf4356",
		"dbda1821b80551c9d65939329250298aa3472ba22feea921c0cf5d620ea67b97",
		"2a871d0798f97d79848a013d4936a73bf4cc922c825d33c1cf7073dff6d409c6",
	}
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
		SupportsEIP1559:      true,
		ChainID:              1337,
		PrivateKeys: []string{
			"ac0974bec39a17e36ba4a6b4d238ff944bacb478cbed5efcae784d7bf4f2ff80",
		},
		URLs:                      []string{"ws://source-chain-ethereum-geth:8546"},
		HTTPURLs:                  []string{"http://source-chain-ethereum-geth:8544"},
		ChainlinkTransactionLimit: 500000,
		Timeout:                   blockchain.JSONStrDuration{Duration: 2 * time.Minute},
		MinimumConfirmations:      1,
		GasEstimationBuffer:       10000,
		DefaultGasLimit:           6000000,
	}

	// SimulatedEVM_NON_DEV_2 represents a simulated network with chain id 2337 which can be used to deploy a non-dev geth node
	SimulatedEVMNonDev2 = blockchain.EVMNetwork{
		Name:                 "dest-chain",
		Simulated:            true,
		SupportsEIP1559:      true,
		ClientImplementation: blockchain.EthereumClientImplementation,
		ChainID:              2337,
		PrivateKeys: []string{
			"59c6995e998f97a5a0044966f0945389dc9e86dae88c7a8412f4603b6b78690d",
		},
		URLs:                      []string{"ws://dest-chain-ethereum-geth:8546"},
		HTTPURLs:                  []string{"http://dest-chain-ethereum-geth:8544"},
		ChainlinkTransactionLimit: 500000,
		Timeout:                   blockchain.JSONStrDuration{Duration: 2 * time.Minute},
		MinimumConfirmations:      1,
		GasEstimationBuffer:       10000,
		DefaultGasLimit:           6000000,
	}

	SimulatedEVMNonDev = blockchain.EVMNetwork{
		Name:                 "geth",
		Simulated:            true,
		SupportsEIP1559:      true,
		ClientImplementation: blockchain.EthereumClientImplementation,
		ChainID:              1337,
		PrivateKeys: []string{
			"ac0974bec39a17e36ba4a6b4d238ff944bacb478cbed5efcae784d7bf4f2ff80",
		},
		URLs:                      []string{"ws://geth-ethereum-geth:8546"},
		HTTPURLs:                  []string{"http://geth-ethereum-geth:8544"},
		ChainlinkTransactionLimit: 500000,
		Timeout:                   blockchain.JSONStrDuration{Duration: 2 * time.Minute},
		MinimumConfirmations:      1,
		GasEstimationBuffer:       10000,
	}

	EthereumMainnet blockchain.EVMNetwork = blockchain.EVMNetwork{
		Name:                      "Ethereum Mainnet",
		SupportsEIP1559:           true,
		ClientImplementation:      blockchain.EthereumClientImplementation,
		ChainID:                   1,
		Simulated:                 false,
		ChainlinkTransactionLimit: 5000,
		Timeout:                   blockchain.JSONStrDuration{Duration: 5 * time.Minute},
		MinimumConfirmations:      1,
		GasEstimationBuffer:       0,
		FinalityTag:               true,
		DefaultGasLimit:           6000000,
	}

	// sepoliaTestnet https://sepolia.dev/
	SepoliaTestnet blockchain.EVMNetwork = blockchain.EVMNetwork{
		Name:                      "Sepolia Testnet",
		SupportsEIP1559:           true,
		ClientImplementation:      blockchain.EthereumClientImplementation,
		ChainID:                   11155111,
		Simulated:                 false,
		ChainlinkTransactionLimit: 5000,
		Timeout:                   blockchain.JSONStrDuration{Duration: time.Minute},
		MinimumConfirmations:      1,
		GasEstimationBuffer:       1000,
		FinalityTag:               true,
		DefaultGasLimit:           6000000,
	}

	// goerliTestnet https://goerli.net/
	GoerliTestnet blockchain.EVMNetwork = blockchain.EVMNetwork{
		Name:                      "Goerli Testnet",
		SupportsEIP1559:           true,
		ClientImplementation:      blockchain.EthereumClientImplementation,
		ChainID:                   5,
		Simulated:                 false,
		ChainlinkTransactionLimit: 5000,
		Timeout:                   blockchain.JSONStrDuration{Duration: 5 * time.Minute},
		MinimumConfirmations:      1,
		GasEstimationBuffer:       1000,
		FinalityTag:               true,
		DefaultGasLimit:           6000000,
	}

	KlaytnMainnet blockchain.EVMNetwork = blockchain.EVMNetwork{
		Name:                      "Klaytn Mainnet",
		SupportsEIP1559:           false,
		ClientImplementation:      blockchain.KlaytnClientImplementation,
		ChainID:                   8217,
		Simulated:                 false,
		ChainlinkTransactionLimit: 5000,
		Timeout:                   blockchain.JSONStrDuration{Duration: time.Minute},
		MinimumConfirmations:      1,
		GasEstimationBuffer:       0,
	}

	// klaytnBaobab https://klaytn.foundation/
	KlaytnBaobab blockchain.EVMNetwork = blockchain.EVMNetwork{
		Name:                      "Klaytn Baobab",
		SupportsEIP1559:           false,
		ClientImplementation:      blockchain.KlaytnClientImplementation,
		ChainID:                   1001,
		Simulated:                 false,
		ChainlinkTransactionLimit: 5000,
		Timeout:                   blockchain.JSONStrDuration{Duration: time.Minute},
		MinimumConfirmations:      1,
		GasEstimationBuffer:       0,
	}

	MetisAndromeda blockchain.EVMNetwork = blockchain.EVMNetwork{
		Name:                      "Metis Andromeda",
		SupportsEIP1559:           false,
		ClientImplementation:      blockchain.MetisClientImplementation,
		ChainID:                   1088,
		Simulated:                 false,
		ChainlinkTransactionLimit: 5000,
		Timeout:                   blockchain.JSONStrDuration{Duration: time.Minute},
		MinimumConfirmations:      1,
		GasEstimationBuffer:       0,
	}

	// metisStardust https://www.metis.io/
	MetisStardust blockchain.EVMNetwork = blockchain.EVMNetwork{
		Name:                      "Metis Stardust",
		SupportsEIP1559:           false,
		ClientImplementation:      blockchain.MetisClientImplementation,
		ChainID:                   588,
		Simulated:                 false,
		ChainlinkTransactionLimit: 5000,
		Timeout:                   blockchain.JSONStrDuration{Duration: time.Minute},
		MinimumConfirmations:      1,
		GasEstimationBuffer:       1000,
	}

	ArbitrumMainnet blockchain.EVMNetwork = blockchain.EVMNetwork{
		Name:                      "Arbitrum Mainnet",
		SupportsEIP1559:           true,
		ClientImplementation:      blockchain.ArbitrumClientImplementation,
		ChainID:                   42161,
		Simulated:                 false,
		ChainlinkTransactionLimit: 5000,
		Timeout:                   blockchain.JSONStrDuration{Duration: 2 * time.Minute},
		MinimumConfirmations:      0,
		GasEstimationBuffer:       0,
		FinalityTag:               true,
		DefaultGasLimit:           100000000,
	}

	// arbitrumGoerli https://developer.offchainlabs.com/docs/public_chains
	ArbitrumGoerli blockchain.EVMNetwork = blockchain.EVMNetwork{
		Name:                      "Arbitrum Goerli",
		SupportsEIP1559:           true,
		ClientImplementation:      blockchain.ArbitrumClientImplementation,
		ChainID:                   421613,
		Simulated:                 false,
		ChainlinkTransactionLimit: 5000,
		Timeout:                   blockchain.JSONStrDuration{Duration: time.Minute},
		MinimumConfirmations:      0,
		GasEstimationBuffer:       0,
		FinalityTag:               true,
		DefaultGasLimit:           100000000,
	}

	OptimismMainnet blockchain.EVMNetwork = blockchain.EVMNetwork{
		Name:                      "Optimism Mainnet",
		SupportsEIP1559:           true,
		ClientImplementation:      blockchain.OptimismClientImplementation,
		ChainID:                   10,
		Simulated:                 false,
		ChainlinkTransactionLimit: 5000,
		Timeout:                   blockchain.JSONStrDuration{Duration: time.Minute},
		MinimumConfirmations:      1,
		GasEstimationBuffer:       0,
		FinalityTag:               true,
		DefaultGasLimit:           6000000,
	}

	// optimismGoerli https://dev.optimism.io/kovan-to-goerli/
	OptimismGoerli blockchain.EVMNetwork = blockchain.EVMNetwork{
		Name:                      "Optimism Goerli",
		SupportsEIP1559:           true,
		ClientImplementation:      blockchain.OptimismClientImplementation,
		ChainID:                   420,
		Simulated:                 false,
		ChainlinkTransactionLimit: 5000,
		Timeout:                   blockchain.JSONStrDuration{Duration: time.Minute},
		MinimumConfirmations:      1,
		GasEstimationBuffer:       0,
		FinalityTag:               true,
		DefaultGasLimit:           6000000,
	}

	RSKMainnet blockchain.EVMNetwork = blockchain.EVMNetwork{
		Name:                      "RSK Mainnet",
		SupportsEIP1559:           false,
		ClientImplementation:      blockchain.RSKClientImplementation,
		ChainID:                   30,
		Simulated:                 false,
		ChainlinkTransactionLimit: 5000,
		Timeout:                   blockchain.JSONStrDuration{Duration: time.Minute},
		MinimumConfirmations:      1,
		GasEstimationBuffer:       1000,
	}

	// rskTestnet https://www.rsk.co/
	RSKTestnet blockchain.EVMNetwork = blockchain.EVMNetwork{
		Name:                      "RSK Testnet",
		SupportsEIP1559:           false,
		ClientImplementation:      blockchain.RSKClientImplementation,
		ChainID:                   31,
		Simulated:                 false,
		ChainlinkTransactionLimit: 5000,
		Timeout:                   blockchain.JSONStrDuration{Duration: time.Minute},
		MinimumConfirmations:      1,
		GasEstimationBuffer:       1000,
	}

	PolygonMainnet blockchain.EVMNetwork = blockchain.EVMNetwork{
		Name:                      "Polygon Mainnet",
		SupportsEIP1559:           true,
		ClientImplementation:      blockchain.PolygonClientImplementation,
		ChainID:                   137,
		Simulated:                 false,
		ChainlinkTransactionLimit: 5000,
		Timeout:                   blockchain.JSONStrDuration{Duration: 2 * time.Minute},
		MinimumConfirmations:      1,
		GasEstimationBuffer:       0,
		FinalityDepth:             550,
		DefaultGasLimit:           6000000,
	}

	// PolygonMumbai https://mumbai.polygonscan.com/
	PolygonMumbai blockchain.EVMNetwork = blockchain.EVMNetwork{
		Name:                      "Polygon Mumbai",
		SupportsEIP1559:           true,
		ClientImplementation:      blockchain.PolygonClientImplementation,
		ChainID:                   80001,
		Simulated:                 false,
		ChainlinkTransactionLimit: 5000,
		Timeout:                   blockchain.JSONStrDuration{Duration: time.Minute},
		MinimumConfirmations:      1,
		GasEstimationBuffer:       1000,
		FinalityDepth:             550,
		DefaultGasLimit:           6000000,
	}

	AvalancheMainnet blockchain.EVMNetwork = blockchain.EVMNetwork{
		Name:                      "Avalanche Mainnet",
		SupportsEIP1559:           true,
		ClientImplementation:      blockchain.EthereumClientImplementation,
		ChainID:                   43114,
		Simulated:                 false,
		ChainlinkTransactionLimit: 5000,
		Timeout:                   blockchain.JSONStrDuration{Duration: time.Minute},
		MinimumConfirmations:      1,
		GasEstimationBuffer:       0,
		FinalityDepth:             35,
		DefaultGasLimit:           6000000,
	}

	AvalancheFuji = blockchain.EVMNetwork{
		Name:                      "Avalanche Fuji",
		SupportsEIP1559:           true,
		ClientImplementation:      blockchain.EthereumClientImplementation,
		ChainID:                   43113,
		Simulated:                 false,
		ChainlinkTransactionLimit: 5000,
		Timeout:                   blockchain.JSONStrDuration{Duration: time.Minute},
		MinimumConfirmations:      1,
		GasEstimationBuffer:       1000,
		FinalityDepth:             35,
		DefaultGasLimit:           6000000,
	}

	Quorum = blockchain.EVMNetwork{
		Name:                      "Quorum",
		SupportsEIP1559:           false,
		ClientImplementation:      blockchain.QuorumClientImplementation,
		ChainID:                   1337,
		Simulated:                 false,
		ChainlinkTransactionLimit: 5000,
		Timeout:                   blockchain.JSONStrDuration{Duration: time.Minute},
		MinimumConfirmations:      1,
		GasEstimationBuffer:       0,
	}

	BaseGoerli blockchain.EVMNetwork = blockchain.EVMNetwork{
		Name:                      "Base Goerli",
		SupportsEIP1559:           true,
		ClientImplementation:      blockchain.OptimismClientImplementation,
		ChainID:                   84531,
		Simulated:                 false,
		ChainlinkTransactionLimit: 5000,
		Timeout:                   blockchain.JSONStrDuration{Duration: time.Minute},
		MinimumConfirmations:      0,
		GasEstimationBuffer:       0,
	}

	CeloAlfajores = blockchain.EVMNetwork{
		Name:                      "Celo Alfajores",
		SupportsEIP1559:           false,
		ClientImplementation:      blockchain.CeloClientImplementation,
		ChainID:                   44787,
		Simulated:                 false,
		ChainlinkTransactionLimit: 5000,
		Timeout:                   blockchain.JSONStrDuration{Duration: time.Minute},
		MinimumConfirmations:      1,
		GasEstimationBuffer:       1000,
	}

	ScrollSepolia = blockchain.EVMNetwork{
		Name:                      "Scroll Sepolia",
		ClientImplementation:      blockchain.ScrollClientImplementation,
		ChainID:                   534351,
		Simulated:                 false,
		ChainlinkTransactionLimit: 5000,
		Timeout:                   blockchain.JSONStrDuration{Duration: time.Minute},
		MinimumConfirmations:      1,
		GasEstimationBuffer:       0,
	}

	ScrollMainnet = blockchain.EVMNetwork{
		Name:                      "Scroll Mainnet",
		ClientImplementation:      blockchain.ScrollClientImplementation,
		ChainID:                   534352,
		Simulated:                 false,
		ChainlinkTransactionLimit: 5000,
		Timeout:                   blockchain.JSONStrDuration{Duration: time.Minute},
		MinimumConfirmations:      1,
		GasEstimationBuffer:       0,
	}

	CeloMainnet = blockchain.EVMNetwork{
		Name:                      "Celo",
		ClientImplementation:      blockchain.CeloClientImplementation,
		ChainID:                   42220,
		Simulated:                 false,
		ChainlinkTransactionLimit: 5000,
		Timeout:                   blockchain.JSONStrDuration{Duration: time.Minute},
		MinimumConfirmations:      1,
		GasEstimationBuffer:       1000,
	}

	BaseMainnet blockchain.EVMNetwork = blockchain.EVMNetwork{
		Name:                      "Base Mainnet",
		SupportsEIP1559:           true,
		ClientImplementation:      blockchain.OptimismClientImplementation,
		ChainID:                   8453,
		Simulated:                 false,
		ChainlinkTransactionLimit: 5000,
		Timeout:                   blockchain.JSONStrDuration{Duration: time.Minute},
		MinimumConfirmations:      0,
		GasEstimationBuffer:       0,
	}

	BSCTestnet blockchain.EVMNetwork = blockchain.EVMNetwork{
		Name:                      "BSC Testnet",
		SupportsEIP1559:           true,
		ClientImplementation:      blockchain.BSCClientImplementation,
		ChainID:                   97,
		Simulated:                 false,
		ChainlinkTransactionLimit: 5000,
		Timeout:                   blockchain.JSONStrDuration{Duration: time.Minute},
		MinimumConfirmations:      3,
		GasEstimationBuffer:       0,
	}

	BSCMainnet blockchain.EVMNetwork = blockchain.EVMNetwork{
		Name:                      "BSC Mainnet",
		SupportsEIP1559:           true,
		ClientImplementation:      blockchain.BSCClientImplementation,
		ChainID:                   56,
		Simulated:                 false,
		ChainlinkTransactionLimit: 5000,
		Timeout:                   blockchain.JSONStrDuration{Duration: time.Minute},
		MinimumConfirmations:      3,
		GasEstimationBuffer:       0,
	}

	MappedNetworks = map[string]blockchain.EVMNetwork{
		"SIMULATED":        SimulatedEVM,
		"SIMULATED_1":      SimulatedEVMNonDev1,
		"SIMULATED_2":      SimulatedEVMNonDev2,
		"SIMULATED_NONDEV": SimulatedEVMNonDev,
		// "GENERAL":         generalEVM, // See above
		"ETHEREUM_MAINNET":  EthereumMainnet,
		"GOERLI":            GoerliTestnet,
		"SEPOLIA":           SepoliaTestnet,
		"KLAYTN_MAINNET":    KlaytnMainnet,
		"KLAYTN_BAOBAB":     KlaytnBaobab,
		"METIS_ANDROMEDA":   MetisAndromeda,
		"METIS_STARDUST":    MetisStardust,
		"ARBITRUM_MAINNET":  ArbitrumMainnet,
		"ARBITRUM_GOERLI":   ArbitrumGoerli,
		"OPTIMISM_MAINNET":  OptimismMainnet,
		"OPTIMISM_GOERLI":   OptimismGoerli,
		"BASE_GOERLI":       BaseGoerli,
		"CELO_ALFAJORES":    CeloAlfajores,
		"CELO_MAINNET":      CeloMainnet,
		"RSK":               RSKTestnet,
		"MUMBAI":            PolygonMumbai,
		"POLYGON_MAINNET":   PolygonMainnet,
		"AVALANCHE_FUJI":    AvalancheFuji,
		"AVALANCHE_MAINNET": AvalancheMainnet,
		"QUORUM":            Quorum,
		"SCROLL_SEPOLIA":    ScrollSepolia,
		"SCROLL_MAINNET":    ScrollMainnet,
		"BASE_MAINNET":      BaseMainnet,
		"BSC_TESTNET":       BSCTestnet,
		"BSC_MAINNET":       BSCMainnet,
	}
)

// determineSelectedNetworks uses `SELECTED_NETWORKS` to determine which networks to run the tests on.
// Use DetermineSelectedNetwork for tests that only use one network
func determineSelectedNetworks() []blockchain.EVMNetwork {
	logging.Init()
	selectedNetworks := make([]blockchain.EVMNetwork, 0)
	rawSelectedNetworks := strings.ToUpper(os.Getenv("SELECTED_NETWORKS"))
	setNetworkNames := strings.Split(rawSelectedNetworks, ",")

	for _, setNetworkName := range setNetworkNames {
		if chosenNetwork, valid := MappedNetworks[setNetworkName]; valid {
			log.Info().
				Interface("SELECTED_NETWORKS", setNetworkNames).
				Str("Network Name", chosenNetwork.Name).
				Msg("Read network choice from 'SELECTED_NETWORKS'")
			setURLs(setNetworkName, &chosenNetwork)
			setKeys(setNetworkName, &chosenNetwork)
			selectedNetworks = append(selectedNetworks, chosenNetwork)
		} else {
			validNetworks := make([]string, 0)
			for validNetwork := range MappedNetworks {
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
	wsEnvURLs, err := utils.GetEnv(wsEnvVar)
	if err != nil {
		log.Fatal().Err(err).Str("env var", wsEnvVar).Msg("Error getting env var")
	}
	httpEnvURLs, err := utils.GetEnv(httpEnvVar)
	if err != nil {
		log.Fatal().Err(err).Str("env var", httpEnvVar).Msg("Error getting env var")
	}
	if wsEnvURLs == "" {
		evmUrls, err := utils.GetEnv("EVM_URLS")
		if err != nil {
			log.Fatal().Err(err).Str("env var", "EVM_URLS").Msg("Error getting env var")
		}
		evmhttpUrls, err := utils.GetEnv("EVM_HTTP_URLS")
		if err != nil {
			log.Fatal().Err(err).Str("env var", "EVM_HTTP_URLS").Msg("Error getting env var")
		}
		wsURLs := strings.Split(evmUrls, ",")
		httpURLs := strings.Split(evmhttpUrls, ",")
		log.Warn().
			Interface("EVM_URLS", wsURLs).
			Interface("EVM_HTTP_URLS", httpURLs).
			Msgf("No '%s' env var defined, defaulting to 'EVM_URLS'", wsEnvVar)
		network.URLs = wsURLs
		network.HTTPURLs = httpURLs
		return
	}

	wsURLs := strings.Split(wsEnvURLs, ",")
	httpURLs := strings.Split(httpEnvURLs, ",")
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
	keysEnv, err := utils.GetEnv(envVar)
	if err != nil {
		log.Fatal().Err(err).Str("env var", envVar).Msg("Error getting env var")
	}
	if keysEnv == "" {
		keys := strings.Split(os.Getenv("EVM_KEYS"), ",")
		log.Warn().
			Interface("EVM_KEYS", keys).
			Msg(fmt.Sprintf("No '%s' env var defined, defaulting to 'EVM_KEYS'", envVar))
		network.PrivateKeys = keys
		return
	}
	keys := strings.Split(keysEnv, ",")
	network.PrivateKeys = keys
	log.Info().Interface(envVar, keys).Msg("Read network Keys")
}
