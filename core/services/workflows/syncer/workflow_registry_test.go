package syncer

import (
	"context"
	"encoding/hex"
	"testing"
	"time"

	"github.com/stretchr/testify/mock"

	"github.com/jonboulle/clockwork"

	"github.com/smartcontractkit/chainlink-common/pkg/capabilities"
	"github.com/smartcontractkit/chainlink-common/pkg/custmsg"
	"github.com/smartcontractkit/chainlink-common/pkg/services/servicetest"
	types "github.com/smartcontractkit/chainlink-common/pkg/types"
	query "github.com/smartcontractkit/chainlink-common/pkg/types/query"
	"github.com/smartcontractkit/chainlink-common/pkg/types/query/primitives"
	"github.com/smartcontractkit/chainlink-common/pkg/values"
	"github.com/smartcontractkit/chainlink/v2/core/internal/testutils"
	"github.com/smartcontractkit/chainlink/v2/core/internal/testutils/pgtest"
	"github.com/smartcontractkit/chainlink/v2/core/logger"
	"github.com/smartcontractkit/chainlink/v2/core/services/keystore/keys/workflowkey"
	"github.com/smartcontractkit/chainlink/v2/core/utils/crypto"
	"github.com/smartcontractkit/chainlink/v2/core/utils/matches"

	"github.com/stretchr/testify/require"
)

type testDonNotifier struct {
	don capabilities.DON
	err error
}

func (t *testDonNotifier) WaitForDon(ctx context.Context) (capabilities.DON, error) {
	return t.don, t.err
}

func Test_Workflow_Registry_Syncer(t *testing.T) {
	var (
		giveContents    = "contents"
		wantContents    = "updated contents"
		contractAddress = "0xdeadbeef"
		giveCfg         = WorkflowEventPollerConfig{
			QueryCount: 20,
		}
		giveURL       = "http://example.com"
		giveHash, err = crypto.Keccak256([]byte(giveURL))

		giveLog = types.Sequence{
			Data: map[string]any{
				"SecretsURLHash": giveHash,
				"Owner":          "0xowneraddr",
			},
			Cursor: "cursor",
		}
	)

	require.NoError(t, err)

	var (
		lggr        = logger.TestLogger(t)
		db          = pgtest.NewSqlxDB(t)
		orm         = &orm{ds: db, lggr: lggr}
		ctx, cancel = context.WithCancel(testutils.Context(t))
		reader      = NewMockContractReader(t)
		emitter     = custmsg.NewLabeler()
		gateway     = func(_ context.Context, _ string) ([]byte, error) {
			return []byte(wantContents), nil
		}
		ticker = make(chan time.Time)

		handler = NewEventHandler(lggr, orm, gateway, nil, nil,
			emitter, clockwork.NewFakeClock(), workflowkey.Key{})
		loader = NewWorkflowRegistryContractLoader(lggr, contractAddress, func(ctx context.Context, bytes []byte) (ContractReader, error) {
			return reader, nil
		}, handler)

		worker = NewWorkflowRegistry(lggr, func(ctx context.Context, bytes []byte) (ContractReader, error) {
			return reader, nil
		}, contractAddress,
			WorkflowEventPollerConfig{
				QueryCount: 20,
			}, handler, loader,
			&testDonNotifier{
				don: capabilities.DON{
					ID: 1,
				},
				err: nil,
			},
			WithTicker(ticker))
	)

	// Cleanup the worker
	defer cancel()

	// Seed the DB with an original entry
	_, err = orm.Create(ctx, giveURL, hex.EncodeToString(giveHash), giveContents)
	require.NoError(t, err)

	// Mock out the contract reader query
	reader.EXPECT().QueryKey(
		matches.AnyContext,
		types.BoundContract{
			Name:    WorkflowRegistryContractName,
			Address: contractAddress,
		},
		query.KeyFilter{
			Key: string(ForceUpdateSecretsEvent),
			Expressions: []query.Expression{
				query.Confidence(primitives.Finalized),
				query.Block("0", primitives.Gte),
			},
		},
		query.LimitAndSort{
			SortBy: []query.SortBy{query.NewSortByTimestamp(query.Asc)},
			Limit:  query.Limit{Count: giveCfg.QueryCount},
		},
		new(values.Value),
	).Return([]types.Sequence{giveLog}, nil)
	reader.EXPECT().QueryKey(
		matches.AnyContext,
		types.BoundContract{
			Name:    WorkflowRegistryContractName,
			Address: contractAddress,
		},
		query.KeyFilter{
			Key: string(WorkflowPausedEvent),
			Expressions: []query.Expression{
				query.Confidence(primitives.Finalized),
				query.Block("0", primitives.Gte),
			},
		},
		query.LimitAndSort{
			SortBy: []query.SortBy{query.NewSortByTimestamp(query.Asc)},
			Limit:  query.Limit{Count: giveCfg.QueryCount},
		},
		new(values.Value),
	).Return([]types.Sequence{}, nil)
	reader.EXPECT().QueryKey(
		matches.AnyContext,
		types.BoundContract{
			Name:    WorkflowRegistryContractName,
			Address: contractAddress,
		},
		query.KeyFilter{
			Key: string(WorkflowDeletedEvent),
			Expressions: []query.Expression{
				query.Confidence(primitives.Finalized),
				query.Block("0", primitives.Gte),
			},
		},
		query.LimitAndSort{
			SortBy: []query.SortBy{query.NewSortByTimestamp(query.Asc)},
			Limit:  query.Limit{Count: giveCfg.QueryCount},
		},
		new(values.Value),
	).Return([]types.Sequence{}, nil)
	reader.EXPECT().QueryKey(
		matches.AnyContext,
		types.BoundContract{
			Name:    WorkflowRegistryContractName,
			Address: contractAddress,
		},
		query.KeyFilter{
			Key: string(WorkflowActivatedEvent),
			Expressions: []query.Expression{
				query.Confidence(primitives.Finalized),
				query.Block("0", primitives.Gte),
			},
		},
		query.LimitAndSort{
			SortBy: []query.SortBy{query.NewSortByTimestamp(query.Asc)},
			Limit:  query.Limit{Count: giveCfg.QueryCount},
		},
		new(values.Value),
	).Return([]types.Sequence{}, nil)
	reader.EXPECT().QueryKey(
		matches.AnyContext,
		types.BoundContract{
			Name:    WorkflowRegistryContractName,
			Address: contractAddress,
		},
		query.KeyFilter{
			Key: string(WorkflowUpdatedEvent),
			Expressions: []query.Expression{
				query.Confidence(primitives.Finalized),
				query.Block("0", primitives.Gte),
			},
		},
		query.LimitAndSort{
			SortBy: []query.SortBy{query.NewSortByTimestamp(query.Asc)},
			Limit:  query.Limit{Count: giveCfg.QueryCount},
		},
		new(values.Value),
	).Return([]types.Sequence{}, nil)
	reader.EXPECT().QueryKey(
		matches.AnyContext,
		types.BoundContract{
			Name:    WorkflowRegistryContractName,
			Address: contractAddress,
		},
		query.KeyFilter{
			Key: string(WorkflowRegisteredEvent),
			Expressions: []query.Expression{
				query.Confidence(primitives.Finalized),
				query.Block("0", primitives.Gte),
			},
		},
		query.LimitAndSort{
			SortBy: []query.SortBy{query.NewSortByTimestamp(query.Asc)},
			Limit:  query.Limit{Count: giveCfg.QueryCount},
		},
		new(values.Value),
	).Return([]types.Sequence{}, nil)
	reader.EXPECT().GetLatestValueWithHeadData(mock.Anything, mock.Anything, mock.Anything, mock.Anything, mock.Anything).Return(&types.Head{
		Height: "0",
	}, nil)
	reader.EXPECT().Start(mock.Anything).Return(nil)
	reader.EXPECT().Bind(mock.Anything, mock.Anything).Return(nil)

	// Go run the worker
	servicetest.Run(t, worker)

	// Send a tick to start a query
	ticker <- time.Now()

	// Require the secrets contents to eventually be updated
	require.Eventually(t, func() bool {
		secrets, err := orm.GetContents(ctx, giveURL)
		require.NoError(t, err)
		return secrets == wantContents
	}, 5*time.Second, time.Second)
}
