package evm

import (
	"database/sql"
	"math/big"
	"testing"
	"time"

	"github.com/ethereum/go-ethereum"
	"github.com/ethereum/go-ethereum/accounts/abi/bind"
	"github.com/ethereum/go-ethereum/common"
	"github.com/ethereum/go-ethereum/core/types"
	"github.com/ethereum/go-ethereum/crypto"
	"github.com/ethereum/go-ethereum/eth/ethconfig"
	"github.com/ethereum/go-ethereum/ethclient/simulated"
	"github.com/onsi/gomega"
	"github.com/pkg/errors"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
	"github.com/stretchr/testify/require"

	"github.com/smartcontractkit/libocr/gethwrappers2/ocr2aggregator"
	"github.com/smartcontractkit/libocr/gethwrappers2/ocrconfigurationstoreevmsimple"
	testoffchainaggregator2 "github.com/smartcontractkit/libocr/gethwrappers2/testocr2aggregator"
	"github.com/smartcontractkit/libocr/offchainreporting2/reportingplugin/median"
	confighelper2 "github.com/smartcontractkit/libocr/offchainreporting2plus/confighelper"
	ocrtypes "github.com/smartcontractkit/libocr/offchainreporting2plus/types"
	ocrtypes2 "github.com/smartcontractkit/libocr/offchainreporting2plus/types"

	"github.com/smartcontractkit/chainlink-common/pkg/logger"
	"github.com/smartcontractkit/chainlink-common/pkg/services/servicetest"

	"github.com/smartcontractkit/chainlink-evm/gethwrappers/generated/link_token_interface"
	"github.com/smartcontractkit/chainlink-evm/pkg/client"
	evmclient "github.com/smartcontractkit/chainlink-evm/pkg/client"
	"github.com/smartcontractkit/chainlink-evm/pkg/client/clienttest"
	"github.com/smartcontractkit/chainlink-evm/pkg/heads/headstest"
	"github.com/smartcontractkit/chainlink-evm/pkg/logpoller"
	evmutils "github.com/smartcontractkit/chainlink-evm/pkg/utils"
	"github.com/smartcontractkit/chainlink/v2/core/chains/evm/logpoller/mocks"
	"github.com/smartcontractkit/chainlink/v2/core/internal/testutils"
	"github.com/smartcontractkit/chainlink/v2/core/internal/testutils/pgtest"
	"github.com/smartcontractkit/chainlink/v2/core/services/ocr2/testhelpers"
	"github.com/smartcontractkit/chainlink/v2/core/utils"
)

func TestConfigPoller(t *testing.T) {
	lggr := logger.Test(t)
	var ethClient *client.SimulatedBackendClient
	var lp logpoller.LogPoller
	var ocrAddress common.Address
	var ocrContract *ocr2aggregator.OCR2Aggregator
	var configStoreContractAddr common.Address
	var configStoreContract *ocrconfigurationstoreevmsimple.OCRConfigurationStoreEVMSimple
	var user *bind.TransactOpts
	var b *simulated.Backend
	var ec simulated.Client
	var linkTokenAddress common.Address
	var accessAddress common.Address
	ctx := testutils.Context(t)

	ld := OCR2AggregatorLogDecoder

	{
		key, err := crypto.GenerateKey()
		require.NoError(t, err)
		user, err = bind.NewKeyedTransactorWithChainID(key, big.NewInt(1337))
		require.NoError(t, err)
		b = simulated.NewBackend(types.GenesisAlloc{
			user.From: {Balance: big.NewInt(1000000000000000000)}},
			simulated.WithBlockGasLimit(5*ethconfig.Defaults.Miner.GasCeil))
		require.NotNil(t, b)
		ec = b.Client()
		linkTokenAddress, _, _, err = link_token_interface.DeployLinkToken(user, ec)
		require.NoError(t, err)
		accessAddress, _, _, err = testoffchainaggregator2.DeploySimpleWriteAccessController(user, ec)
		require.NoError(t, err, "failed to deploy test access controller contract")
		ocrAddress, _, ocrContract, err = ocr2aggregator.DeployOCR2Aggregator(
			user,
			ec,
			linkTokenAddress,
			big.NewInt(0),
			big.NewInt(10),
			accessAddress,
			accessAddress,
			9,
			"TEST",
		)
		require.NoError(t, err)
		configStoreContractAddr, _, configStoreContract, err = ocrconfigurationstoreevmsimple.DeployOCRConfigurationStoreEVMSimple(user, ec)
		require.NoError(t, err)
		b.Commit()

		db := pgtest.NewSqlxDB(t)
		ethClient = evmclient.NewSimulatedBackendClient(t, b, testutils.SimulatedChainID)

		lorm := logpoller.NewORM(testutils.SimulatedChainID, db, lggr)

		lpOpts := logpoller.Opts{
			PollPeriod:               100 * time.Millisecond,
			FinalityDepth:            1,
			BackfillBatchSize:        2,
			RPCBatchSize:             2,
			KeepFinalizedBlocksDepth: 1000,
		}
		ht := headstest.NewSimulatedHeadTracker(ethClient, lpOpts.UseFinalityTag, lpOpts.FinalityDepth)
		lp = logpoller.NewLogPoller(lorm, ethClient, lggr, ht, lpOpts)
		servicetest.Run(t, lp)
	}

	t.Run("LatestConfig errors if there is no config in logs and config store is unconfigured", func(t *testing.T) {
		cp, err := NewConfigPoller(ctx, lggr, CPConfig{ethClient, lp, ocrAddress, nil, ld})
		require.NoError(t, err)

		_, err = cp.LatestConfig(testutils.Context(t), 0)
		require.Error(t, err)
		assert.Contains(t, err.Error(), "no logs found for config on contract")
	})

	t.Run("happy path (with config store)", func(t *testing.T) {
		cp, err := NewConfigPoller(ctx, lggr, CPConfig{ethClient, lp, ocrAddress, &configStoreContractAddr, ld})
		require.NoError(t, err)
		// Should have no config to begin with.
		_, configDigest, err := cp.LatestConfigDetails(testutils.Context(t))
		require.NoError(t, err)
		require.Equal(t, ocrtypes2.ConfigDigest{}, configDigest)
		// Should error because there are no logs for config at block 0
		_, err = cp.LatestConfig(testutils.Context(t), 0)
		require.Error(t, err)
		assert.Contains(t, err.Error(), "config details missing while trying to lookup config in store")

		// Set the config
		contractConfig := setConfig(t, median.OffchainConfig{
			AlphaReportInfinite: false,
			AlphaReportPPB:      0,
			AlphaAcceptInfinite: true,
			AlphaAcceptPPB:      0,
			DeltaC:              10,
		}, ocrContract, user)
		b.Commit()
		latest, err := ec.BlockByNumber(testutils.Context(t), nil)
		require.NoError(t, err)
		// Ensure we capture this config set log.
		require.NoError(t, lp.Replay(testutils.Context(t), latest.Number().Int64()-1))

		// Send blocks until we see the config updated.
		var configBlock uint64
		var digest [32]byte
		gomega.NewGomegaWithT(t).Eventually(func() bool {
			b.Commit()
			configBlock, digest, err = cp.LatestConfigDetails(testutils.Context(t))
			require.NoError(t, err)
			return ocrtypes2.ConfigDigest{} != digest
		}, testutils.WaitTimeout(t), 100*time.Millisecond).Should(gomega.BeTrue())

		// Assert the config returned is the one we configured.
		newConfig, err := cp.LatestConfig(testutils.Context(t), configBlock)
		require.NoError(t, err)
		// Note we don't check onchainConfig, as that is populated in the contract itself.
		assert.Equal(t, digest, [32]byte(newConfig.ConfigDigest))
		assert.Equal(t, contractConfig.Signers, newConfig.Signers)
		assert.Equal(t, contractConfig.Transmitters, newConfig.Transmitters)
		assert.Equal(t, contractConfig.F, newConfig.F)
		assert.Equal(t, contractConfig.OffchainConfigVersion, newConfig.OffchainConfigVersion)
		assert.Equal(t, contractConfig.OffchainConfig, newConfig.OffchainConfig)
	})

	{
		var err error
		ocrAddress, _, ocrContract, err = ocr2aggregator.DeployOCR2Aggregator(
			user,
			ec,
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
	}

	t.Run("LatestConfigDetails, when logs have been pruned and config store contract is configured", func(t *testing.T) {
		// Give it a log poller that will never return logs
		mp := new(mocks.LogPoller)
		mp.On("RegisterFilter", mock.Anything, mock.Anything).Return(nil)
		mp.On("LatestLogByEventSigWithConfs", mock.Anything, mock.Anything, mock.Anything, mock.Anything).Return(nil, sql.ErrNoRows)

		t.Run("if callLatestConfigDetails succeeds", func(t *testing.T) {
			cp, err := newConfigPoller(ctx, lggr, ethClient, mp, ocrAddress, &configStoreContractAddr, ld)
			require.NoError(t, err)

			t.Run("when config has not been set, returns zero values", func(t *testing.T) {
				changedInBlock, configDigest, err := cp.LatestConfigDetails(testutils.Context(t))
				require.NoError(t, err)

				assert.Equal(t, 0, int(changedInBlock))
				assert.Equal(t, ocrtypes.ConfigDigest{}, configDigest)
			})
			t.Run("when config has been set, returns config details", func(t *testing.T) {
				setConfig(t, median.OffchainConfig{
					AlphaReportInfinite: false,
					AlphaReportPPB:      0,
					AlphaAcceptInfinite: true,
					AlphaAcceptPPB:      0,
					DeltaC:              10,
				}, ocrContract, user)
				b.Commit()

				changedInBlock, configDigest, err := cp.LatestConfigDetails(testutils.Context(t))
				require.NoError(t, err)

				latest, err := ec.BlockByNumber(testutils.Context(t), nil)
				require.NoError(t, err)

				onchainDetails, err := ocrContract.LatestConfigDetails(nil)
				require.NoError(t, err)

				assert.Equal(t, latest.Number().Int64(), int64(changedInBlock))
				assert.Equal(t, onchainDetails.ConfigDigest, [32]byte(configDigest))
			})
		})
		t.Run("returns error if callLatestConfigDetails fails", func(t *testing.T) {
			failingClient := new(clienttest.Client)
			failingClient.On("ConfiguredChainID").Return(big.NewInt(42))
			failingClient.On("CallContract", mock.Anything, mock.Anything, mock.Anything).Return(nil, errors.New("something exploded"))
			cp, err := newConfigPoller(ctx, lggr, failingClient, mp, ocrAddress, &configStoreContractAddr, ld)
			require.NoError(t, err)

			cp.configStoreContractAddr = &configStoreContractAddr
			cp.configStoreContract = configStoreContract

			_, _, err = cp.LatestConfigDetails(testutils.Context(t))
			assert.EqualError(t, err, "something exploded")

			failingClient.AssertExpectations(t)
		})
	})

	{
		var err error
		// deploy it again to reset to empty config
		ocrAddress, _, ocrContract, err = ocr2aggregator.DeployOCR2Aggregator(
			user,
			ec,
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
	}

	t.Run("LatestConfig, when logs have been pruned and config store contract is configured", func(t *testing.T) {
		// Give it a log poller that will never return logs
		mp := mocks.NewLogPoller(t)
		mp.On("RegisterFilter", mock.Anything, mock.Anything).Return(nil)
		mp.On("Logs", mock.Anything, mock.Anything, mock.Anything, mock.Anything, mock.Anything).Return(nil, nil)
		mp.On("LatestLogByEventSigWithConfs", mock.Anything, mock.Anything, mock.Anything, mock.Anything).Return(nil, sql.ErrNoRows)

		t.Run("if callReadConfig succeeds", func(t *testing.T) {
			cp, err := newConfigPoller(ctx, lggr, ethClient, mp, ocrAddress, &configStoreContractAddr, ld)
			require.NoError(t, err)

			t.Run("when config has not been set, returns error", func(t *testing.T) {
				_, err := cp.LatestConfig(testutils.Context(t), 0)
				require.Error(t, err)

				assert.Contains(t, err.Error(), "config details missing while trying to lookup config in store")
			})
			t.Run("when config has been set, returns config", func(t *testing.T) {
				b.Commit()
				onchainDetails, err := ocrContract.LatestConfigDetails(nil)
				require.NoError(t, err)

				contractConfig := setConfig(t, median.OffchainConfig{
					AlphaReportInfinite: false,
					AlphaReportPPB:      0,
					AlphaAcceptInfinite: true,
					AlphaAcceptPPB:      0,
					DeltaC:              10,
				}, ocrContract, user)

				signerAddresses, err := OnchainPublicKeyToAddress(contractConfig.Signers)
				require.NoError(t, err)
				transmitterAddresses, err := AccountToAddress(contractConfig.Transmitters)
				require.NoError(t, err)

				configuration := ocrconfigurationstoreevmsimple.OCRConfigurationStoreEVMSimpleConfigurationEVMSimple{
					Signers:               signerAddresses,
					Transmitters:          transmitterAddresses,
					OnchainConfig:         contractConfig.OnchainConfig,
					OffchainConfig:        contractConfig.OffchainConfig,
					ContractAddress:       ocrAddress,
					OffchainConfigVersion: contractConfig.OffchainConfigVersion,
					ConfigCount:           1,
					F:                     contractConfig.F,
				}

				addConfig(t, user, configStoreContract, configuration)

				b.Commit()
				onchainDetails, err = ocrContract.LatestConfigDetails(nil)
				require.NoError(t, err)

				newConfig, err := cp.LatestConfig(testutils.Context(t), 0)
				require.NoError(t, err)

				assert.Equal(t, onchainDetails.ConfigDigest, [32]byte(newConfig.ConfigDigest))
				assert.Equal(t, contractConfig.Signers, newConfig.Signers)
				assert.Equal(t, contractConfig.Transmitters, newConfig.Transmitters)
				assert.Equal(t, contractConfig.F, newConfig.F)
				assert.Equal(t, contractConfig.OffchainConfigVersion, newConfig.OffchainConfigVersion)
				assert.Equal(t, contractConfig.OffchainConfig, newConfig.OffchainConfig)
			})
		})
		t.Run("returns error if callReadConfig fails", func(t *testing.T) {
			failingClient := new(clienttest.Client)
			failingClient.On("ConfiguredChainID").Return(big.NewInt(42))
			failingClient.On("CallContract", mock.Anything, mock.MatchedBy(func(callArgs ethereum.CallMsg) bool {
				// initial call to retrieve config store address from aggregator
				return *callArgs.To == ocrAddress
			}), mock.Anything).Return(nil, errors.New("something exploded")).Once()
			cp, err := newConfigPoller(ctx, lggr, failingClient, mp, ocrAddress, &configStoreContractAddr, ld)
			require.NoError(t, err)

			_, err = cp.LatestConfig(testutils.Context(t), 0)
			assert.EqualError(t, err, "failed to get latest config details: something exploded")

			failingClient.AssertExpectations(t)
		})
	})
}

func setConfig(t *testing.T, pluginConfig median.OffchainConfig, ocrContract *ocr2aggregator.OCR2Aggregator, user *bind.TransactOpts) ocrtypes2.ContractConfig {
	// Create minimum number of nodes.
	var oracles []confighelper2.OracleIdentityExtra
	for i := 0; i < 4; i++ {
		oracles = append(oracles, confighelper2.OracleIdentityExtra{
			OracleIdentity: confighelper2.OracleIdentity{
				OnchainPublicKey:  evmutils.RandomAddress().Bytes(),
				TransmitAccount:   ocrtypes2.Account(evmutils.RandomAddress().Hex()),
				OffchainPublicKey: evmutils.RandomBytes32(),
				PeerID:            utils.MustNewPeerID(),
			},
			ConfigEncryptionPublicKey: evmutils.RandomBytes32(),
		})
	}
	// Gnerate OnchainConfig
	onchainConfig, err := testhelpers.GenerateDefaultOCR2OnchainConfig(big.NewInt(0), big.NewInt(10))
	require.NoError(t, err)
	// Change the offramp config
	signers, transmitters, threshold, onchainConfig, offchainConfigVersion, offchainConfig, err := confighelper2.ContractSetConfigArgsForTests(
		2*time.Second,        // deltaProgress
		1*time.Second,        // deltaResend
		1*time.Second,        // deltaRound
		500*time.Millisecond, // deltaGrace
		2*time.Second,        // deltaStage
		3,
		[]int{1, 1, 1, 1},
		oracles,
		pluginConfig.Encode(),
		nil,
		50*time.Millisecond,
		50*time.Millisecond,
		50*time.Millisecond,
		50*time.Millisecond,
		50*time.Millisecond,
		1, // faults
		onchainConfig,
	)
	require.NoError(t, err)
	signerAddresses, err := OnchainPublicKeyToAddress(signers)
	require.NoError(t, err)
	transmitterAddresses, err := AccountToAddress(transmitters)
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

func addConfig(t *testing.T, user *bind.TransactOpts, configStoreContract *ocrconfigurationstoreevmsimple.OCRConfigurationStoreEVMSimple, config ocrconfigurationstoreevmsimple.OCRConfigurationStoreEVMSimpleConfigurationEVMSimple) {
	_, err := configStoreContract.AddConfig(user, config)
	require.NoError(t, err)
}
