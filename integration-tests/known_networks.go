package networks

import (
	"os"
	"path/filepath"
	"strings"
	"time"

	"github.com/joho/godotenv"
	"github.com/rs/zerolog/log"
	"github.com/smartcontractkit/chainlink-testing-framework/blockchain"
)

// Pre-configured test networks and their connections
// Some networks with public RPC endpoints are already filled out, but make use of environment variables to use info like
// private RPC endpoints and private keys.
var (
	// SimulatedEVMNetwork represents a simulated network
	SimulatedEVMNetwork = blockchain.SimulatedEVMNetwork

	// MetisTestNetwork holds default values for the Metis Stardust testnet
	MetisTestNetwork *blockchain.EVMNetwork
)

// LoadNetworks utilizes a .env file to load all env vars and assign values to preset networks
func LoadNetworks(dotEnvPath string) {
	absPath, err := filepath.Abs(dotEnvPath)
	if err != nil {
		log.Error().Err(err).Msg("Error loading .env file, proceeding with default values")
	} else {
		err = godotenv.Load(absPath)
		if err != nil {
			log.Error().Str("Path", absPath).Err(err).Msg("Error loading .env file, proceeding with default values")
		}
	}

	MetisTestNetwork = &blockchain.EVMNetwork{
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

	log.Info().Msg("Loaded Networks")
}
