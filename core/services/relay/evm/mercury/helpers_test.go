package mercury

import (
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
	"github.com/stretchr/testify/require"

	evmclient "github.com/smartcontractkit/chainlink/v2/core/chains/evm/client"
	"github.com/smartcontractkit/chainlink/v2/core/chains/evm/logpoller"
	"github.com/smartcontractkit/chainlink/v2/core/gethwrappers/llo-feeds/generated/verifier"
	"github.com/smartcontractkit/chainlink/v2/core/gethwrappers/llo-feeds/generated/verifier_proxy"
	"github.com/smartcontractkit/chainlink/v2/core/internal/testutils"
	"github.com/smartcontractkit/chainlink/v2/core/internal/testutils/pgtest"
	"github.com/smartcontractkit/chainlink/v2/core/logger"
	pgmocks "github.com/smartcontractkit/chainlink/v2/core/services/pg/mocks"
	v1 "github.com/smartcontractkit/chainlink/v2/core/services/relay/evm/mercury/v1"
	"github.com/smartcontractkit/chainlink/v2/core/utils"
)

var (
	sampleFeedID       = [32]uint8{28, 145, 107, 74, 167, 229, 124, 167, 182, 138, 225, 191, 69, 101, 63, 86, 182, 86, 253, 58, 163, 53, 239, 127, 174, 105, 107, 102, 63, 27, 132, 114}
	sampleClientPubKey = hexutil.MustDecode("0x724ff6eae9e900270edfff233e16322a70ec06e1a6e62a81ef13921f398f6c93")
)

func buildSampleV1Report(p int64) []byte {
	feedID := sampleFeedID
	timestamp := uint32(42)
	bp := big.NewInt(p)
	bid := big.NewInt(243)
	ask := big.NewInt(244)
	currentBlockNumber := uint64(143)
	currentBlockHash := utils.NewHash()
	currentBlockTimestamp := uint64(123)
	validFromBlockNum := uint64(142)

	b, err := v1.ReportTypes.Pack(feedID, timestamp, bp, bid, ask, currentBlockNumber, currentBlockHash, currentBlockTimestamp, validFromBlockNum)
	if err != nil {
		panic(err)
	}
	return b
}

var sampleReports [][]byte

func init() {
	sampleReports = make([][]byte, 4)
	for i := 0; i < len(sampleReports); i++ {
		sampleReports[i] = buildSampleV1Report(int64(i))
	}
}

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
