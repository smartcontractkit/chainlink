package mercury

import (
	"database/sql"
	"fmt"
	"math/big"
	"testing"
	"time"

	"github.com/ethereum/go-ethereum/common"
	"github.com/ethereum/go-ethereum/common/hexutil"
	"github.com/onsi/gomega"
	"github.com/pkg/errors"
	confighelper2 "github.com/smartcontractkit/libocr/offchainreporting2plus/confighelper"
	"github.com/smartcontractkit/libocr/offchainreporting2plus/types"
	ocrtypes "github.com/smartcontractkit/libocr/offchainreporting2plus/types"
	ocrtypes2 "github.com/smartcontractkit/libocr/offchainreporting2plus/types"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
	"github.com/stretchr/testify/require"
	"github.com/umbracle/ethgo/abi"

	evmClientMocks "github.com/smartcontractkit/chainlink/v2/core/chains/evm/client/mocks"
	"github.com/smartcontractkit/chainlink/v2/core/chains/evm/logpoller/mocks"
	"github.com/smartcontractkit/chainlink/v2/core/internal/testutils"
	"github.com/smartcontractkit/chainlink/v2/core/services/pg"
	"github.com/smartcontractkit/chainlink/v2/core/utils"
)

func TestMercuryConfigPoller(t *testing.T) {
	feedID := utils.NewHash()
	feedIDBytes := [32]byte(feedID)

	th := SetupTH(t, feedID)
	th.configPoller.Start()
	t.Cleanup(func() {
		require.NoError(t, th.configPoller.Close())
	})
	th.subscription.On("Events").Return(nil)

	notify := th.configPoller.Notify()
	assert.Empty(t, notify)

	// Should have no config to begin with.
	_, config, err := th.configPoller.LatestConfigDetails(testutils.Context(t))
	require.NoError(t, err)
	require.Equal(t, ocrtypes2.ConfigDigest{}, config)

	// Create minimum number of nodes.
	n := 4
	var oracles []confighelper2.OracleIdentityExtra
	for i := 0; i < n; i++ {
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
	f := uint8(1)
	// Setup config on contract
	configType := abi.MustNewType("tuple()")
	onchainConfigVal, err := abi.Encode(map[string]interface{}{}, configType)
	require.NoError(t, err)
	signers, _, threshold, onchainConfig, offchainConfigVersion, offchainConfig, err := confighelper2.ContractSetConfigArgsForTests(
		2*time.Second,        // DeltaProgress
		20*time.Second,       // DeltaResend
		100*time.Millisecond, // DeltaRound
		0,                    // DeltaGrace
		1*time.Minute,        // DeltaStage
		100,                  // rMax
		[]int{len(oracles)},  // S
		oracles,
		[]byte{},             // reportingPluginConfig []byte,
		0,                    // Max duration query
		250*time.Millisecond, // Max duration observation
		250*time.Millisecond, // MaxDurationReport
		250*time.Millisecond, // MaxDurationShouldAcceptFinalizedReport
		250*time.Millisecond, // MaxDurationShouldTransmitAcceptedReport
		int(f),               // f
		onchainConfigVal,
	)
	require.NoError(t, err)
	signerAddresses, err := onchainPublicKeyToAddress(signers)
	require.NoError(t, err)
	offchainTransmitters := make([][32]byte, n)
	encodedTransmitter := make([]ocrtypes2.Account, n)
	for i := 0; i < n; i++ {
		offchainTransmitters[i] = oracles[i].OffchainPublicKey
		encodedTransmitter[i] = ocrtypes2.Account(fmt.Sprintf("%x", oracles[i].OffchainPublicKey[:]))
	}
	_, err = th.verifierContract.SetConfig(th.user, feedIDBytes, signerAddresses, offchainTransmitters, f, onchainConfig, offchainConfigVersion, offchainConfig)
	require.NoError(t, err, "failed to setConfig with feed ID")
	th.backend.Commit()

	latest, err := th.backend.BlockByNumber(testutils.Context(t), nil)
	require.NoError(t, err)
	// Ensure we capture this config set log.
	require.NoError(t, th.logPoller.Replay(testutils.Context(t), latest.Number().Int64()-1))

	// Send blocks until we see the config updated.
	var configBlock uint64
	var digest [32]byte
	gomega.NewGomegaWithT(t).Eventually(func() bool {
		th.backend.Commit()
		configBlock, digest, err = th.configPoller.LatestConfigDetails(testutils.Context(t))
		require.NoError(t, err)
		return ocrtypes2.ConfigDigest{} != digest
	}, testutils.WaitTimeout(t), 100*time.Millisecond).Should(gomega.BeTrue())

	// Assert the config returned is the one we configured.
	newConfig, err := th.configPoller.LatestConfig(testutils.Context(t), configBlock)
	require.NoError(t, err)
	// Note we don't check onchainConfig, as that is populated in the contract itself.
	assert.Equal(t, digest, [32]byte(newConfig.ConfigDigest))
	assert.Equal(t, signers, newConfig.Signers)
	assert.Equal(t, threshold, newConfig.F)
	assert.Equal(t, encodedTransmitter, newConfig.Transmitters)
	assert.Equal(t, offchainConfigVersion, newConfig.OffchainConfigVersion)
	assert.Equal(t, offchainConfig, newConfig.OffchainConfig)
}

func Test_MercuryConfigPoller_ConfigPersisted(t *testing.T) {
	feedID := utils.NewHash()
	feedIDBytes := [32]byte(feedID)
	th := SetupTH(t, feedID)

	offchainConfig := []byte{}

	t.Run("callIsConfigPersisted returns false", func(t *testing.T) {
		t.Run("when contract method missing, does not enable persistConfig", func(t *testing.T) {
			th.configPoller.addr = utils.ZeroAddress

			persistConfig, err := th.configPoller.callIsConfigPersisted(testutils.Context(t))
			require.NoError(t, err)
			assert.False(t, persistConfig)
		})
		t.Run("when contract method returns false, does not enable persistConfig", func(t *testing.T) {
			th.configPoller.addr = th.verifierContract.Address()

			persistConfig, err := th.configPoller.callIsConfigPersisted(testutils.Context(t))
			require.NoError(t, err)
			assert.False(t, persistConfig)
		})
	})

	t.Run("callIsConfigPersisted returns true", func(t *testing.T) {
		th.replaceVerifier(t, true)

		persistConfig, err := th.configPoller.callIsConfigPersisted(testutils.Context(t))
		require.NoError(t, err)
		assert.True(t, persistConfig)
	})

	t.Run("LatestConfigDetails, when logs have been pruned and persistConfig is true", func(t *testing.T) {
		// Give it a log poller that will never return logs
		mp := new(mocks.LogPoller)
		mp.On("RegisterFilter", mock.Anything).Return(nil)
		mp.On("LatestLogByEventSigWithConfs", mock.Anything, mock.Anything, mock.Anything, mock.Anything).Return(nil, sql.ErrNoRows)

		cp := th.configPoller
		cp.destChainLogPoller = mp

		t.Run("if callLatestConfigDetails succeeds", func(t *testing.T) {
			cp.persistConfig.Store(true)

			t.Run("when config has not been set, returns zero values", func(t *testing.T) {
				changedInBlock, configDigest, err := cp.LatestConfigDetails(testutils.Context(t))
				require.NoError(t, err)

				assert.Equal(t, 0, int(changedInBlock))
				assert.Equal(t, ocrtypes.ConfigDigest{}, configDigest)
			})
			t.Run("when config has been set, returns config details", func(t *testing.T) {
				th.setConfig(t, offchainConfig)
				th.backend.Commit()

				changedInBlock, configDigest, err := cp.LatestConfigDetails(testutils.Context(t))
				require.NoError(t, err)

				latest, err := th.backend.BlockByNumber(testutils.Context(t), nil)
				require.NoError(t, err)

				onchainDetails, err := th.verifierContract.LatestConfigDetails(nil, feedID)
				require.NoError(t, err)

				assert.Equal(t, latest.Number().Int64(), int64(changedInBlock))
				assert.Equal(t, onchainDetails.ConfigDigest, [32]byte(configDigest))
			})
		})
		t.Run("returns error if callLatestConfigDetails fails", func(t *testing.T) {
			failingClient := new(evmClientMocks.Client)
			failingClient.On("ConfiguredChainID").Return(big.NewInt(42))
			failingClient.On("CallContract", mock.Anything, mock.Anything, mock.Anything).Return(nil, errors.New("something exploded!"))

			th.configPoller.client = failingClient

			_, _, err := th.configPoller.LatestConfigDetails(testutils.Context(t))
			assert.EqualError(t, err, "something exploded!")

			failingClient.AssertExpectations(t)
		})
	})

	th.replaceVerifier(t, true)
	// ocrAddressPersistConfigEnabled, _, ocrContractPersistConfigEnabled, err = ocr2aggregator.DeployOCR2Aggregator(
	//     user,
	//     b,
	//     linkTokenAddress,
	//     big.NewInt(0),
	//     big.NewInt(10),
	//     accessAddress,
	//     accessAddress,
	//     9,
	//     "TEST",
	//     true,
	// )
	// require.NoError(t, err)
	th.backend.Commit()

	t.Run("LatestConfig, when logs have been pruned and persistConfig is true", func(t *testing.T) {
		// Give it a log poller that will never return logs
		mp := new(mocks.LogPoller)
		mp.On("RegisterFilter", mock.Anything).Return(nil)
		mp.On("Logs", mock.Anything, mock.Anything, mock.Anything, mock.Anything, mock.Anything).Return(nil, nil)

		latest, err := th.backend.BlockByNumber(testutils.Context(t), nil)
		require.NoError(t, err)
		blockNum := uint64(latest.Number().Int64())

		t.Run("if callLatestConfig succeeds", func(t *testing.T) {
			t.Run("when config has not been set, returns zero values", func(t *testing.T) {
				contractConfig, err := th.configPoller.LatestConfig(testutils.Context(t), blockNum)
				require.NoError(t, err)

				assert.Equal(t, ocrtypes.ConfigDigest{}, contractConfig.ConfigDigest)
			})
			t.Run("when config has been set, returns config details", func(t *testing.T) {
				contractConfig := th.setConfig(t, offchainConfig)
				th.backend.Commit()
				blockNum++

				newConfig, err := th.configPoller.LatestConfig(testutils.Context(t), blockNum)
				require.NoError(t, err)

				onchainDetails, err := th.verifierContract.LatestConfigDetails(nil, feedID)
				require.NoError(t, err)

				assert.Equal(t, onchainDetails.ConfigDigest, [32]byte(newConfig.ConfigDigest))
				assert.Equal(t, contractConfig.Signers, newConfig.Signers)
				assert.Equal(t, contractConfig.Transmitters, newConfig.Transmitters)
				assert.Equal(t, contractConfig.F, newConfig.F)
				assert.Equal(t, contractConfig.OffchainConfigVersion, newConfig.OffchainConfigVersion)
				assert.Equal(t, contractConfig.OffchainConfig, newConfig.OffchainConfig)
			})
		})
		t.Run("returns error if callLatestConfig fails", func(t *testing.T) {
			failingClient := new(evmClientMocks.Client)
			failingClient.On("ConfiguredChainID").Return(big.NewInt(42))
			failingClient.On("CallContract", mock.Anything, mock.Anything, mock.Anything).Return(nil, errors.New("something exploded!"))

			th.configPoller.client = failingClient

			_, err = th.configPoller.LatestConfig(testutils.Context(t), blockNum)
			assert.EqualError(t, err, "something exploded!")

			failingClient.AssertExpectations(t)
		})
	})
}

func TestNotify(t *testing.T) {
	feedIDStr := "8257737fdf4f79639585fd0ed01bea93c248a9ad940e98dd27f41c9b6230fed1"
	feedIDBytes, err := hexutil.Decode("0x" + feedIDStr)
	require.NoError(t, err)
	feedID := common.BytesToHash(feedIDBytes)

	eventCh := make(chan pg.Event)

	th := SetupTH(t, feedID)
	th.subscription.On("Events").Return((<-chan pg.Event)(eventCh))

	notify := th.configPoller.Notify()
	assert.Empty(t, notify)

	eventCh <- pg.Event{} // Empty event
	assert.Empty(t, notify)

	eventCh <- pg.Event{Payload: "address"} // missing topic values
	assert.Empty(t, notify)

	eventCh <- pg.Event{Payload: "address:val1"} // missing feedId topic value
	assert.Empty(t, notify)

	eventCh <- pg.Event{Payload: "address:8257737fdf4f79639585fd0ed01bea93c248a9ad940e98dd27f41c9b6230fed1,val2"} // wrong index
	assert.Empty(t, notify)

	eventCh <- pg.Event{Payload: "address:val1,val2,8257737fdf4f79639585fd0ed01bea93c248a9ad940e98dd27f41c9b6230fed1"} // wrong index
	assert.Empty(t, notify)

	eventCh <- pg.Event{Payload: "address:val1,0x8257737fdf4f79639585fd0ed01bea93c248a9ad940e98dd27f41c9b6230fed1"} // 0x prefix
	assert.Empty(t, notify)

	eventCh <- pg.Event{Payload: "address:val1,8257737fdf4f79639585fd0ed01bea93c248a9ad940e98dd27f41c9b6230fed1"}
	assert.Eventually(t, func() bool { <-notify; return true }, time.Second, 10*time.Millisecond)

	eventCh <- pg.Event{Payload: "address:val1,8257737fdf4f79639585fd0ed01bea93c248a9ad940e98dd27f41c9b6230fed1"} // try second time
	assert.Eventually(t, func() bool { <-notify; return true }, time.Second, 10*time.Millisecond)

	eventCh <- pg.Event{Payload: "address:val1,8257737fdf4f79639585fd0ed01bea93c248a9ad940e98dd27f41c9b6230fed1:additional"}
	assert.Eventually(t, func() bool { <-notify; return true }, time.Second, 10*time.Millisecond)
}

func onchainPublicKeyToAddress(publicKeys []types.OnchainPublicKey) (addresses []common.Address, err error) {
	for _, signer := range publicKeys {
		if len(signer) != 20 {
			return []common.Address{}, errors.Errorf("address is not 20 bytes %s", signer)
		}
		addresses = append(addresses, common.BytesToAddress(signer))
	}
	return addresses, nil
}
