package networks

import (
	"fmt"
	"os"
	"strings"
	"time"

	"github.com/rs/zerolog"
	"github.com/rs/zerolog/log"

	"github.com/smartcontractkit/chainlink-testing-framework/blockchain"
)

// Pre-configured test networks and their connections
// Some networks with public RPC endpoints are already filled out, but make use of environment variables to use info like
// private RPC endpoints and private keys.
var (
	// SelectedNetwork uses the SELECTED_NETWORK env var to determine which network to run the test on
	SelectedNetwork *blockchain.EVMNetwork = determineSelectedNetwork()
	// SimulatedEVM represents a simulated network
	SimulatedEVM *blockchain.EVMNetwork = blockchain.SimulatedEVMNetwork
	// generalEVM is a customizable network through environment variables
	generalEVM *blockchain.EVMNetwork = blockchain.LoadNetworkFromEnvironment()

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
		"GENERAL":         generalEVM,
		"GOERLI":          GoerliTestnet,
		"SEPOLIA":         SepoliaTestnet,
		"KLAYTN_BAOBAB":   KlaytnBaobab,
		"METIS_STARDUST":  MetisStardust,
		"ARBITRUM_GOERLI": ArbitrumGoerli,
		"OPTIMISM_GOERLI": OptimismGoerli,
	}
)

// DetermineNetwork determines which network
func determineSelectedNetwork() *blockchain.EVMNetwork {
	log.Logger = log.Output(zerolog.ConsoleWriter{Out: os.Stderr})
	setNetwork := strings.ToUpper(os.Getenv("SELECTED_NETWORK"))
	if chosenNetwork, valid := mappedNetworks[setNetwork]; valid {
		log.Info().
			Str("SELECTED_NETWORK", setNetwork).
			Str("Network Name", chosenNetwork.Name).
			Msg("Read network choice from 'SELECTED_NETWORK'")
		chosenNetwork.URLs = getURLs(setNetwork)
		chosenNetwork.PrivateKeys = getKeys(setNetwork)
		return chosenNetwork
	}
	validNetworks := make([]string, 0)
	for validNetwork := range mappedNetworks {
		validNetworks = append(validNetworks, validNetwork)
	}
	log.Fatal().
		Str("SELECTED_NETWORK", setNetwork).
		Str("Valid Networks", strings.Join(validNetworks, ", ")).
		Msg("SELECTED_NETWORK value of is invalid. Use a listed valid one")
	return nil
}

func getURLs(prefix string) []string {
	log.Logger = log.Output(zerolog.ConsoleWriter{Out: os.Stderr})
	prefix = strings.Trim(prefix, "_")
	envVar := fmt.Sprintf("%s_URLS", prefix)
	if os.Getenv(envVar) == "" {
		urls := strings.Split(os.Getenv("EVM_URLS"), ",")
		log.Warn().
			Interface("EVM_URLS", urls).
			Msg(fmt.Sprintf("No '%s' env var defined, defaulting to 'EVM_URLS'", envVar))
		return urls
	}
	urls := strings.Split(os.Getenv(envVar), ",")
	log.Info().Interface(envVar, urls).Msg("Read network URLs")
	return urls
}

func getKeys(prefix string) []string {
	log.Logger = log.Output(zerolog.ConsoleWriter{Out: os.Stderr})
	strings.Trim(prefix, "_")
	envVar := fmt.Sprintf("%s_KEYS", prefix)
	if os.Getenv(envVar) == "" {
		keys := strings.Split(os.Getenv("EVM_PRIVATE_KEYS"), ",")
		log.Warn().Interface("EVM_PRIVATE_KEYS", keys).Msg(fmt.Sprintf("No '%s' env var defined, defaulting to 'EVM_PRIVATE_KEYS'", envVar))
		return keys
	}
	keys := strings.Split(os.Getenv(envVar), ",")
	log.Info().Interface(envVar, keys).Msg("Read network Keys")
	return keys
}
