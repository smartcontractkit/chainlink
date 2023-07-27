package mercury

import (
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
	configPoller     *configPoller
	user             *bind.TransactOpts
	backend          *backends.SimulatedBackend
	verifierContract *mercury_verifier.MercuryVerifier
	logPoller        logpoller.LogPoller
	eventBroadcaster *pgmocks.EventBroadcaster
	subscription     *pgmocks.Subscription
	logger           logger.Logger
}

func (th TestHarness) replaceVerifier(t *testing.T, persistConfig bool) *mercury_verifier.MercuryVerifier {
	th.verifierContract = deployVerifier(t, persistConfig)
	th.configPoller.addr = th.verifierContract.Address()
	return th.verifierContract
}

func (th TestHarness) setConfig(t *testing.T, offchainConfig []byte) ocrtypes2.ContractConfig {
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
		nil,
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

func SetupTH(t *testing.T, feedID common.Hash) TestHarness {
	key, err := crypto.GenerateKey()
	require.NoError(t, err)
	user, err := bind.NewKeyedTransactorWithChainID(key, big.NewInt(1337))
	require.NoError(t, err)
	b := backends.NewSimulatedBackend(core.GenesisAlloc{
		th.configPoller.From: {Balance: big.NewInt(1000000000000000000)}},
		5*ethconfig.Defaults.Miner.GasCeil)

	verifierContract := deployVerifier(t, false)
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

	configPoller, err := NewConfigPoller(lggr, ethClient, lp, verifierAddress, feedID, eventBroadcaster)
	require.NoError(t, err)

	return TestHarness{
		configPoller:     configPoller,
		user:             user,
		backend:          b,
		verifierContract: verifierContract,
		logPoller:        lp,
		eventBroadcaster: eventBroadcaster,
		subscription:     subscription,
		logger:           lggr,
	}
}

func deployVerifier(t *testing.T, persistConfig bool) (verifierContract *mercury_verifier.MercuryVerifier) {
	proxyAddress, _, verifierProxy, err := mercury_verifier_proxy.DeployMercuryVerifierProxy(user, b, common.Address{})
	require.NoError(t, err, "failed to deploy test mercury verifier proxy contract")
	verifierAddress, _, verifierContract, err := mercury_verifier.DeployMercuryVerifier(user, b, proxyAddress, persistConfig)
	require.NoError(t, err, "failed to deploy test mercury verifier contract")
	_, err = verifierProxy.InitializeVerifier(user, verifierAddress)
	require.NoError(t, err)
	return verifierContract
}
