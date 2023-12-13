package soak

import (
	"testing"

	"github.com/stretchr/testify/require"

	"github.com/smartcontractkit/chainlink-testing-framework/logging"

	"github.com/smartcontractkit/chainlink/integration-tests/actions"
	"github.com/smartcontractkit/chainlink/integration-tests/testsetups"
)

func TestOCRSoak(t *testing.T) {
	l := logging.GetTestLogger(t)
	// Use this variable to pass in any custom EVM specific TOML values to your Chainlink nodes
	customNetworkTOML := `ChainID = '1029'
	ChainType = 'bttc'
	FinalityDepth = 500 
	# blocks are generated every 2-4s 
	LogPollInterval = '2s
	
	[GasEstimator] 
	Mode = 'BlockHistory'
	EIP1559DynamicFees = false
	PriceMax = '115792089237316195423570985008687907853269984665.640564039457584007913129639935 tether'
	PriceMin = '0.009 ether'
	
	[GasEstimator.BlockHistory]
	# how many blocks we want to keep in memory to calculate gas price
	# # Average block time of 2s
	BlockHistorySize = 24
	
	[Transactions]
	ResendAfterThreshold = '30s'
	
	[HeadTracker]
	# re-org for bttc is really high so we want to check for the block where reorg happens
	HistoryDepth = 500
	
	[NodePool]
	SyncThreshold = 10
	
	[OCR]
	ContractConfirmations = 1`
	// Uncomment below for debugging TOML issues on the node
	// network := networks.MustGetSelectedNetworksFromEnv()[0]
	// fmt.Println("Using Chainlink TOML\n---------------------")
	// fmt.Println(networks.AddNetworkDetailedConfig(config.BaseOCR1Config, customNetworkTOML, network))
	// fmt.Println("---------------------")

	ocrSoakTest, err := testsetups.NewOCRSoakTest(t, false)
	require.NoError(t, err, "Error creating soak test")
	if !ocrSoakTest.Interrupted() {
		ocrSoakTest.DeployEnvironment(customNetworkTOML)
	}
	if ocrSoakTest.Environment().WillUseRemoteRunner() {
		return
	}
	t.Cleanup(func() {
		if err := actions.TeardownRemoteSuite(ocrSoakTest.TearDownVals(t)); err != nil {
			l.Error().Err(err).Msg("Error tearing down environment")
		}
	})
	if ocrSoakTest.Interrupted() {
		err = ocrSoakTest.LoadState()
		require.NoError(t, err, "Error loading state")
		ocrSoakTest.Resume()
	} else {
		ocrSoakTest.Setup()
		ocrSoakTest.Run()
	}
}
