package functions_test

import (
	"encoding/binary"
	"math/big"
	"testing"
	"time"

	"github.com/ethereum/go-ethereum/accounts/abi/bind"
	"github.com/ethereum/go-ethereum/accounts/abi/bind/backends"
	"github.com/ethereum/go-ethereum/core"
	"github.com/ethereum/go-ethereum/crypto"
	"github.com/ethereum/go-ethereum/eth/ethconfig"
	"github.com/onsi/gomega"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"

	"github.com/smartcontractkit/libocr/gethwrappers2/ocr2aggregator"
	testoffchainaggregator2 "github.com/smartcontractkit/libocr/gethwrappers2/testocr2aggregator"
	confighelper2 "github.com/smartcontractkit/libocr/offchainreporting2plus/confighelper"
	ocrtypes2 "github.com/smartcontractkit/libocr/offchainreporting2plus/types"

	functionsConfig "github.com/smartcontractkit/chainlink/v2/core/services/ocr2/plugins/functions/config"

	evmclient "github.com/smartcontractkit/chainlink/v2/core/chains/evm/client"
	"github.com/smartcontractkit/chainlink/v2/core/chains/evm/logpoller"
	"github.com/smartcontractkit/chainlink/v2/core/gethwrappers/generated/link_token_interface"
	"github.com/smartcontractkit/chainlink/v2/core/internal/testutils"
	"github.com/smartcontractkit/chainlink/v2/core/internal/testutils/pgtest"
	"github.com/smartcontractkit/chainlink/v2/core/logger"
	"github.com/smartcontractkit/chainlink/v2/core/services/relay/evm"
	"github.com/smartcontractkit/chainlink/v2/core/services/relay/evm/functions"
	"github.com/smartcontractkit/chainlink/v2/core/utils"
)

func TestFunctionsConfigPoller(t *testing.T) {
	t.Run("FunctionsPlugin", func(t *testing.T) {
		runTest(t, functions.FunctionsPlugin, functions.FunctionsDigestPrefix)
	})
	t.Run("ThresholdPlugin", func(t *testing.T) {
		runTest(t, functions.ThresholdPlugin, functions.ThresholdDigestPrefix)
	})
	// TODO: Test config poller for S4Plugin (requires S4Plugin to be implemented & corresponding updates to pluginConfig)
}

func runTest(t *testing.T, pluginType functions.FunctionsPluginType, expectedDigestPrefix ocrtypes2.ConfigDigestPrefix) {
	key, err := crypto.GenerateKey()
	require.NoError(t, err)
	user, err := bind.NewKeyedTransactorWithChainID(key, big.NewInt(1337))
	require.NoError(t, err)
	b := backends.NewSimulatedBackend(core.GenesisAlloc{
		user.From: {Balance: big.NewInt(1000000000000000000)}},
		5*ethconfig.Defaults.Miner.GasCeil)
	defer b.Close()
	linkTokenAddress, _, _, err := link_token_interface.DeployLinkToken(user, b)
	require.NoError(t, err)
	accessAddress, _, _, err := testoffchainaggregator2.DeploySimpleWriteAccessController(user, b)
	require.NoError(t, err, "failed to deploy test access controller contract")
	ocrAddress, _, ocrContract, err := ocr2aggregator.DeployOCR2Aggregator(
		user,
		b,
		linkTokenAddress,
		big.NewInt(0),
		big.NewInt(10),
		accessAddress,
		accessAddress,
		9,
		"TEST",
	)
	require.NoError(t, err)
	b.Commit()
	db := pgtest.NewSqlxDB(t)
	defer db.Close()
	cfg := pgtest.NewQConfig(false)
	ethClient := evmclient.NewSimulatedBackendClient(t, b, big.NewInt(1337))
	defer ethClient.Close()
	lggr := logger.TestLogger(t)
	ctx := testutils.Context(t)
	lorm := logpoller.NewORM(big.NewInt(1337), db, lggr, cfg)
	lp := logpoller.NewLogPoller(lorm, ethClient, lggr, 100*time.Millisecond, 1, 2, 2, 1000)
	defer lp.Close()
	require.NoError(t, lp.Start(ctx))
	logPoller, err := functions.NewFunctionsConfigPoller(pluginType, lp, ocrAddress, lggr)
	require.NoError(t, err)
	// Should have no config to begin with.
	_, config, err := logPoller.LatestConfigDetails(testutils.Context(t))
	require.NoError(t, err)
	require.Equal(t, ocrtypes2.ConfigDigest{}, config)

	pluginConfig := &functionsConfig.ReportingPluginConfigWrapper{
		Config: &functionsConfig.ReportingPluginConfig{
			MaxQueryLengthBytes:       10000,
			MaxObservationLengthBytes: 10000,
			MaxReportLengthBytes:      10000,
			MaxRequestBatchSize:       10,
			DefaultAggregationMethod:  functionsConfig.AggregationMethod(0),
			UniqueReports:             true,
			ThresholdPluginConfig: &functionsConfig.ThresholdReportingPluginConfig{
				MaxQueryLengthBytes:       10000,
				MaxObservationLengthBytes: 10000,
				MaxReportLengthBytes:      10000,
				RequestCountLimit:         100,
				RequestTotalBytesLimit:    100000,
				RequireLocalRequestCheck:  true,
			},
		},
	}

	// Set the config
	contractConfig := setFunctionsConfig(t, pluginConfig, ocrContract, user)
	b.Commit()
	latest, err := b.BlockByNumber(testutils.Context(t), nil)
	require.NoError(t, err)
	// Ensure we capture this config set log.
	require.NoError(t, lp.Replay(testutils.Context(t), latest.Number().Int64()-1))

	// Send blocks until we see the config updated.
	var configBlock uint64
	var digest [32]byte
	gomega.NewGomegaWithT(t).Eventually(func() bool {
		b.Commit()
		configBlock, digest, err = logPoller.LatestConfigDetails(testutils.Context(t))
		require.NoError(t, err)
		return ocrtypes2.ConfigDigest{} != digest
	}, testutils.WaitTimeout(t), 100*time.Millisecond).Should(gomega.BeTrue())

	// Assert the config returned is the one we configured.
	newConfig, err := logPoller.LatestConfig(testutils.Context(t), configBlock)
	require.NoError(t, err)

	// Get actual configDigest value from contracts
	configFromContract, err := ocrContract.LatestConfigDetails(nil)
	require.NoError(t, err)
	onChainConfigDigest := configFromContract.ConfigDigest

	assert.Equal(t, contractConfig.Signers, newConfig.Signers)
	assert.Equal(t, contractConfig.Transmitters, newConfig.Transmitters)
	assert.Equal(t, contractConfig.F, newConfig.F)
	assert.Equal(t, contractConfig.OffchainConfigVersion, newConfig.OffchainConfigVersion)
	assert.Equal(t, contractConfig.OffchainConfig, newConfig.OffchainConfig)

	var expectedConfigDigest [32]byte
	copy(expectedConfigDigest[:], onChainConfigDigest[:])
	binary.BigEndian.PutUint16(expectedConfigDigest[:2], uint16(expectedDigestPrefix))

	assert.Equal(t, expectedConfigDigest, digest)
	assert.Equal(t, expectedConfigDigest, [32]byte(newConfig.ConfigDigest))
}

func setFunctionsConfig(t *testing.T, pluginConfig *functionsConfig.ReportingPluginConfigWrapper, ocrContract *ocr2aggregator.OCR2Aggregator, user *bind.TransactOpts) ocrtypes2.ContractConfig {
	// Create minimum number of nodes.
	var oracles []confighelper2.OracleIdentityExtra
	for i := 0; i < 4; i++ {
		oracles = append(oracles, confighelper2.OracleIdentityExtra{
			OracleIdentity: confighelper2.OracleIdentity{
				OnchainPublicKey:  utils.RandomAddress().Bytes(),
				TransmitAccount:   ocrtypes2.Account(utils.RandomAddress().String()),
				OffchainPublicKey: utils.RandomBytes32(),
				PeerID:            utils.MustNewPeerID(),
			},
			ConfigEncryptionPublicKey: utils.RandomBytes32(),
		})
	}

	pluginConfigBytes, err := functionsConfig.EncodeReportingPluginConfig(pluginConfig)
	require.NoError(t, err)

	signers, transmitters, threshold, onchainConfig, offchainConfigVersion, offchainConfig, err := confighelper2.ContractSetConfigArgsForTests(
		2*time.Second,        // deltaProgress
		1*time.Second,        // deltaResend
		1*time.Second,        // deltaRound
		500*time.Millisecond, // deltaGrace
		2*time.Second,        // deltaStage
		3,
		[]int{1, 1, 1, 1},
		oracles,
		pluginConfigBytes,
		50*time.Millisecond,
		50*time.Millisecond,
		50*time.Millisecond,
		50*time.Millisecond,
		50*time.Millisecond,
		1, // faults
		nil,
	)

	require.NoError(t, err)
	signerAddresses, err := evm.OnchainPublicKeyToAddress(signers)
	require.NoError(t, err)
	transmitterAddresses, err := evm.AccountToAddress(transmitters)
	require.NoError(t, err)
	_, err = ocrContract.SetConfig(user, signerAddresses, transmitterAddresses, threshold, onchainConfig, offchainConfigVersion, offchainConfig)
	require.NoError(t, err)
	return ocrtypes2.ContractConfig{
		Signers:               signers,
		Transmitters:          transmitters,
		F:                     threshold,
		OnchainConfig:         onchainConfig,
		OffchainConfigVersion: offchainConfigVersion,
		OffchainConfig:        offchainConfig,
	}
}
