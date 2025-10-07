package testutils

import (
	"context"
	"encoding/json"
	"math/big"
	"testing"
	"time"

	"github.com/ethereum/go-ethereum/accounts/abi/bind"
	"github.com/ethereum/go-ethereum/common"
	"github.com/ethereum/go-ethereum/core"
	ethtypes "github.com/ethereum/go-ethereum/core/types"
	"github.com/ethereum/go-ethereum/ethclient/simulated"
	"github.com/stretchr/testify/require"
	"go.uber.org/atomic"

	"github.com/smartcontractkit/chainlink-common/pkg/logger"
	"github.com/smartcontractkit/chainlink-common/pkg/types"
	"github.com/smartcontractkit/chainlink-evm/pkg/assets"
	evmclient "github.com/smartcontractkit/chainlink-evm/pkg/client"
	"github.com/smartcontractkit/chainlink-evm/pkg/config"
	"github.com/smartcontractkit/chainlink-evm/pkg/heads/headstest"
	"github.com/smartcontractkit/chainlink-evm/pkg/keys"
	"github.com/smartcontractkit/chainlink-evm/pkg/keys/keystest"
	"github.com/smartcontractkit/chainlink-evm/pkg/logpoller"
	"github.com/smartcontractkit/chainlink-evm/pkg/testutils"
	evmtypes "github.com/smartcontractkit/chainlink-evm/pkg/types"
	"github.com/smartcontractkit/chainlink/v2/core/services/relay/evm"
)

// Test harness with EVM backend and chainlink core services like
// Log Poller and Head Tracker
type EVMBackendTH struct {
	// Backend details
	Lggr      logger.Logger
	ChainID   *big.Int
	Backend   evmtypes.Backend
	EVMClient evmclient.Client

	ContractsOwner     *bind.TransactOpts
	ContractsOwnerSign func(bytes []byte) ([]byte, error)

	HeadTracker logpoller.HeadTracker
	LogPoller   logpoller.LogPoller
}

var startID = atomic.NewInt64(1000)

// Test harness to create a simulated backend for testing a LOOPCapability
func NewEVMBackendTH(t *testing.T) *EVMBackendTH {
	lggr := logger.Test(t)

	memKS := keystest.NewMemoryChainStore()
	ownerAddress := memKS.MustCreate(t)
	chainStore := keys.NewChainStore(memKS, testutils.SimulatedChainID)

	contractsOwner := &bind.TransactOpts{
		From: ownerAddress,
		Signer: func(addr common.Address, tx *ethtypes.Transaction) (*ethtypes.Transaction, error) {
			return chainStore.SignTx(testutils.Context(t), addr, tx)
		},
	}

	// Setup simulated go-ethereum EVM backend
	genesisData := core.GenesisAlloc{
		ownerAddress: {Balance: assets.Ether(100000).ToInt()},
	}

	chainID := big.NewInt(startID.Add(1))
	backend := simulated.NewBackend(genesisData)

	h, err := backend.Client().HeaderByNumber(testutils.Context(t), nil)
	require.NoError(t, err)
	//nolint:gosec // G115
	blockTime := time.UnixMilli(int64(h.Time))
	err = backend.AdjustTime(time.Since(blockTime) - 24*time.Hour)
	require.NoError(t, err)
	backend.Commit()

	// Setup backend client
	client := evmclient.NewSimulatedBackendClient(t, backend, chainID)

	th := &EVMBackendTH{
		Lggr:      lggr,
		ChainID:   chainID,
		Backend:   backend,
		EVMClient: client,

		ContractsOwner: contractsOwner,
		ContractsOwnerSign: func(bytes []byte) ([]byte, error) {
			return memKS.Sign(testutils.Context(t), ownerAddress.String(), bytes)
		},
	}
	th.HeadTracker, th.LogPoller = th.SetupCoreServices(t)

	return th
}

// Setup core services like log poller and head tracker for the simulated backend
func (th *EVMBackendTH) SetupCoreServices(t *testing.T) (logpoller.HeadTracker, logpoller.LogPoller) {
	db := testutils.NewSqlxDB(t)
	const finalityDepth = 2
	ht := headstest.NewSimulatedHeadTracker(th.EVMClient, false, finalityDepth)
	lp := logpoller.NewLogPoller(
		logpoller.NewORM(th.EVMClient.ConfiguredChainID(), db, th.Lggr),
		th.EVMClient,
		th.Lggr,
		ht,
		logpoller.Opts{
			PollPeriod:               100 * time.Millisecond,
			FinalityDepth:            finalityDepth,
			BackfillBatchSize:        3,
			RPCBatchSize:             2,
			KeepFinalizedBlocksDepth: 1000,
		},
	)
	require.NoError(t, ht.Start(testutils.Context(t)))
	require.NoError(t, lp.Start(testutils.Context(t)))
	t.Cleanup(func() { ht.Close() })
	t.Cleanup(func() { lp.Close() })
	// Sleep 200ms to allow LP to load filters
	time.Sleep(time.Millisecond * 200)
	return ht, lp
}

func (th *EVMBackendTH) NewContractReader(ctx context.Context, t *testing.T, cfg []byte) (types.ContractReader, error) {
	crCfg := &config.ChainReaderConfig{}
	if err := json.Unmarshal(cfg, crCfg); err != nil {
		return nil, err
	}

	svc, err := evm.NewChainReaderService(ctx, th.Lggr, th.LogPoller, th.HeadTracker, th.EVMClient, *crCfg)
	if err != nil {
		return nil, err
	}

	return svc, err
}
