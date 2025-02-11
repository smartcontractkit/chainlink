package testutils

import (
	"context"
	"encoding/json"
	"math/big"
	"testing"
	"time"

	"github.com/ethereum/go-ethereum/accounts/abi/bind"
	"github.com/ethereum/go-ethereum/core"
	"github.com/ethereum/go-ethereum/eth/ethconfig"
	"github.com/stretchr/testify/require"

	"github.com/smartcontractkit/chainlink-common/pkg/logger"
	"github.com/smartcontractkit/chainlink-common/pkg/types"
	"github.com/smartcontractkit/chainlink-integrations/evm/assets"
	evmclient "github.com/smartcontractkit/chainlink-integrations/evm/client"
	evmtypes "github.com/smartcontractkit/chainlink-integrations/evm/types"
	"github.com/smartcontractkit/chainlink/v2/core/chains/evm/headtracker"
	"github.com/smartcontractkit/chainlink/v2/core/chains/evm/logpoller"
	"github.com/smartcontractkit/chainlink/v2/core/internal/cltest"
	"github.com/smartcontractkit/chainlink/v2/core/internal/testutils"
	"github.com/smartcontractkit/chainlink/v2/core/internal/testutils/pgtest"
	"github.com/smartcontractkit/chainlink/v2/core/services/keystore/keys/ethkey"
	"github.com/smartcontractkit/chainlink/v2/core/services/relay/evm"
	evmrelaytypes "github.com/smartcontractkit/chainlink/v2/core/services/relay/evm/types"
)

// Test harness with EVM backend and chainlink core services like
// Log Poller and Head Tracker
type EVMBackendTH struct {
	// Backend details
	Lggr      logger.Logger
	ChainID   *big.Int
	Backend   evmtypes.Backend
	EVMClient evmclient.Client

	ContractsOwner    *bind.TransactOpts
	ContractsOwnerKey ethkey.KeyV2

	HeadTracker logpoller.HeadTracker
	LogPoller   logpoller.LogPoller
}

// Test harness to create a simulated backend for testing a LOOPCapability
func NewEVMBackendTH(t *testing.T) *EVMBackendTH {
	lggr := logger.Test(t)

	ownerKey := cltest.MustGenerateRandomKey(t)
	contractsOwner, err := bind.NewKeyedTransactorWithChainID(ownerKey.ToEcdsaPrivKey(), testutils.SimulatedChainID)
	require.NoError(t, err)

	// Setup simulated go-ethereum EVM backend
	genesisData := core.GenesisAlloc{
		contractsOwner.From: {Balance: assets.Ether(100000).ToInt()},
	}
	chainID := testutils.SimulatedChainID
	backend := cltest.NewSimulatedBackend(t, genesisData, ethconfig.Defaults.Miner.GasCeil)
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

		ContractsOwner:    contractsOwner,
		ContractsOwnerKey: ownerKey,
	}
	th.HeadTracker, th.LogPoller = th.SetupCoreServices(t)

	return th
}

// Setup core services like log poller and head tracker for the simulated backend
func (th *EVMBackendTH) SetupCoreServices(t *testing.T) (logpoller.HeadTracker, logpoller.LogPoller) {
	db := pgtest.NewSqlxDB(t)
	const finalityDepth = 2
	ht := headtracker.NewSimulatedHeadTracker(th.EVMClient, false, finalityDepth)
	lp := logpoller.NewLogPoller(
		logpoller.NewORM(testutils.SimulatedChainID, db, th.Lggr),
		th.EVMClient,
		th.Lggr,
		ht,
		logpoller.Opts{
			PollPeriod:               100 * time.Millisecond,
			FinalityDepth:            finalityDepth,
			BackfillBatchSize:        3,
			RpcBatchSize:             2,
			KeepFinalizedBlocksDepth: 1000,
		},
	)
	require.NoError(t, ht.Start(testutils.Context(t)))
	require.NoError(t, lp.Start(testutils.Context(t)))
	t.Cleanup(func() { ht.Close() })
	t.Cleanup(func() { lp.Close() })
	return ht, lp
}

func (th *EVMBackendTH) NewContractReader(ctx context.Context, t *testing.T, cfg []byte) (types.ContractReader, error) {
	crCfg := &evmrelaytypes.ChainReaderConfig{}
	if err := json.Unmarshal(cfg, crCfg); err != nil {
		return nil, err
	}

	svc, err := evm.NewChainReaderService(ctx, th.Lggr, th.LogPoller, th.HeadTracker, th.EVMClient, *crCfg)
	if err != nil {
		return nil, err
	}

	return svc, err
}
