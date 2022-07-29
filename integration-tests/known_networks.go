package networks

import (
	"os"
	"strings"
	"time"

	"github.com/smartcontractkit/chainlink-testing-framework/blockchain"

	"github.com/smartcontractkit/chainlink-env/logging"
)

func init() {
	logging.Init()
}

// Pre-configured test networks and their connections
// Some networks with public RPC endpoints are already filled out, but make use of environment variables to use info like
// private RPC endpoints and private keys.
var (
	// SimulatedEVM represents a simulated network
	SimulatedEVM *blockchain.EVMNetwork = blockchain.SimulatedEVMNetwork

	// SepoliaTestnet https://sepolia.dev/
	SepoliaTestnet *blockchain.EVMNetwork = &blockchain.EVMNetwork{
		Name:                      "Sepolia Testnet",
		ChainID:                   11155111,
		URLs:                      strings.Split(os.Getenv("EVM_URLS"), ","),
		Simulated:                 false,
		PrivateKeys:               strings.Split(os.Getenv("EVM_PRIVATE_KEYS"), ","),
		ChainlinkTransactionLimit: 5000,
		Timeout:                   time.Minute,
		MinimumConfirmations:      1,
		GasEstimationBuffer:       1000,
	}

	// GoerliTestnet https://goerli.net/
	GoerliTestnet *blockchain.EVMNetwork = &blockchain.EVMNetwork{
		Name:                      "Goerli Testnet",
		ChainID:                   5,
		URLs:                      strings.Split(os.Getenv("EVM_URLS"), ","),
		Simulated:                 false,
		PrivateKeys:               strings.Split(os.Getenv("EVM_PRIVATE_KEYS"), ","),
		ChainlinkTransactionLimit: 5000,
		Timeout:                   time.Minute,
		MinimumConfirmations:      1,
		GasEstimationBuffer:       1000,
	}

	// KlaytnBaobab https://klaytn.foundation/
	KlaytnBaobab *blockchain.EVMNetwork = &blockchain.EVMNetwork{
		Name:                      "Klaytn Baobab",
		ChainID:                   1001,
		URLs:                      strings.Split(os.Getenv("EVM_URLS"), ","),
		Simulated:                 false,
		PrivateKeys:               strings.Split(os.Getenv("EVM_PRIVATE_KEYS"), ","),
		ChainlinkTransactionLimit: 5000,
		Timeout:                   time.Minute,
		MinimumConfirmations:      1,
		GasEstimationBuffer:       0,
	}

	// MetisStardust https://www.metis.io/
	MetisStardust *blockchain.EVMNetwork = &blockchain.EVMNetwork{
		Name:                      "Metis Stardust",
		ChainID:                   588,
		URLs:                      []string{"wss://stardust-ws.metis.io/"},
		Simulated:                 false,
		PrivateKeys:               strings.Split(os.Getenv("EVM_PRIVATE_KEYS"), ","),
		ChainlinkTransactionLimit: 5000,
		Timeout:                   time.Minute,
		MinimumConfirmations:      1,
		GasEstimationBuffer:       1000,
	}

	ArbitrumGoerli *blockchain.EVMNetwork = &blockchain.EVMNetwork{
		Name:                      "Arbitrum Goerli",
		ChainID:                   421613,
		URLs:                      strings.Split(os.Getenv("EVM_URLS"), ","),
		Simulated:                 false,
		PrivateKeys:               strings.Split(os.Getenv("EVM_PRIVATE_KEYS"), ","),
		ChainlinkTransactionLimit: 5000,
		Timeout:                   time.Minute,
		MinimumConfirmations:      0,
		GasEstimationBuffer:       0,
	}

	ArbitrumRinkeby *blockchain.EVMNetwork = &blockchain.EVMNetwork{
		Name:                      "Arbitrum Rinkeby",
		ChainID:                   421611,
		URLs:                      strings.Split(os.Getenv("EVM_URLS"), ","),
		Simulated:                 false,
		PrivateKeys:               strings.Split(os.Getenv("EVM_PRIVATE_KEYS"), ","),
		ChainlinkTransactionLimit: 5000,
		Timeout:                   time.Minute,
		MinimumConfirmations:      0,
		GasEstimationBuffer:       0,
	}
)

// GeneralEVM loads general EVM settings from env vars
func GeneralEVM() *blockchain.EVMNetwork {
	return blockchain.LoadNetworkFromEnvironment()
}
