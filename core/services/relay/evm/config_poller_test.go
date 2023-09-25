package evm

import (
	"database/sql"
	"math/big"
	"testing"
	"time"

	"github.com/ethereum/go-ethereum"
	"github.com/smartcontractkit/libocr/bigbigendian"
	"github.com/smartcontractkit/libocr/gethwrappers2/ocrconfigurationstoreevmsimple"

	"github.com/ethereum/go-ethereum/accounts/abi/bind"
	"github.com/ethereum/go-ethereum/accounts/abi/bind/backends"
	"github.com/ethereum/go-ethereum/core"
	"github.com/ethereum/go-ethereum/crypto"
	"github.com/ethereum/go-ethereum/eth/ethconfig"
	"github.com/onsi/gomega"
	"github.com/pkg/errors"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
	"github.com/stretchr/testify/require"

	"github.com/smartcontractkit/libocr/gethwrappers2/ocr2aggregator"
	testoffchainaggregator2 "github.com/smartcontractkit/libocr/gethwrappers2/testocr2aggregator"
	"github.com/smartcontractkit/libocr/offchainreporting2/reportingplugin/median"
	confighelper2 "github.com/smartcontractkit/libocr/offchainreporting2plus/confighelper"
	ocrtypes "github.com/smartcontractkit/libocr/offchainreporting2plus/types"
	ocrtypes2 "github.com/smartcontractkit/libocr/offchainreporting2plus/types"

	evmclient "github.com/smartcontractkit/chainlink/v2/core/chains/evm/client"
	evmClientMocks "github.com/smartcontractkit/chainlink/v2/core/chains/evm/client/mocks"
	"github.com/smartcontractkit/chainlink/v2/core/chains/evm/logpoller"
	"github.com/smartcontractkit/chainlink/v2/core/chains/evm/logpoller/mocks"
	"github.com/smartcontractkit/chainlink/v2/core/gethwrappers/generated/link_token_interface"
	"github.com/smartcontractkit/chainlink/v2/core/internal/testutils"
	"github.com/smartcontractkit/chainlink/v2/core/internal/testutils/pgtest"
	"github.com/smartcontractkit/chainlink/v2/core/logger"
	evmRelayTypes "github.com/smartcontractkit/chainlink/v2/core/services/relay/evm/types"
	"github.com/smartcontractkit/chainlink/v2/core/utils"
)

func TestConfigPoller(t *testing.T) {
	key, err := crypto.GenerateKey()
	require.NoError(t, err)
	user, err := bind.NewKeyedTransactorWithChainID(key, big.NewInt(1337))
	require.NoError(t, err)
	b := backends.NewSimulatedBackend(core.GenesisAlloc{
		user.From: {Balance: big.NewInt(1000000000000000000)}},
		5*ethconfig.Defaults.Miner.GasCeil)
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
	configStoreContractAddr, _, configStoreContract, err := ocrconfigurationstoreevmsimple.DeployOCRConfigurationStoreEVMSimple(user, b)
	require.NoError(t, err)
	b.Commit()

	db := pgtest.NewSqlxDB(t)
	cfg := pgtest.NewQConfig(false)
	ethClient := evmclient.NewSimulatedBackendClient(t, b, testutils.SimulatedChainID)
	lggr := logger.TestLogger(t)
	ctx := testutils.Context(t)
	lorm := logpoller.NewORM(testutils.SimulatedChainID, db, lggr, cfg)
	lp := logpoller.NewLogPoller(lorm, ethClient, lggr, 100*time.Millisecond, 1, 2, 2, 1000)
	require.NoError(t, lp.Start(ctx))
	t.Cleanup(func() { lp.Close() })

	t.Run("happy path", func(t *testing.T) {
		var cp evmRelayTypes.ConfigPoller
		cp, err = NewConfigPoller(lggr, ethClient, lp, ocrAddress, &configStoreContractAddr)
		require.NoError(t, err)
		// Should have no config to begin with.
		_, config, err := cp.LatestConfigDetails(testutils.Context(t))
		require.NoError(t, err)
		require.Equal(t, ocrtypes2.ConfigDigest{}, config)
		_, err = cp.LatestConfig(testutils.Context(t), 0)
		require.Error(t, err)
		// Set the config
		contractConfig := setConfig(t, median.OffchainConfig{
			AlphaReportInfinite: false,
			AlphaReportPPB:      0,
			AlphaAcceptInfinite: true,
			AlphaAcceptPPB:      0,
			DeltaC:              10,
		}, ocrContract, user)
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

	ocrAddress, _, ocrContract, err = ocr2aggregator.DeployOCR2Aggregator(
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

	t.Run("LatestConfigDetails, when logs have been pruned and config store contract is configured", func(t *testing.T) {
		// Give it a log poller that will never return logs
		mp := new(mocks.LogPoller)
		mp.On("RegisterFilter", mock.Anything).Return(nil)
		mp.On("LatestLogByEventSigWithConfs", mock.Anything, mock.Anything, mock.Anything, mock.Anything).Return(nil, sql.ErrNoRows)

		t.Run("if callLatestConfigDetails succeeds", func(t *testing.T) {
			var cp evmRelayTypes.ConfigPoller
			cp, err = newConfigPoller(lggr, ethClient, mp, ocrAddress, &configStoreContractAddr)
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

				latest, err := b.BlockByNumber(testutils.Context(t), nil)
				require.NoError(t, err)

				onchainDetails, err := ocrContract.LatestConfigDetails(nil)
				require.NoError(t, err)

				assert.Equal(t, latest.Number().Int64(), int64(changedInBlock))
				assert.Equal(t, onchainDetails.ConfigDigest, [32]byte(configDigest))
			})
		})
		t.Run("returns error if callLatestConfigDetails fails", func(t *testing.T) {
			failingClient := new(evmClientMocks.Client)
			failingClient.On("ConfiguredChainID").Return(big.NewInt(42))
			failingClient.On("CallContract", mock.Anything, mock.Anything, mock.Anything).Return(nil, errors.New("something exploded"))
			var cp *configPoller
			cp, err = newConfigPoller(lggr, failingClient, mp, ocrAddress, &configStoreContractAddr)
			require.NoError(t, err)

			cp.configStoreContractAddr = &configStoreContractAddr
			cp.configStoreContract = configStoreContract

			_, _, err = cp.LatestConfigDetails(testutils.Context(t))
			assert.EqualError(t, err, "something exploded")

			failingClient.AssertExpectations(t)
		})
	})

	// deploy it again to reset to empty config
	ocrAddress, _, ocrContract, err = ocr2aggregator.DeployOCR2Aggregator(
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

	t.Run("LatestConfig, when logs have been pruned and config store contract is configured", func(t *testing.T) {
		// Give it a log poller that will never return logs
		mp := mocks.NewLogPoller(t)
		mp.On("RegisterFilter", mock.Anything).Return(nil)
		mp.On("Logs", mock.Anything, mock.Anything, mock.Anything, mock.Anything, mock.Anything).Return(nil, nil)
		mp.On("LatestLogByEventSigWithConfs", mock.Anything, mock.Anything, mock.Anything, mock.Anything).Return(nil, sql.ErrNoRows)

		t.Run("if callReadConfig succeeds", func(t *testing.T) {
			cp, err := newConfigPoller(lggr, ethClient, mp, ocrAddress, &configStoreContractAddr)
			require.NoError(t, err)

			t.Run("when config has not been set, returns zero values", func(t *testing.T) {
				contractConfig, err := cp.LatestConfig(testutils.Context(t), 0)
				require.NoError(t, err)

				assert.Equal(t, ocrtypes.ConfigDigest{}, contractConfig.ConfigDigest)
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
			failingClient := new(evmClientMocks.Client)
			failingClient.On("ConfiguredChainID").Return(big.NewInt(42))
			failingClient.On("CallContract", mock.Anything, mock.MatchedBy(func(callArgs ethereum.CallMsg) bool {
				// initial call to retrieve config store address from aggregator
				return *callArgs.To == ocrAddress
			}), mock.Anything).Return(nil, errors.New("something exploded")).Once()
			cp, err := newConfigPoller(lggr, failingClient, mp, ocrAddress, &configStoreContractAddr)
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
				OnchainPublicKey:  utils.RandomAddress().Bytes(),
				TransmitAccount:   ocrtypes2.Account(utils.RandomAddress().Hex()),
				OffchainPublicKey: utils.RandomBytes32(),
				PeerID:            utils.MustNewPeerID(),
			},
			ConfigEncryptionPublicKey: utils.RandomBytes32(),
		})
	}
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
		50*time.Millisecond,
		50*time.Millisecond,
		50*time.Millisecond,
		50*time.Millisecond,
		50*time.Millisecond,
		1, // faults
		generateDefaultOCR2OnchainConfig(big.NewInt(0), big.NewInt(10)),
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

func generateDefaultOCR2OnchainConfig(minValue *big.Int, maxValue *big.Int) []byte {
	serializedConfig := make([]byte, 0)

	s1, err := bigbigendian.SerializeSigned(1, big.NewInt(1)) //version
	if err != nil {
		panic(err)
	}
	serializedConfig = append(serializedConfig, s1...)

	s2, err := bigbigendian.SerializeSigned(24, minValue) //min
	if err != nil {
		panic(err)
	}
	serializedConfig = append(serializedConfig, s2...)

	s3, err := bigbigendian.SerializeSigned(24, maxValue) //max
	if err != nil {
		panic(err)
	}
	serializedConfig = append(serializedConfig, s3...)

	return serializedConfig
}
