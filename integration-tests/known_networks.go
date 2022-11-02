package networks

import (
	"encoding/json"
	"fmt"
	"os"
	"strings"
	"time"

	"github.com/rs/zerolog/log"
	"gopkg.in/guregu/null.v4"

	"github.com/smartcontractkit/chainlink-testing-framework/blockchain"
	"github.com/smartcontractkit/chainlink-testing-framework/logging"

	"github.com/smartcontractkit/chainlink/core/chains/evm/types"
	"github.com/smartcontractkit/chainlink/core/utils"
)

// Pre-configured test networks and their connections
// Some networks with public RPC endpoints are already filled out, but make use of environment variables to use info like
// private RPC endpoints and private keys.
var (
	// SelectedNetworks uses the SELECTED_NETWORKS env var to determine which network to run the test on.
	// For use in tests that utilize multiple chains. For tests on one chain, see SelectedNetwork
	// For CCIP use index 1 and 2 of SELECTED_NETWORKS to denote source and destination network respectively
	SelectedNetworks []*blockchain.EVMNetwork = determineSelectedNetworks()
	// SelectedNetwork uses the first listed network in SELECTED_NETWORKS, for use in tests on only one chain
	SelectedNetwork *blockchain.EVMNetwork = SelectedNetworks[0]

	// SimulatedEVM represents a simulated network
	SimulatedEVM *blockchain.EVMNetwork = blockchain.SimulatedEVMNetwork
	// generalEVM is a customizable network through environment variables
	generalEVM *blockchain.EVMNetwork = blockchain.LoadNetworkFromEnvironment()

	NetworkAlpha = &blockchain.EVMNetwork{
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
		Timeout:                   2 * time.Minute,
		MinimumConfirmations:      1,
		GasEstimationBuffer:       10000,
	}

	NetworkBeta = &blockchain.EVMNetwork{
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
		Timeout:                   2 * time.Minute,
		MinimumConfirmations:      1,
		GasEstimationBuffer:       10000,
	}

	// sepoliaTestnet https://sepolia.dev/
	SepoliaTestnet *blockchain.EVMNetwork = &blockchain.EVMNetwork{
		Name:                      "Sepolia Testnet",
		ClientImplementation:      blockchain.EthereumClientImplementation,
		ChainID:                   11155111,
		Simulated:                 false,
		ChainlinkTransactionLimit: 5000,
		Timeout:                   time.Minute,
		MinimumConfirmations:      1,
		GasEstimationBuffer:       1000,
	}

	// goerliTestnet https://goerli.net/
	GoerliTestnet *blockchain.EVMNetwork = &blockchain.EVMNetwork{
		Name:                      "Goerli Testnet",
		ClientImplementation:      blockchain.EthereumClientImplementation,
		ChainID:                   5,
		Simulated:                 false,
		ChainlinkTransactionLimit: 5000,
		Timeout:                   time.Minute * 5,
		MinimumConfirmations:      1,
		GasEstimationBuffer:       1000,
	}

	// klaytnBaobab https://klaytn.foundation/
	KlaytnBaobab *blockchain.EVMNetwork = &blockchain.EVMNetwork{
		Name:                      "Klaytn Baobab",
		ClientImplementation:      blockchain.KlaytnClientImplementation,
		ChainID:                   1001,
		Simulated:                 false,
		ChainlinkTransactionLimit: 5000,
		Timeout:                   time.Minute,
		MinimumConfirmations:      1,
		GasEstimationBuffer:       0,
	}

	// metisStardust https://www.metis.io/
	MetisStardust *blockchain.EVMNetwork = &blockchain.EVMNetwork{
		Name:                      "Metis Stardust",
		ClientImplementation:      blockchain.MetisClientImplementation,
		ChainID:                   588,
		Simulated:                 false,
		ChainlinkTransactionLimit: 5000,
		Timeout:                   time.Minute,
		MinimumConfirmations:      1,
		GasEstimationBuffer:       1000,
	}

	// arbitrumGoerli https://developer.offchainlabs.com/docs/public_chains
	ArbitrumGoerli *blockchain.EVMNetwork = &blockchain.EVMNetwork{
		Name:                      "Arbitrum Goerli",
		ClientImplementation:      blockchain.ArbitrumClientImplementation,
		ChainID:                   421613,
		Simulated:                 false,
		ChainlinkTransactionLimit: 5000,
		Timeout:                   time.Minute,
		MinimumConfirmations:      0,
		GasEstimationBuffer:       0,
	}

	// optimismGoerli https://dev.optimism.io/kovan-to-goerli/
	OptimismGoerli *blockchain.EVMNetwork = &blockchain.EVMNetwork{
		Name:                      "Optimism Goerli",
		ClientImplementation:      blockchain.OptimismClientImplementation,
		ChainID:                   420,
		Simulated:                 false,
		ChainlinkTransactionLimit: 5000,
		Timeout:                   time.Minute,
		MinimumConfirmations:      0,
		GasEstimationBuffer:       0,
	}

	mappedNetworks = map[string]*blockchain.EVMNetwork{
		"SIMULATED":       SimulatedEVM,
		"SIMULATED_ALPHA": NetworkAlpha,
		"SIMULATED_BETA":  NetworkBeta,
		"GENERAL":         generalEVM,
		"GOERLI":          GoerliTestnet,
		"SEPOLIA":         SepoliaTestnet,
		"KLAYTN_BAOBAB":   KlaytnBaobab,
		"METIS_STARDUST":  MetisStardust,
		"ARBITRUM_GOERLI": ArbitrumGoerli,
		"OPTIMISM_GOERLI": OptimismGoerli,
	}
)

// determineSelectedNetworks uses `SELECTED_NETWORKS` to determine which network(s) to run the tests on
func determineSelectedNetworks() []*blockchain.EVMNetwork {
	logging.Init()
	selectedNetworks := make([]*blockchain.EVMNetwork, 0)
	setNetworkNames := strings.Split(strings.ToUpper(os.Getenv("SELECTED_NETWORKS")), ",")

	for _, setNetworkName := range setNetworkNames {
		if chosenNetwork, valid := mappedNetworks[setNetworkName]; valid {
			log.Info().
				Interface("SELECTED_NETWORKS", setNetworkNames).
				Str("Network Name", chosenNetwork.Name).
				Msg("Read network choice from 'SELECTED_NETWORKS'")
			setURLs(setNetworkName, chosenNetwork)
			setKeys(setNetworkName, chosenNetwork)
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

	if strings.Contains(prefix, "SIMULATED") { // Use defaults or read from env values for SIMULATED
		return
	}

	envVar := fmt.Sprintf("%s_URLS", prefix)
	if os.Getenv(envVar) == "" {
		urls := strings.Split(os.Getenv("EVM_URLS"), ",")
		log.Warn().
			Interface("EVM_URLS", urls).
			Msg(fmt.Sprintf("No '%s' env var defined, defaulting to 'EVM_URLS'", envVar))
		network.URLs = urls
		return
	}
	urls := strings.Split(os.Getenv(envVar), ",")
	network.URLs = urls
	log.Info().Interface(envVar, urls).Msg("Read network URLs")
}

// setKeys sets a network's private key(s) based on env vars
func setKeys(prefix string, network *blockchain.EVMNetwork) {
	prefix = strings.Trim(prefix, "_")
	prefix = strings.ToUpper(prefix)

	if strings.Contains(prefix, "SIMULATED") { // Use defaults or read from env values for SIMULATED
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

func DeriveEVMNodesFromNetworkSettings(networks ...blockchain.EVMNetwork) (string, error) {
	var evmNodes []types.NewNode
	for i, n := range networks {
		evmNodes = append(evmNodes, types.NewNode{
			Name:       fmt.Sprintf("network_%d", i),
			EVMChainID: *utils.NewBigI(n.ChainID),
			WSURL:      null.StringFrom(n.URLs[0]),
			HTTPURL:    null.StringFrom(n.HTTPURLs[0]),
			SendOnly:   false,
		})
	}
	if len(evmNodes) > 0 {
		evmNodes, err := json.Marshal(evmNodes)
		if err != nil {
			return "", err
		}
		return string(evmNodes), nil
	}
	return "", nil
}
