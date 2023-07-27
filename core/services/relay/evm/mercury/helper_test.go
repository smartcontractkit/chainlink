package mercury

import (
	"fmt"
	"math/big"
	"testing"
	"time"

	"github.com/ethereum/go-ethereum/accounts/abi/bind"
	"github.com/ethereum/go-ethereum/accounts/abi/bind/backends"
	"github.com/ethereum/go-ethereum/common"
	"github.com/ethereum/go-ethereum/core"
	"github.com/ethereum/go-ethereum/crypto"
	"github.com/ethereum/go-ethereum/eth/ethconfig"
	"github.com/stretchr/testify/require"
	"github.com/umbracle/ethgo/abi"

	evmclient "github.com/smartcontractkit/chainlink/v2/core/chains/evm/client"
	"github.com/smartcontractkit/chainlink/v2/core/chains/evm/logpoller"
	"github.com/smartcontractkit/chainlink/v2/core/gethwrappers/generated/mercury_verifier"
	"github.com/smartcontractkit/chainlink/v2/core/gethwrappers/generated/mercury_verifier_proxy"
	"github.com/smartcontractkit/chainlink/v2/core/internal/testutils"
	"github.com/smartcontractkit/chainlink/v2/core/internal/testutils/pgtest"
	"github.com/smartcontractkit/chainlink/v2/core/logger"
	pgmocks "github.com/smartcontractkit/chainlink/v2/core/services/pg/mocks"
	"github.com/smartcontractkit/chainlink/v2/core/utils"
	confighelper2 "github.com/smartcontractkit/libocr/offchainreporting2plus/confighelper"
	ocrtypes2 "github.com/smartcontractkit/libocr/offchainreporting2plus/types"
)

type TestHarness struct {
	configPoller         *configPoller
	user                 *bind.TransactOpts
	backend              *backends.SimulatedBackend
	verifierContract     *mercury_verifier.MercuryVerifier
	verifierContractAddr common.Address
	logPoller            logpoller.LogPoller
	eventBroadcaster     *pgmocks.EventBroadcaster
	subscription         *pgmocks.Subscription
	logger               logger.Logger
}

func (th TestHarness) replaceVerifier(t *testing.T, persistConfig bool) *mercury_verifier.MercuryVerifier {
	th.verifierContract, th.verifierContractAddr = deployVerifier(t, th.user, th.backend, persistConfig)
	th.configPoller.addr = th.verifierContract.Address()
	return th.verifierContract
}

func (th TestHarness) setConfig(t *testing.T, feedID [32]byte) ocrtypes2.ContractConfig {
	offchainConfig := []byte{}

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

	signers, _, f, onchainConfig, offchainConfigVersion, offchainConfig, err := confighelper2.ContractSetConfigArgsForTests(
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
	encodedTransmitters := make([]ocrtypes2.Account, n)
	for i := 0; i < n; i++ {
		offchainTransmitters[i] = oracles[i].OffchainPublicKey
		encodedTransmitters[i] = ocrtypes2.Account(fmt.Sprintf("%x", oracles[i].OffchainPublicKey[:]))
	}

	_, err = th.verifierContract.SetConfig(th.user, feedID, signerAddresses, offchainTransmitters, f, onchainConfig, offchainConfigVersion, offchainConfig)
	require.NoError(t, err, "failed to setConfig with feed ID")
	return ocrtypes2.ContractConfig{
		Signers:               signers,
		Transmitters:          encodedTransmitters,
		F:                     f,
		OnchainConfig:         onchainConfig,
		OffchainConfigVersion: offchainConfigVersion,
		OffchainConfig:        offchainConfig,
	}
}

func SetupTH(t *testing.T, feedID common.Hash) TestHarness {
	key, err := crypto.GenerateKey()
	require.NoError(t, err)
	user, err := bind.NewKeyedTransactorWithChainID(key, big.NewInt(1337))
	require.NoError(t, err)
	b := backends.NewSimulatedBackend(core.GenesisAlloc{
		user.From: {Balance: big.NewInt(1000000000000000000)}},
		5*ethconfig.Defaults.Miner.GasCeil)

	verifierContract, addr := deployVerifier(t, user, b, false)
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

	configPoller, err := NewConfigPoller(lggr, ethClient, lp, addr, feedID, eventBroadcaster)
	require.NoError(t, err)

	return TestHarness{
		configPoller:         configPoller,
		user:                 user,
		backend:              b,
		verifierContract:     verifierContract,
		verifierContractAddr: addr,
		logPoller:            lp,
		eventBroadcaster:     eventBroadcaster,
		subscription:         subscription,
		logger:               lggr,
	}
}

func deployVerifier(t *testing.T, user *bind.TransactOpts, b *backends.SimulatedBackend, persistConfig bool) (verifierContract *mercury_verifier.MercuryVerifier, verifierAddress common.Address) {
	proxyAddress, _, verifierProxy, err := mercury_verifier_proxy.DeployMercuryVerifierProxy(user, b, common.Address{})
	require.NoError(t, err, "failed to deploy test mercury verifier proxy contract")
	verifierAddress, _, verifierContract, err = mercury_verifier.DeployMercuryVerifier(user, b, proxyAddress, persistConfig)
	require.NoError(t, err, "failed to deploy test mercury verifier contract")
	_, err = verifierProxy.InitializeVerifier(user, verifierAddress)
	require.NoError(t, err)
	return verifierContract, verifierAddress
}
