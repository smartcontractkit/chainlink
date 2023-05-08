package mercury

import (
	"fmt"
	"math/big"
	"testing"
	"time"

	"github.com/ethereum/go-ethereum/accounts/abi/bind"
	"github.com/ethereum/go-ethereum/accounts/abi/bind/backends"
	"github.com/ethereum/go-ethereum/common"
	"github.com/ethereum/go-ethereum/common/hexutil"
	"github.com/ethereum/go-ethereum/core"
	"github.com/ethereum/go-ethereum/crypto"
	"github.com/ethereum/go-ethereum/eth/ethconfig"
	"github.com/onsi/gomega"
	"github.com/pkg/errors"
	confighelper2 "github.com/smartcontractkit/libocr/offchainreporting2/confighelper"
	"github.com/smartcontractkit/libocr/offchainreporting2/types"
	ocrtypes2 "github.com/smartcontractkit/libocr/offchainreporting2/types"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
	"github.com/umbracle/ethgo/abi"

	evmclient "github.com/smartcontractkit/chainlink/v2/core/chains/evm/client"
	"github.com/smartcontractkit/chainlink/v2/core/chains/evm/logpoller"
	"github.com/smartcontractkit/chainlink/v2/core/gethwrappers/generated/mercury_verifier"
	"github.com/smartcontractkit/chainlink/v2/core/gethwrappers/generated/mercury_verifier_proxy"
	"github.com/smartcontractkit/chainlink/v2/core/internal/testutils"
	"github.com/smartcontractkit/chainlink/v2/core/internal/testutils/pgtest"
	"github.com/smartcontractkit/chainlink/v2/core/logger"
	"github.com/smartcontractkit/chainlink/v2/core/services/pg"
	pgmocks "github.com/smartcontractkit/chainlink/v2/core/services/pg/mocks"
	"github.com/smartcontractkit/chainlink/v2/core/utils"
)

func TestMercuryConfigPoller(t *testing.T) {
	feedID := utils.NewHash()
	feedIDBytes := [32]byte(feedID)
	key, err := crypto.GenerateKey()
	require.NoError(t, err)
	user, err := bind.NewKeyedTransactorWithChainID(key, big.NewInt(1337))
	require.NoError(t, err)
	b := backends.NewSimulatedBackend(core.GenesisAlloc{
		user.From: {Balance: big.NewInt(1000000000000000000)}},
		5*ethconfig.Defaults.Miner.GasCeil)

	proxyAddress, _, _, err := mercury_verifier_proxy.DeployMercuryVerifierProxy(user, b, common.Address{})
	require.NoError(t, err, "failed to deploy test mercury verifier proxy contract")
	verifierAddress, _, verifierContract, err := mercury_verifier.DeployMercuryVerifier(user, b, proxyAddress)
	require.NoError(t, err, "failed to deploy test mercury verifier contract")
	b.Commit()

	db := pgtest.NewSqlxDB(t)
	cfg := pgtest.NewQConfig(false)
	ethClient := evmclient.NewSimulatedBackendClient(t, b, big.NewInt(1337))
	lggr := logger.TestLogger(t)
	ctx := testutils.Context(t)
	lorm := logpoller.NewORM(big.NewInt(1337), db, lggr, cfg)
	lp := logpoller.NewLogPoller(lorm, ethClient, lggr, 100*time.Millisecond, 1, 2, 2, 1000)
	eventBroadcaster := pgmocks.NewEventBroadcaster(t)
	subscription := pgmocks.NewSubscription(t)
	require.NoError(t, lp.Start(ctx))
	t.Cleanup(func() { lp.Close() })

	eventBroadcaster.On("Subscribe", "insert_on_evm_logs", "").Return(subscription, nil)
	subscription.On("Events").Return(nil)

	logPoller, err := NewConfigPoller(lggr, lp, verifierAddress, feedID, eventBroadcaster)
	require.NoError(t, err)

	notify := logPoller.Notify()
	assert.Empty(t, notify)

	// Should have no config to begin with.
	_, config, err := logPoller.LatestConfigDetails(testutils.Context(t))
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
	_, err = verifierContract.SetConfig(user, feedIDBytes, signerAddresses, offchainTransmitters, f, onchainConfig, offchainConfigVersion, offchainConfig)
	require.NoError(t, err, "failed to setConfig with feed ID")
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
	// Note we don't check onchainConfig, as that is populated in the contract itself.
	assert.Equal(t, digest, [32]byte(newConfig.ConfigDigest))
	assert.Equal(t, signers, newConfig.Signers)
	assert.Equal(t, threshold, newConfig.F)
	assert.Equal(t, encodedTransmitter, newConfig.Transmitters)
	assert.Equal(t, offchainConfigVersion, newConfig.OffchainConfigVersion)
	assert.Equal(t, offchainConfig, newConfig.OffchainConfig)
}

func TestNotify(t *testing.T) {
	feedIDStr := "8257737fdf4f79639585fd0ed01bea93c248a9ad940e98dd27f41c9b6230fed1"
	feedIDBytes, err := hexutil.Decode("0x" + feedIDStr)
	require.NoError(t, err)
	feedID := common.BytesToHash(feedIDBytes)

	key, err := crypto.GenerateKey()
	require.NoError(t, err)
	user, err := bind.NewKeyedTransactorWithChainID(key, big.NewInt(1337))
	require.NoError(t, err)
	b := backends.NewSimulatedBackend(core.GenesisAlloc{
		user.From: {Balance: big.NewInt(1000000000000000000)}},
		5*ethconfig.Defaults.Miner.GasCeil)

	db := pgtest.NewSqlxDB(t)
	cfg := pgtest.NewQConfig(false)
	ethClient := evmclient.NewSimulatedBackendClient(t, b, big.NewInt(1337))
	lggr := logger.TestLogger(t)
	ctx := testutils.Context(t)
	lorm := logpoller.NewORM(big.NewInt(1337), db, lggr, cfg)
	lp := logpoller.NewLogPoller(lorm, ethClient, lggr, 100*time.Millisecond, 1, 2, 2, 1000)
	eventBroadcaster := pgmocks.NewEventBroadcaster(t)
	subscription := pgmocks.NewSubscription(t)
	require.NoError(t, lp.Start(ctx))
	t.Cleanup(func() { lp.Close() })

	eventCh := make(chan pg.Event)

	eventBroadcaster.On("Subscribe", "insert_on_evm_logs", "").Return(subscription, nil)
	subscription.On("Events").Return((<-chan pg.Event)(eventCh))

	logPoller, err := NewConfigPoller(lggr, lp, common.Address{1}, feedID, eventBroadcaster)
	require.NoError(t, err)

	notify := logPoller.Notify()
	assert.Empty(t, notify)

	eventCh <- pg.Event{} // No topic values
	assert.Empty(t, notify)

	eventCh <- pg.Event{Payload: "val1"} // missing feedId topic value
	assert.Empty(t, notify)

	eventCh <- pg.Event{Payload: "8257737fdf4f79639585fd0ed01bea93c248a9ad940e98dd27f41c9b6230fed1,val2"} // wrong index
	assert.Empty(t, notify)

	eventCh <- pg.Event{Payload: "val1,val2,8257737fdf4f79639585fd0ed01bea93c248a9ad940e98dd27f41c9b6230fed1"} // wrong index
	assert.Empty(t, notify)

	eventCh <- pg.Event{Payload: "val1,0x8257737fdf4f79639585fd0ed01bea93c248a9ad940e98dd27f41c9b6230fed1"} // 0x prefix
	assert.Empty(t, notify)

	eventCh <- pg.Event{Payload: "val1,8257737fdf4f79639585fd0ed01bea93c248a9ad940e98dd27f41c9b6230fed1"}
	assert.Eventually(t, func() bool { <-notify; return true }, time.Second, 10*time.Millisecond)

	eventCh <- pg.Event{Payload: "val1,8257737fdf4f79639585fd0ed01bea93c248a9ad940e98dd27f41c9b6230fed1"}
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
