package networks

import (
	"os"
	"strings"
	"time"

	"github.com/smartcontractkit/chainlink-testing-framework/blockchain"
)

// Pre-configured test networks and their connections
// Some networks with public RPC endpoints are already filled out, but make use of environment variables to use info like
// private RPC endpoints and private keys.
var (
	// SimulatedEVM represents a simulated network
	SimulatedEVM = blockchain.SimulatedEVMNetwork

	// MetisStardust holds default values for the Metis Stardust testnet
	MetisStardust *blockchain.EVMNetwork = &blockchain.EVMNetwork{
		Name:                      "Metis Stardust Network",
		ChainID:                   588,
		URLs:                      []string{"wss://stardust-ws.metis.io/"},
		Simulated:                 false,
		PrivateKeys:               strings.Split(os.Getenv("EVM_PRIVATE_KEYS"), ","),
		ChainlinkTransactionLimit: 5000,
		Timeout:                   time.Minute,
		MinimumConfirmations:      1,
		GasEstimationBuffer:       0,
	}
)
