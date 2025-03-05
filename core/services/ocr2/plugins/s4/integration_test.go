package s4_test

import (
	"context"
	"crypto/ecdsa"
	"fmt"
	"maps"
	"math/rand"
	"testing"
	"time"

	"github.com/smartcontractkit/chainlink-integrations/evm/utils/big"
	"github.com/smartcontractkit/chainlink/v2/core/internal/cltest"
	"github.com/smartcontractkit/chainlink/v2/core/internal/testutils"
	"github.com/smartcontractkit/chainlink/v2/core/internal/testutils/pgtest"
	"github.com/smartcontractkit/chainlink/v2/core/logger"
	"github.com/smartcontractkit/chainlink/v2/core/services/ocr2/plugins/s4"
	s4_svc "github.com/smartcontractkit/chainlink/v2/core/services/s4"

	commonlogger "github.com/smartcontractkit/chainlink-common/pkg/logger"

	"github.com/smartcontractkit/libocr/commontypes"
	"github.com/smartcontractkit/libocr/offchainreporting2/types"

	"github.com/ethereum/go-ethereum/common"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
	"go.uber.org/multierr"
)

// Disclaimer: this is not a true integration test, it's more of a S4 feature test, on purpose.
// The purpose of the test is to make sure that S4 plugin works as expected in conjunction with Postgres ORM.
// Because of simplification, this emulates OCR2 rounds, not involving libocr.
// A proper integration test would be done per product, e.g. as a part of Functions integration test.

type don struct {
	size    int
	config  *s4.PluginConfig
	logger  logger.SugaredLogger
	orms    []s4_svc.ORM
	plugins []types.ReportingPlugin
}

func newDON(t *testing.T, size int, config *s4.PluginConfig) *don {
	t.Helper()

	logger := logger.TestLogger(t)
	db := pgtest.NewSqlxDB(t)

	orms := make([]s4_svc.ORM, size)
	plugins := make([]types.ReportingPlugin, size)

	for i := 0; i < size; i++ {
		ns := fmt.Sprintf("s4_int_test_%d", i)
		orm := s4_svc.NewPostgresORM(db, s4_svc.SharedTableName, ns)
		orms[i] = orm

		ocrLogger := commonlogger.NewOCRWrapper(logger, true, func(msg string) {})
		plugin, err := s4.NewReportingPlugin(ocrLogger, config, orm)
		require.NoError(t, err)
		plugins[i] = plugin
	}

	return &don{
		size:    size,
		config:  config,
		logger:  logger,
		orms:    orms,
		plugins: plugins,
	}
}

func (d *don) simulateOCR(ctx context.Context, rounds int) []error {
	errors := make([]error, d.size)

	for i := 0; i < rounds && ctx.Err() == nil; i++ {
		leaderIndex := i % d.size
		leader := d.plugins[leaderIndex]
		query, err := leader.Query(ctx, types.ReportTimestamp{})
		if err != nil {
			errors[leaderIndex] = multierr.Combine(errors[leaderIndex], err)
			continue
		}

		aos := make([]types.AttributedObservation, 0)
		for i := 0; i < d.size; i++ {
			observation, err2 := d.plugins[i].Observation(ctx, types.ReportTimestamp{}, query)
			if err2 != nil {
				errors[i] = multierr.Combine(errors[i], err2)
				continue
			}
			aos = append(aos, types.AttributedObservation{
				Observation: observation,
				Observer:    commontypes.OracleID(i),
			})
		}
		if len(aos) < d.size-1 {
			continue
		}

		_, report, err := leader.Report(ctx, types.ReportTimestamp{}, query, aos)
		if err != nil {
			errors[leaderIndex] = multierr.Combine(errors[leaderIndex], err)
			continue
		}

		for i := 0; i < d.size; i++ {
			_, err2 := d.plugins[i].ShouldAcceptFinalizedReport(ctx, types.ReportTimestamp{}, report)
			errors[i] = multierr.Combine(errors[i], err2)
		}
	}

	return errors
}

func compareSnapshots(s1, s2 []*s4_svc.SnapshotRow) bool {
	if len(s1) != len(s2) {
		return false
	}
	m1 := make(map[string]struct{}, len(s1))
	m2 := make(map[string]struct{}, len(s2))
	for i := 0; i < len(s1); i++ {
		k1 := fmt.Sprintf("%s_%d_%d", s1[i].Address.String(), s1[i].SlotId, s1[i].Version)
		k2 := fmt.Sprintf("%s_%d_%d", s2[i].Address.String(), s2[i].SlotId, s2[i].Version)
		m1[k1] = struct{}{}
		m2[k2] = struct{}{}
	}
	return maps.Equal(m1, m2)
}

func filter[T any](ss []T, test func(T) bool) (ret []T) {
	for _, s := range ss {
		if test(s) {
			ret = append(ret, s)
		}
	}
	return
}

func checkNoErrors(t *testing.T, errors []error) {
	t.Helper()

	for i, err := range errors {
		assert.NoError(t, err, "oracle %d", i)
	}
}

func checkNoUnconfirmedRows(ctx context.Context, t *testing.T, orm s4_svc.ORM, limit uint) {
	t.Helper()

	rows, err := orm.GetUnconfirmedRows(ctx, limit)
	assert.NoError(t, err)
	assert.Empty(t, rows)
}

func TestS4Integration_HappyDON(t *testing.T) {
	don := newDON(t, 4, createPluginConfig(100))
	ctx := testutils.Context(t)

	// injecting new records
	rows := generateTestOrmRows(t, 10, time.Minute)
	for _, row := range rows {
		err := don.orms[0].Update(ctx, row)
		require.NoError(t, err)
	}
	originSnapshot, err := don.orms[0].GetSnapshot(ctx, s4_svc.NewFullAddressRange())
	require.NoError(t, err)

	// S4 to propagate all records in one OCR round
	errors := don.simulateOCR(ctx, 1)
	checkNoErrors(t, errors)

	for i := 0; i < don.size; i++ {
		snapshot, err := don.orms[i].GetSnapshot(ctx, s4_svc.NewFullAddressRange())
		require.NoError(t, err)
		equal := compareSnapshots(originSnapshot, snapshot)
		assert.True(t, equal, "oracle %d", i)
		checkNoUnconfirmedRows(ctx, t, don.orms[i], 10)
	}
}

func TestS4Integration_HappyDON_4X(t *testing.T) {
	don := newDON(t, 4, createPluginConfig(100))
	ctx := testutils.Context(t)

	// injecting new records to all nodes
	for o := 0; o < don.size; o++ {
		rows := generateTestOrmRows(t, 10, time.Minute)
		for _, row := range rows {
			err := don.orms[o].Update(ctx, row)
			require.NoError(t, err)
		}
	}

	// S4 to propagate all records in one OCR round
	errors := don.simulateOCR(ctx, 1)
	checkNoErrors(t, errors)

	firstSnapshot, err := don.orms[0].GetSnapshot(ctx, s4_svc.NewFullAddressRange())
	require.NoError(t, err)

	for i := 1; i < don.size; i++ {
		snapshot, err := don.orms[i].GetSnapshot(ctx, s4_svc.NewFullAddressRange())
		require.NoError(t, err)
		equal := compareSnapshots(firstSnapshot, snapshot)
		assert.True(t, equal, "oracle %d", i)
		checkNoUnconfirmedRows(ctx, t, don.orms[i], 10)
	}
}

func TestS4Integration_WrongSignature(t *testing.T) {
	don := newDON(t, 4, createPluginConfig(100))
	ctx := testutils.Context(t)

	// injecting new records
	rows := generateTestOrmRows(t, 10, time.Minute)
	rows[0].Signature = rows[1].Signature
	for _, row := range rows {
		err := don.orms[0].Update(ctx, row)
		require.NoError(t, err)
	}
	originSnapshot, err := don.orms[0].GetSnapshot(ctx, s4_svc.NewFullAddressRange())
	require.NoError(t, err)
	originSnapshot = filter(originSnapshot, func(row *s4_svc.SnapshotRow) bool {
		return row.Address.Cmp(rows[0].Address) != 0 || row.SlotId != rows[0].SlotId
	})
	require.Len(t, originSnapshot, len(rows)-1)

	// S4 to propagate valid records in one OCR round
	errors := don.simulateOCR(ctx, 1)
	checkNoErrors(t, errors)

	for i := 1; i < don.size; i++ {
		snapshot, err2 := don.orms[i].GetSnapshot(ctx, s4_svc.NewFullAddressRange())
		require.NoError(t, err2)
		equal := compareSnapshots(originSnapshot, snapshot)
		assert.True(t, equal, "oracle %d", i)
	}

	// record with a wrong signature must remain unconfirmed
	ur, err := don.orms[0].GetUnconfirmedRows(ctx, 10)
	require.NoError(t, err)
	require.Len(t, ur, 1)
}

func TestS4Integration_MaxObservations(t *testing.T) {
	config := createPluginConfig(100)
	config.MaxObservationEntries = 5
	don := newDON(t, 4, config)
	ctx := testutils.Context(t)

	// injecting new records
	rows := generateTestOrmRows(t, 10, time.Minute)
	for _, row := range rows {
		err := don.orms[0].Update(ctx, row)
		require.NoError(t, err)
	}
	originSnapshot, err := don.orms[0].GetSnapshot(ctx, s4_svc.NewFullAddressRange())
	require.NoError(t, err)

	// It requires at least two rounds due to MaxObservationEntries = rows / 2
	errors := don.simulateOCR(ctx, 2)
	checkNoErrors(t, errors)

	for i := 1; i < don.size; i++ {
		snapshot, err := don.orms[i].GetSnapshot(ctx, s4_svc.NewFullAddressRange())
		require.NoError(t, err)
		equal := compareSnapshots(originSnapshot, snapshot)
		assert.True(t, equal, "oracle %d", i)
	}
}

func TestS4Integration_Expired(t *testing.T) {
	config := createPluginConfig(100)
	config.MaxObservationEntries = 5
	don := newDON(t, 4, config)
	ctx := testutils.Context(t)

	// injecting expiring records
	rows := generateTestOrmRows(t, 10, time.Millisecond)
	for _, row := range rows {
		err := don.orms[0].Update(ctx, row)
		require.NoError(t, err)
	}

	// within one round, all records will be GC-ed
	time.Sleep(testutils.TestInterval)
	errors := don.simulateOCR(ctx, 1)
	checkNoErrors(t, errors)

	for i := 0; i < don.size; i++ {
		snapshot, err := don.orms[i].GetSnapshot(ctx, s4_svc.NewFullAddressRange())
		require.NoError(t, err)
		require.Len(t, snapshot, 0)
	}
}

func TestS4Integration_NSnapshotShards(t *testing.T) {
	config := createPluginConfig(10000)
	config.NSnapshotShards = 4
	don := newDON(t, 4, config)
	ctx := testutils.Context(t)

	// injecting lots of new records (to be close to normal address distribution)
	rows := generateTestOrmRows(t, 1000, time.Minute)
	for _, row := range rows {
		err := don.orms[0].Update(ctx, row)
		require.NoError(t, err)
	}
	originSnapshot, err := don.orms[0].GetSnapshot(ctx, s4_svc.NewFullAddressRange())
	require.NoError(t, err)

	// this still requires one round, because Observation takes all unconfirmed rows
	errors := don.simulateOCR(ctx, 1)
	checkNoErrors(t, errors)

	for i := 1; i < don.size; i++ {
		snapshot, err := don.orms[i].GetSnapshot(ctx, s4_svc.NewFullAddressRange())
		require.NoError(t, err)
		equal := compareSnapshots(originSnapshot, snapshot)
		assert.True(t, equal, "oracle %d", i)
		checkNoUnconfirmedRows(ctx, t, don.orms[i], 1000)
	}
}

func TestS4Integration_OneNodeOutOfSync(t *testing.T) {
	don := newDON(t, 4, createPluginConfig(100))
	ctx := testutils.Context(t)

	// injecting same confirmed records to all nodes but the last one
	rows := generateConfirmedTestOrmRows(t, 10, time.Minute)
	for o := 0; o < don.size-1; o++ {
		for _, row := range rows {
			err := don.orms[o].Update(ctx, row)
			require.NoError(t, err)
		}
	}

	// all records will be propagated to the last node when it is a leader
	// leader selection is round-robin, so the 4th iteration picks the last node
	errors := don.simulateOCR(ctx, 4)
	checkNoErrors(t, errors)

	firstSnapshot, err := don.orms[0].GetSnapshot(ctx, s4_svc.NewFullAddressRange())
	require.NoError(t, err)
	lastSnapshot, err := don.orms[don.size-1].GetSnapshot(ctx, s4_svc.NewFullAddressRange())
	require.NoError(t, err)
	equal := compareSnapshots(firstSnapshot, lastSnapshot)
	assert.True(t, equal)
	checkNoUnconfirmedRows(ctx, t, don.orms[don.size-1], 10)
}

func TestS4Integration_RandomState(t *testing.T) {
	don := newDON(t, 4, createPluginConfig(1000))
	ctx := testutils.Context(t)

	type user struct {
		privateKey *ecdsa.PrivateKey
		address    *big.Big
	}

	nUsers := 100
	users := make([]user, nUsers)
	for i := 0; i < nUsers; i++ {
		pk, addr := testutils.NewPrivateKeyAndAddress(t)
		users[i] = user{pk, big.New(addr.Big())}
	}

	// generating test records
	for o := 0; o < don.size; o++ {
		for u := 0; u < nUsers; u++ {
			user := users[u]
			row := &s4_svc.Row{
				Address:    user.address,
				SlotId:     uint(u),
				Version:    uint64(rand.Intn(don.size)),
				Confirmed:  rand.Intn(2) == 0,
				Expiration: time.Now().UTC().Add(time.Minute).UnixMilli(),
				Payload:    cltest.MustRandomBytes(t, 64),
			}
			env := &s4_svc.Envelope{
				Address:    common.BytesToAddress(user.address.Bytes()).Bytes(),
				SlotID:     row.SlotId,
				Version:    row.Version,
				Expiration: row.Expiration,
				Payload:    row.Payload,
			}
			sig, err := env.Sign(user.privateKey)
			require.NoError(t, err)
			row.Signature = sig
			err = don.orms[o].Update(ctx, row)
			require.NoError(t, err)
		}
	}

	// for any state, all nodes should converge to the same snapshot
	errors := don.simulateOCR(ctx, 4)
	checkNoErrors(t, errors)

	firstSnapshot, err := don.orms[0].GetSnapshot(ctx, s4_svc.NewFullAddressRange())
	require.NoError(t, err)
	require.NotEmpty(t, firstSnapshot)
	checkNoUnconfirmedRows(ctx, t, don.orms[0], 1000)

	for i := 1; i < don.size; i++ {
		snapshot, err := don.orms[i].GetSnapshot(ctx, s4_svc.NewFullAddressRange())
		require.NoError(t, err)
		equal := compareSnapshots(firstSnapshot, snapshot)
		assert.True(t, equal, "oracle %d", i)
		checkNoUnconfirmedRows(ctx, t, don.orms[i], 1000)
	}
}
