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
	"github.com/smartcontractkit/chainlink/v2/core/gethwrappers/llo-feeds/generated/verifier"
	"github.com/smartcontractkit/chainlink/v2/core/gethwrappers/llo-feeds/generated/verifier_proxy"
	"github.com/smartcontractkit/chainlink/v2/core/internal/testutils"
	"github.com/smartcontractkit/chainlink/v2/core/internal/testutils/pgtest"
	"github.com/smartcontractkit/chainlink/v2/core/logger"
	pgmocks "github.com/smartcontractkit/chainlink/v2/core/services/pg/mocks"
)

type TestHarness struct {
	configPoller     *ConfigPoller
	user             *bind.TransactOpts
	backend          *backends.SimulatedBackend
	verifierAddress  common.Address
	verifierContract *verifier.Verifier
	logPoller        logpoller.LogPoller
	eventBroadcaster *pgmocks.EventBroadcaster
	subscription     *pgmocks.Subscription
}

func SetupTH(t *testing.T, feedID common.Hash) TestHarness {
	key, err := crypto.GenerateKey()
	require.NoError(t, err)
	user, err := bind.NewKeyedTransactorWithChainID(key, big.NewInt(1337))
	require.NoError(t, err)
	b := backends.NewSimulatedBackend(core.GenesisAlloc{
		user.From: {Balance: big.NewInt(1000000000000000000)}},
		5*ethconfig.Defaults.Miner.GasCeil)

	proxyAddress, _, verifierProxy, err := verifier_proxy.DeployVerifierProxy(user, b, common.Address{})
	require.NoError(t, err, "failed to deploy test mercury verifier proxy contract")
	verifierAddress, _, verifierContract, err := verifier.DeployVerifier(user, b, proxyAddress)
	require.NoError(t, err, "failed to deploy test mercury verifier contract")
	_, err = verifierProxy.InitializeVerifier(user, verifierAddress)
	require.NoError(t, err)
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

	configPoller, err := NewConfigPoller(lggr, lp, verifierAddress, feedID, eventBroadcaster)
	require.NoError(t, err)

	configPoller.Start()

	return TestHarness{
		configPoller:     configPoller,
		user:             user,
		backend:          b,
		verifierAddress:  verifierAddress,
		verifierContract: verifierContract,
		logPoller:        lp,
		eventBroadcaster: eventBroadcaster,
		subscription:     subscription,
	}
}
