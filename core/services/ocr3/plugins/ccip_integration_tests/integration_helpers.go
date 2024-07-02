package ccip_integration_tests

import (
	"context"
	"math/big"
	"testing"
	"time"

	"github.com/ethereum/go-ethereum/accounts/abi/bind/backends"
	"github.com/ethereum/go-ethereum/common"

	types2 "github.com/smartcontractkit/chainlink-common/pkg/types"

	"github.com/stretchr/testify/require"

	"github.com/smartcontractkit/chainlink/v2/core/chains/evm/client"
	"github.com/smartcontractkit/chainlink/v2/core/chains/evm/logpoller"
	"github.com/smartcontractkit/chainlink/v2/core/internal/testutils/pgtest"
	"github.com/smartcontractkit/chainlink/v2/core/logger"
	"github.com/smartcontractkit/chainlink/v2/core/services/relay/evm"
	evmtypes "github.com/smartcontractkit/chainlink/v2/core/services/relay/evm/types"
)

const chainID = 1337

type TestSetupData struct {
	LogPoller   logpoller.LogPoller
	ChainReader evm.ChainReaderService
}

func SetupReaderTestData(ctx context.Context, t *testing.T, simulatedBackend *backends.SimulatedBackend, address common.Address, chainReaderConfig evmtypes.ChainReaderConfig, contractName string) TestSetupData {
	lggr := logger.TestLogger(t)
	db := pgtest.NewSqlxDB(t)
	lpOpts := logpoller.Opts{
		PollPeriod:               time.Millisecond,
		FinalityDepth:            1,
		BackfillBatchSize:        1,
		RpcBatchSize:             1,
		KeepFinalizedBlocksDepth: 10000,
	}
	cl := client.NewSimulatedBackendClient(t, simulatedBackend, big.NewInt(chainID))
	lp := logpoller.NewLogPoller(logpoller.NewORM(big.NewInt(chainID), db, lggr), cl, lggr, lpOpts)
	require.NoError(t, lp.Start(ctx))

	cr, err := evm.NewChainReaderService(ctx, lggr, lp, cl, chainReaderConfig)
	require.NoError(t, err)
	err = cr.Bind(ctx, []types2.BoundContract{
		{
			Address: address.String(),
			Name:    contractName,
			Pending: false,
		},
	})
	require.NoError(t, err)
	require.NoError(t, cr.Start(ctx))
	for {
		if err := cr.Ready(); err == nil {
			break
		}
	}
	return TestSetupData{
		LogPoller:   lp,
		ChainReader: cr,
	}
}
