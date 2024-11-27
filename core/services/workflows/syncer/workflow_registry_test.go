package syncer

import (
	"context"
	"encoding/hex"
	"strconv"
	"testing"
	"time"

	"github.com/smartcontractkit/chainlink-common/pkg/services/servicetest"
	types "github.com/smartcontractkit/chainlink-common/pkg/types"
	query "github.com/smartcontractkit/chainlink-common/pkg/types/query"
	"github.com/smartcontractkit/chainlink-common/pkg/types/query/primitives"
	"github.com/smartcontractkit/chainlink-common/pkg/values"
	"github.com/smartcontractkit/chainlink/v2/core/internal/testutils"
	"github.com/smartcontractkit/chainlink/v2/core/internal/testutils/pgtest"
	"github.com/smartcontractkit/chainlink/v2/core/logger"
	"github.com/smartcontractkit/chainlink/v2/core/utils/crypto"
	"github.com/smartcontractkit/chainlink/v2/core/utils/matches"

	"github.com/stretchr/testify/require"
)

func Test_Workflow_Registry_Syncer(t *testing.T) {
	var (
		giveContents = "contents"
		wantContents = "updated contents"
		giveCfg      = ContractEventPollerConfig{
			ContractName:    ContractName,
			ContractAddress: "0xdeadbeef",
			StartBlockNum:   0,
			QueryCount:      20,
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
		gateway     = func(_ context.Context, _ string) ([]byte, error) {
			return []byte(wantContents), nil
		}
		ticker = make(chan time.Time)
		worker = NewWorkflowRegistry(lggr, orm, reader, gateway, giveCfg.ContractAddress, nil, nil, WithTicker(ticker))
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
			Name:    giveCfg.ContractName,
			Address: giveCfg.ContractAddress,
		},
		query.KeyFilter{
			Key: string(ForceUpdateSecretsEvent),
			Expressions: []query.Expression{
				query.Confidence(primitives.Finalized),
				query.Block(strconv.FormatUint(giveCfg.StartBlockNum, 10), primitives.Gte),
			},
		},
		query.LimitAndSort{
			SortBy: []query.SortBy{query.NewSortByTimestamp(query.Asc)},
			Limit:  query.Limit{Count: giveCfg.QueryCount},
		},
		new(values.Value),
	).Return([]types.Sequence{giveLog}, nil)

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
