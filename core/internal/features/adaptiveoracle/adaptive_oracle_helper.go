// Package adaptiveoracle contains a PoC smoke test validating that the v0.3-adaptive-oracle
// contract design (from https://github.com/smartcontractkit/svr-auction-don) plays nice with the
// existing off-chain OCR2 median plugin: a real running DON, driving real signed/transmitted
// reports into a DualAggregator, should end up flowing through AdaptiveOracle's transformAnswer
// hook exactly as it does in the forge test suite in that repo.
//
// This mirrors the pattern already used for core/internal/features/svr's DualAggregator test: an
// in-process cltest harness against a simulated backend, rather than the heavier devenv/products
// Docker-based framework used for the standard OCR2Aggregator smoke test. That's an intentional
// scope choice for this PoC -- the question being answered here is "does this on-chain design
// integrate with the off-chain logic", not "is this production deployment tooling".
package adaptiveoracle

import (
	"fmt"
	"math/big"
	"testing"

	"github.com/ethereum/go-ethereum/accounts/abi/bind"
	"github.com/ethereum/go-ethereum/common"
	"github.com/ethereum/go-ethereum/core/types"
	"github.com/ethereum/go-ethereum/eth/ethconfig"
	"github.com/ethereum/go-ethereum/ethclient/simulated"
	gethnode "github.com/ethereum/go-ethereum/node"
	"github.com/stretchr/testify/require"

	"github.com/smartcontractkit/freeport"
	testoffchainaggregator2 "github.com/smartcontractkit/libocr/gethwrappers2/testocr2aggregator"
	confighelper2 "github.com/smartcontractkit/libocr/offchainreporting2plus/confighelper"

	commonconfig "github.com/smartcontractkit/chainlink-common/pkg/config"
	"github.com/smartcontractkit/chainlink-evm/gethwrappers/generated/link_token_interface"
	"github.com/smartcontractkit/chainlink-evm/pkg/config/toml"
	evmtestutils "github.com/smartcontractkit/chainlink-evm/pkg/testutils"

	generated "github.com/smartcontractkit/chainlink/v2/core/internal/features/adaptiveoracle/generated"
	"github.com/smartcontractkit/chainlink/v2/core/internal/features/ocr2"
	"github.com/smartcontractkit/chainlink/v2/core/logger"
	"github.com/smartcontractkit/chainlink/v2/core/services/ocr2/testhelpers"
	"github.com/smartcontractkit/chainlink/v2/core/services/ocrbootstrap"
)

// AdaptiveOracleDecimals matches across DualAggregator, ReferenceRateAdapterMock, and AdaptiveOracle
// for this test -- no underlying rate feed is configured, so AdaptiveOracle's output decimals equal
// the market/reference decimals (see AdaptiveOracle::setUnderlyingRateFeed in the contract).
const AdaptiveOracleDecimals = 9

// Contracts bundles the deployed contracts and their addresses for the adaptive oracle stack.
type Contracts struct {
	DualAggregatorAddress    common.Address
	DualAggregator           *generated.DualAggregator
	AdaptiveOracleAddress    common.Address
	AdaptiveOracle           *generated.AdaptiveOracle
	AdaptiveRateLogicAddress common.Address
	ReferenceRateAdapter     *generated.ReferenceRateAdapterMock
}

// SetupAdaptiveOracleContracts deploys the full v0.3-adaptive-oracle contract stack (DualAggregator,
// AdaptiveOracle, AdaptiveRateLogic, and a ReferenceRateAdapterMock standing in for a real reference
// rate source) on a simulated backend, and wires them together exactly as the forge test suite does.
func SetupAdaptiveOracleContracts(t *testing.T, referenceRate int64) (*bind.TransactOpts, *simulated.Backend, *toml.Node, *Contracts) {
	owner := evmtestutils.MustNewSimTransactor(t)
	startingBalance := new(big.Int)
	startingBalance, _ = startingBalance.SetString("100000000000000000000", 10) // 100 eth
	genesisData := types.GenesisAlloc{owner.From: {Balance: startingBalance}}
	gasLimit := ethconfig.Defaults.Miner.GasCeil * 2

	httpPort := freeport.GetOne(t)
	wsPort := freeport.GetOne(t)
	host := "localhost"
	nodeConfig := toml.Node{
		Name:              new("simulated-node"),
		WSURL:             new(commonconfig.URL{Scheme: "ws", Host: fmt.Sprintf("%s:%d", host, wsPort)}),
		HTTPURL:           new(commonconfig.URL{Scheme: "http", Host: fmt.Sprintf("%s:%d", host, httpPort)}),
		Order:             new(int32(1)),
		IsLoadBalancedRPC: new(false),
	}

	b := simulated.NewBackend(genesisData, simulated.WithBlockGasLimit(gasLimit), withRPCServer(host, httpPort, wsPort, []string{"eth", "net", "web3"}))

	linkTokenAddress, _, linkContract, err := link_token_interface.DeployLinkToken(owner, b.Client())
	require.NoError(t, err)
	b.Commit()

	accessAddress, _, _, err := testoffchainaggregator2.DeploySimpleWriteAccessController(owner, b.Client())
	require.NoError(t, err, "failed to deploy test access controller contract")
	b.Commit()

	minAnswer, maxAnswer := new(big.Int), new(big.Int)
	minAnswer.Exp(big.NewInt(-2), big.NewInt(191), nil)
	maxAnswer.Exp(big.NewInt(2), big.NewInt(191), nil)
	maxAnswer.Sub(maxAnswer, big.NewInt(1))

	dualAggregatorAddress, _, dualAggregator, err := generated.DeployDualAggregator(
		owner,
		b.Client(),
		linkTokenAddress,
		minAnswer,
		maxAnswer,
		accessAddress,
		accessAddress,
		AdaptiveOracleDecimals,
		"ADAPTIVE-MARKET",
		common.Address{}, // secondaryProxy_: unused by this smoke test
		0,                // cutoffTime_
		4,                // maxSyncIterations_
	)
	require.NoError(t, err)
	b.Commit()

	referenceRateAdapterAddress, _, referenceRateAdapter, err := generated.DeployReferenceRateAdapterMock(owner, b.Client(), AdaptiveOracleDecimals)
	require.NoError(t, err)
	b.Commit()
	_, err = referenceRateAdapter.SetRate(owner, big.NewInt(referenceRate), true)
	require.NoError(t, err)
	b.Commit()

	adaptiveRateLogicAddress, _, _, err := generated.DeployAdaptiveRateLogic(owner, b.Client())
	require.NoError(t, err)
	b.Commit()

	adaptiveOracleAddress, _, adaptiveOracle, err := generated.DeployAdaptiveOracle(owner, b.Client(), owner.From, AdaptiveOracleDecimals, "ADAPTIVE")
	require.NoError(t, err)
	b.Commit()

	_, err = adaptiveOracle.SetAggregator(owner, dualAggregatorAddress)
	require.NoError(t, err)
	b.Commit()
	_, err = adaptiveOracle.SetReferenceRateAdapter(owner, referenceRateAdapterAddress)
	require.NoError(t, err)
	b.Commit()
	_, err = adaptiveOracle.SetAdaptiveRateLogic(owner, adaptiveRateLogicAddress)
	require.NoError(t, err)
	b.Commit()
	_, err = dualAggregator.SetAdaptiveOracle(owner, adaptiveOracleAddress)
	require.NoError(t, err)
	b.Commit()

	// Ensure we have finality depth worth of blocks to start.
	for range 20 {
		b.Commit()
	}

	_, err = linkContract.Transfer(owner, dualAggregatorAddress, big.NewInt(1000))
	require.NoError(t, err)
	b.Commit()

	return owner, b, &nodeConfig, &Contracts{
		DualAggregatorAddress:    dualAggregatorAddress,
		DualAggregator:           dualAggregator,
		AdaptiveOracleAddress:    adaptiveOracleAddress,
		AdaptiveOracle:           adaptiveOracle,
		AdaptiveRateLogicAddress: adaptiveRateLogicAddress,
		ReferenceRateAdapter:     referenceRateAdapter,
	}
}

func withRPCServer(host string, httpPort, wsPort int, modules []string) func(nodeConf *gethnode.Config, ethConf *ethconfig.Config) {
	return func(nodeConf *gethnode.Config, ethConf *ethconfig.Config) {
		nodeConf.HTTPHost = host
		nodeConf.HTTPPort = httpPort
		nodeConf.HTTPModules = modules

		nodeConf.WSHost = host
		nodeConf.WSPort = wsPort
		nodeConf.WSModules = modules
	}
}

// InitAdaptiveOracle sets payees/config on the DualAggregator (identical OCR2Abstract semantics to
// the standard OCR2Aggregator) and starts the bootstrap node's job.
func InitAdaptiveOracle(
	t *testing.T,
	lggr logger.Logger,
	b *simulated.Backend,
	dualAggregator *generated.DualAggregator,
	owner *bind.TransactOpts,
	bootstrapNode *ocr2.Node,
	oracles []confighelper2.OracleIdentityExtra,
	transmitters []common.Address,
	payees []common.Address,
	specFn func(int64) string,
) (blockBeforeConfig *types.Block) {
	_, err := dualAggregator.SetPayees(owner, transmitters, payees)
	require.NoError(t, err)
	b.Commit()
	blockBeforeConfig, err = b.Client().BlockByNumber(t.Context(), nil)
	require.NoError(t, err)

	signers, effectiveTransmitters, threshold, _, encodedConfigVersion, encodedConfig, err := confighelper2.ContractSetConfigArgsForEthereumIntegrationTest(
		oracles,
		1,
		1000000000/100, // threshold PPB
	)
	require.NoError(t, err)

	minAnswer, maxAnswer := new(big.Int), new(big.Int)
	minAnswer.Exp(big.NewInt(-2), big.NewInt(191), nil)
	maxAnswer.Exp(big.NewInt(2), big.NewInt(191), nil)
	maxAnswer.Sub(maxAnswer, big.NewInt(1))

	onchainConfig, err := testhelpers.GenerateDefaultOCR2OnchainConfig(minAnswer, maxAnswer)
	require.NoError(t, err)

	lggr.Debugw("Setting Config on DualAggregator",
		"signers", signers,
		"transmitters", transmitters,
		"effectiveTransmitters", effectiveTransmitters,
		"threshold", threshold,
	)
	_, err = dualAggregator.SetConfig(owner, signers, effectiveTransmitters, threshold, onchainConfig, encodedConfigVersion, encodedConfig)
	require.NoError(t, err)
	b.Commit()

	require.NoError(t, bootstrapNode.App.Start(t.Context()))

	ocrJob, err := ocrbootstrap.ValidatedBootstrapSpecToml(specFn(blockBeforeConfig.Number().Int64()))
	require.NoError(t, err)
	require.NoError(t, bootstrapNode.App.AddJobV2(t.Context(), &ocrJob))
	return
}

// medianJobToml returns an offchainreporting2/median job spec identical in shape to the one used by
// the base OCR2 features test, pointed at the DualAggregator and reading a fixed price from a bridge.
func medianJobToml(contractAddress common.Address, kbID, transmitterID, bridgeName string, fromBlock int64) string {
	return fmt.Sprintf(`
type               = "offchainreporting2"
relay              = "evm"
schemaVersion      = 1
pluginType         = "median"
name               = "adaptive oracle market rate spec"
contractID         = "%s"
ocrKeyBundleID     = "%s"
transmitterID      = "%s"
contractConfigConfirmations = 1
contractConfigTrackerPollInterval = "1s"
observationSource  = """
    ds1          [type=bridge name="%s"];
    ds1_parse    [type=jsonparse path="data"];
    ds1 -> ds1_parse -> answer1;
    answer1 [type=median index=0];
"""
[relayConfig]
chainID = 1337
fromBlock = %d
[pluginConfig]
juelsPerFeeCoinSource = """
    ds1          [type=bridge name="%s"];
    ds1_parse    [type=jsonparse path="data"];
    ds1 -> ds1_parse -> answer1;
    answer1 [type=median index=0];
"""
[pluginConfig.juelsPerFeeCoinCache]
updateInterval = "1m"
`, contractAddress, kbID, transmitterID, bridgeName, fromBlock, bridgeName)
}
