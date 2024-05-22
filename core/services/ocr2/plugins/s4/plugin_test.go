package s4_test

import (
	"errors"
	"testing"
	"time"

	"github.com/smartcontractkit/chainlink/v2/core/chains/evm/utils/big"
	"github.com/smartcontractkit/chainlink/v2/core/internal/cltest"
	"github.com/smartcontractkit/chainlink/v2/core/internal/testutils"
	"github.com/smartcontractkit/chainlink/v2/core/logger"
	"github.com/smartcontractkit/chainlink/v2/core/services/ocr2/plugins/s4"
	s4_svc "github.com/smartcontractkit/chainlink/v2/core/services/s4"
	s4_mocks "github.com/smartcontractkit/chainlink/v2/core/services/s4/mocks"

	commonlogger "github.com/smartcontractkit/chainlink-common/pkg/logger"

	"github.com/ethereum/go-ethereum/common"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
	"google.golang.org/protobuf/proto"

	"github.com/smartcontractkit/libocr/offchainreporting2plus/types"
)

func createPluginConfig(maxEntries uint) *s4.PluginConfig {
	return &s4.PluginConfig{
		MaxObservationEntries:   maxEntries,
		MaxReportEntries:        maxEntries,
		MaxDeleteExpiredEntries: maxEntries,
		NSnapshotShards:         1,
	}
}

func generateTestRows(t *testing.T, n int, ttl time.Duration) []*s4.Row {
	ormRows := generateTestOrmRows(t, n, ttl)
	rows := make([]*s4.Row, n)
	for i := 0; i < n; i++ {
		rows[i] = &s4.Row{
			Address:    ormRows[i].Address.Bytes(),
			Slotid:     uint32(ormRows[i].SlotId),
			Version:    ormRows[i].Version,
			Expiration: ormRows[i].Expiration,
			Payload:    ormRows[i].Payload,
			Signature:  ormRows[i].Signature,
		}
	}
	return rows
}

func generateTestOrmRow(t *testing.T, ttl time.Duration, version uint64, confimed bool) *s4_svc.Row {
	priv, addr := testutils.NewPrivateKeyAndAddress(t)
	row := &s4_svc.Row{
		Address:    big.New(addr.Big()),
		SlotId:     0,
		Version:    version,
		Confirmed:  confimed,
		Expiration: time.Now().Add(ttl).UnixMilli(),
		Payload:    cltest.MustRandomBytes(t, 64),
	}
	env := &s4_svc.Envelope{
		Address:    addr.Bytes(),
		SlotID:     row.SlotId,
		Version:    row.Version,
		Expiration: row.Expiration,
		Payload:    row.Payload,
	}
	sig, err := env.Sign(priv)
	assert.NoError(t, err)
	row.Signature = sig
	return row
}

func generateTestOrmRows(t *testing.T, n int, ttl time.Duration) []*s4_svc.Row {
	rows := make([]*s4_svc.Row, n)
	for i := 0; i < n; i++ {
		rows[i] = generateTestOrmRow(t, ttl, 0, false)
	}
	return rows
}

func generateConfirmedTestOrmRows(t *testing.T, n int, ttl time.Duration) []*s4_svc.Row {
	rows := make([]*s4_svc.Row, n)
	for i := 0; i < n; i++ {
		rows[i] = generateTestOrmRow(t, ttl, uint64(i), true)
	}
	return rows
}

func compareRows(t *testing.T, protoRows []*s4.Row, ormRows []*s4_svc.Row) {
	assert.Equal(t, len(ormRows), len(protoRows))
	for i, row := range protoRows {
		assert.Equal(t, row.Address, ormRows[i].Address.Bytes())
		assert.Equal(t, row.Version, ormRows[i].Version)
		assert.Equal(t, row.Expiration, ormRows[i].Expiration)
		assert.Equal(t, row.Payload, ormRows[i].Payload)
		assert.Equal(t, row.Signature, ormRows[i].Signature)
	}
}

func compareSnapshotRows(t *testing.T, snapshot []*s4.SnapshotRow, ormVersions []*s4_svc.SnapshotRow) {
	assert.Equal(t, len(ormVersions), len(snapshot))
	for i, row := range snapshot {
		assert.Equal(t, row.Address, ormVersions[i].Address.Bytes())
		assert.Equal(t, row.Slotid, uint32(ormVersions[i].SlotId))
		assert.Equal(t, row.Version, ormVersions[i].Version)
	}
}

func rowsToShapshotRows(rows []*s4_svc.Row) []*s4_svc.SnapshotRow {
	versions := make([]*s4_svc.SnapshotRow, len(rows))
	for i, r := range rows {
		versions[i] = &s4_svc.SnapshotRow{
			Address: r.Address,
			SlotId:  r.SlotId,
			Version: r.Version,
		}
	}
	return versions
}

func TestPlugin_NewReportingPlugin(t *testing.T) {
	t.Parallel()

	logger := commonlogger.NewOCRWrapper(logger.TestLogger(t), true, func(msg string) {})
	orm := s4_mocks.NewORM(t)

	t.Run("ErrInvalidIntervals", func(t *testing.T) {
		config := createPluginConfig(1)
		config.NSnapshotShards = 0

		_, err := s4.NewReportingPlugin(logger, config, orm)
		assert.ErrorIs(t, err, s4_svc.ErrInvalidIntervals)
	})

	t.Run("MaxObservationEntries is zero", func(t *testing.T) {
		config := createPluginConfig(1)
		config.MaxObservationEntries = 0

		_, err := s4.NewReportingPlugin(logger, config, orm)
		assert.ErrorContains(t, err, "max number of observation entries cannot be zero")
	})

	t.Run("MaxReportEntries is zero", func(t *testing.T) {
		config := createPluginConfig(1)
		config.MaxReportEntries = 0

		_, err := s4.NewReportingPlugin(logger, config, orm)
		assert.ErrorContains(t, err, "max number of report entries cannot be zero")
	})

	t.Run("MaxDeleteExpiredEntries is zero", func(t *testing.T) {
		config := createPluginConfig(1)
		config.MaxDeleteExpiredEntries = 0

		_, err := s4.NewReportingPlugin(logger, config, orm)
		assert.ErrorContains(t, err, "max number of delete expired entries cannot be zero")
	})

	t.Run("happy", func(t *testing.T) {
		config := createPluginConfig(1)
		p, err := s4.NewReportingPlugin(logger, config, orm)
		assert.NoError(t, err)
		assert.NotNil(t, p)
	})
}

func TestPlugin_Close(t *testing.T) {
	t.Parallel()

	logger := commonlogger.NewOCRWrapper(logger.TestLogger(t), true, func(msg string) {})
	config := createPluginConfig(10)
	orm := s4_mocks.NewORM(t)
	plugin, err := s4.NewReportingPlugin(logger, config, orm)
	assert.NoError(t, err)

	err = plugin.Close()
	assert.NoError(t, err)
}

func TestPlugin_ShouldTransmitAcceptedReport(t *testing.T) {
	t.Parallel()

	logger := commonlogger.NewOCRWrapper(logger.TestLogger(t), true, func(msg string) {})
	config := createPluginConfig(10)
	orm := s4_mocks.NewORM(t)
	plugin, err := s4.NewReportingPlugin(logger, config, orm)
	assert.NoError(t, err)

	should, err := plugin.ShouldTransmitAcceptedReport(testutils.Context(t), types.ReportTimestamp{}, nil)
	assert.NoError(t, err)
	assert.False(t, should)
}

func TestPlugin_ShouldAcceptFinalizedReport(t *testing.T) {
	t.Parallel()

	logger := commonlogger.NewOCRWrapper(logger.TestLogger(t), true, func(msg string) {})
	config := createPluginConfig(10)
	orm := s4_mocks.NewORM(t)
	plugin, err := s4.NewReportingPlugin(logger, config, orm)
	assert.NoError(t, err)

	t.Run("happy", func(t *testing.T) {
		ormRows := make([]*s4_svc.Row, 0)
		rows := generateTestRows(t, 10, time.Minute)
		orm.On("Update", mock.Anything, mock.Anything).Run(func(args mock.Arguments) {
			updateRow := args.Get(1).(*s4_svc.Row)
			ormRows = append(ormRows, updateRow)
		}).Return(nil).Times(10)

		report, err := proto.Marshal(&s4.Rows{
			Rows: rows,
		})
		assert.NoError(t, err)

		should, err := plugin.ShouldAcceptFinalizedReport(testutils.Context(t), types.ReportTimestamp{}, report)
		assert.NoError(t, err)
		assert.False(t, should)
		assert.Equal(t, 10, len(ormRows))
		compareRows(t, rows, ormRows)
	})

	t.Run("error", func(t *testing.T) {
		testErr := errors.New("some error")
		rows := generateTestRows(t, 1, time.Minute)
		orm.On("Update", mock.Anything, mock.Anything).Return(testErr).Once()

		report, err := proto.Marshal(&s4.Rows{
			Rows: rows,
		})
		assert.NoError(t, err)

		should, err := plugin.ShouldAcceptFinalizedReport(testutils.Context(t), types.ReportTimestamp{}, report)
		assert.NoError(t, err) // errors just logged
		assert.False(t, should)
	})

	t.Run("don't save expired", func(t *testing.T) {
		ormRows := make([]*s4_svc.Row, 0)
		rows := generateTestRows(t, 2, -time.Minute)

		report, err := proto.Marshal(&s4.Rows{
			Rows: rows,
		})
		assert.NoError(t, err)

		should, err := plugin.ShouldAcceptFinalizedReport(testutils.Context(t), types.ReportTimestamp{}, report)
		assert.NoError(t, err)
		assert.False(t, should)
		assert.Equal(t, 0, len(ormRows))
	})
}

func TestPlugin_Query(t *testing.T) {
	t.Parallel()

	logger := commonlogger.NewOCRWrapper(logger.TestLogger(t), true, func(msg string) {})
	config := createPluginConfig(10)
	orm := s4_mocks.NewORM(t)
	plugin, err := s4.NewReportingPlugin(logger, config, orm)
	assert.NoError(t, err)

	t.Run("happy", func(t *testing.T) {
		ormRows := generateTestOrmRows(t, 10, time.Minute)
		rows := rowsToShapshotRows(ormRows)

		orm.On("GetSnapshot", mock.Anything, mock.Anything).Return(rows, nil).Once()

		queryBytes, err := plugin.Query(testutils.Context(t), types.ReportTimestamp{})
		assert.NoError(t, err)

		query := &s4.Query{}
		err = proto.Unmarshal(queryBytes, query)
		assert.NoError(t, err)
		assert.Equal(t, s4_svc.MinAddress, s4.UnmarshalAddress(query.AddressRange.MinAddress))
		assert.Equal(t, s4_svc.MaxAddress, s4.UnmarshalAddress(query.AddressRange.MaxAddress))

		compareSnapshotRows(t, query.Rows, rows)
	})

	t.Run("empty", func(t *testing.T) {
		empty := make([]*s4_svc.SnapshotRow, 0)
		orm.On("GetSnapshot", mock.Anything, mock.Anything).Return(empty, nil).Once()

		query, err := plugin.Query(testutils.Context(t), types.ReportTimestamp{})
		assert.NoError(t, err)
		assert.NotNil(t, query)
	})

	t.Run("query with shards", func(t *testing.T) {
		config.NSnapshotShards = 16

		ormRows := generateTestOrmRows(t, 256, time.Minute)
		for i := 0; i < 256; i++ {
			var thisAddress common.Address
			thisAddress[0] = byte(i)
			ormRows[i].Address = big.New(thisAddress.Big())
		}
		versions := rowsToShapshotRows(ormRows)

		ar, err := s4_svc.NewInitialAddressRangeForIntervals(config.NSnapshotShards)
		assert.NoError(t, err)

		for i := 0; i <= int(config.NSnapshotShards); i++ {
			from := i * 16
			to := from + 16
			if i == int(config.NSnapshotShards) {
				from = 0
				to = 16
			}
			orm.On("GetSnapshot", mock.Anything, mock.Anything).Return(versions[from:to], nil).Once()

			query, err := plugin.Query(testutils.Context(t), types.ReportTimestamp{})
			assert.NoError(t, err)

			qq := &s4.Query{}
			err = proto.Unmarshal(query, qq)
			assert.NoError(t, err)

			assert.Len(t, qq.Rows, 16)
			for _, r := range qq.Rows {
				thisAddress := s4.UnmarshalAddress(r.Address)
				assert.True(t, ar.Contains(thisAddress))
			}

			ar.Advance()
		}
	})
}

func TestPlugin_Observation(t *testing.T) {
	t.Parallel()

	logger := commonlogger.NewOCRWrapper(logger.TestLogger(t), true, func(msg string) {})
	config := createPluginConfig(10)
	orm := s4_mocks.NewORM(t)
	plugin, err := s4.NewReportingPlugin(logger, config, orm)
	assert.NoError(t, err)

	t.Run("all unconfirmed", func(t *testing.T) {
		ormRows := generateTestOrmRows(t, int(config.MaxObservationEntries), time.Minute)
		for _, or := range ormRows {
			or.Confirmed = false
		}
		orm.On("DeleteExpired", mock.Anything, uint(10), mock.Anything, mock.Anything).Return(int64(10), nil).Once()
		orm.On("GetUnconfirmedRows", mock.Anything, config.MaxObservationEntries).Return(ormRows, nil).Once()

		observation, err := plugin.Observation(testutils.Context(t), types.ReportTimestamp{}, []byte{})
		assert.NoError(t, err)

		rows := &s4.Rows{}
		err = proto.Unmarshal(observation, rows)
		assert.NoError(t, err)
		assert.Len(t, rows.Rows, int(config.MaxObservationEntries))
	})

	t.Run("unconfirmed with query", func(t *testing.T) {
		numUnconfirmed := int(config.MaxObservationEntries / 2)
		ormRows := generateTestOrmRows(t, int(config.MaxObservationEntries), time.Minute)
		snapshot := make([]*s4_svc.SnapshotRow, len(ormRows))
		for i, or := range ormRows {
			or.Confirmed = i < numUnconfirmed // First half are confirmed
			or.Version = uint64(i)
			snapshot[i] = &s4_svc.SnapshotRow{
				Address:   or.Address,
				SlotId:    or.SlotId,
				Version:   or.Version,
				Confirmed: or.Confirmed,
			}
		}
		orm.On("DeleteExpired", mock.Anything, uint(10), mock.Anything, mock.Anything).Return(int64(10), nil).Once()
		orm.On("GetUnconfirmedRows", mock.Anything, config.MaxObservationEntries).Return(ormRows[numUnconfirmed:], nil).Once()
		orm.On("GetSnapshot", mock.Anything, mock.Anything).Return(snapshot, nil).Once()

		snapshotRows := rowsToShapshotRows(ormRows)
		query := &s4.Query{
			Rows: make([]*s4.SnapshotRow, len(snapshotRows)),
		}
		numHigherVersion := 2
		for i, v := range snapshotRows {
			query.Rows[i] = &s4.SnapshotRow{
				Address: v.Address.Bytes(),
				Slotid:  uint32(v.SlotId),
				Version: v.Version,
			}
			if i < numHigherVersion {
				ormRows[i].Version++
				snapshot[i].Version++
				orm.On("Get", mock.Anything, v.Address, v.SlotId).Return(ormRows[i], nil).Once()
			}
		}
		queryBytes, err := proto.Marshal(query)
		assert.NoError(t, err)

		observation, err := plugin.Observation(testutils.Context(t), types.ReportTimestamp{}, queryBytes)
		assert.NoError(t, err)

		rows := &s4.Rows{}
		err = proto.Unmarshal(observation, rows)
		assert.NoError(t, err)
		assert.Len(t, rows.Rows, numUnconfirmed+numHigherVersion)

		for i := 0; i < numUnconfirmed; i++ {
			assert.Equal(t, ormRows[numUnconfirmed+i].Version, rows.Rows[i].Version)
		}
		for i := 0; i < numHigherVersion; i++ {
			assert.Equal(t, ormRows[i].Version, rows.Rows[numUnconfirmed+i].Version)
		}
	})

	t.Run("missing from query", func(t *testing.T) {
		vLow, vHigh := uint64(2), uint64(5)
		ormRows := generateTestOrmRows(t, 3, time.Minute)
		// Follower node has 3 confirmed entries with latest versions.
		snapshot := make([]*s4_svc.SnapshotRow, len(ormRows))
		for i, or := range ormRows {
			or.Confirmed = true
			or.Version = vHigh
			snapshot[i] = &s4_svc.SnapshotRow{
				Address:   or.Address,
				SlotId:    or.SlotId,
				Version:   or.Version,
				Confirmed: or.Confirmed,
			}
		}

		// Query snapshot has:
		//   - First entry with same version
		//	 - Second entry with lower version
		//   - Third entry missing
		query := &s4.Query{
			Rows: []*s4.SnapshotRow{
				&s4.SnapshotRow{
					Address: snapshot[0].Address.Bytes(),
					Slotid:  uint32(snapshot[0].SlotId),
					Version: vHigh,
				},
				&s4.SnapshotRow{
					Address: snapshot[1].Address.Bytes(),
					Slotid:  uint32(snapshot[1].SlotId),
					Version: vLow,
				},
			},
		}
		queryBytes, err := proto.Marshal(query)
		assert.NoError(t, err)

		orm.On("DeleteExpired", mock.Anything, uint(10), mock.Anything, mock.Anything).Return(int64(10), nil).Once()
		orm.On("GetUnconfirmedRows", mock.Anything, config.MaxObservationEntries).Return([]*s4_svc.Row{}, nil).Once()
		orm.On("GetSnapshot", mock.Anything, mock.Anything).Return(snapshot, nil).Once()
		orm.On("Get", mock.Anything, snapshot[1].Address, snapshot[1].SlotId).Return(ormRows[1], nil).Once()
		orm.On("Get", mock.Anything, snapshot[2].Address, snapshot[2].SlotId).Return(ormRows[2], nil).Once()

		observation, err := plugin.Observation(testutils.Context(t), types.ReportTimestamp{}, queryBytes)
		assert.NoError(t, err)

		rows := &s4.Rows{}
		err = proto.Unmarshal(observation, rows)
		assert.NoError(t, err)
		assert.Len(t, rows.Rows, 2)
	})
}

func TestPlugin_Report(t *testing.T) {
	t.Parallel()

	logger := commonlogger.NewOCRWrapper(logger.TestLogger(t), true, func(msg string) {})
	config := createPluginConfig(10)
	orm := s4_mocks.NewORM(t)
	plugin, err := s4.NewReportingPlugin(logger, config, orm)
	assert.NoError(t, err)

	rows := generateTestRows(t, 10, time.Minute)
	observation, err := proto.Marshal(&s4.Rows{Rows: rows})
	assert.NoError(t, err)

	aos := []types.AttributedObservation{
		{
			Observation: observation,
		},
		{
			Observation: observation,
		},
	}
	ok, report, err := plugin.Report(testutils.Context(t), types.ReportTimestamp{}, nil, aos)
	assert.NoError(t, err)
	assert.True(t, ok)

	reportRows := &s4.Rows{}
	err = proto.Unmarshal(report, reportRows)
	assert.NoError(t, err)
	assert.Len(t, reportRows.Rows, 10)

	ok2, report2, err2 := plugin.Report(testutils.Context(t), types.ReportTimestamp{}, nil, aos)
	assert.NoError(t, err2)
	assert.True(t, ok2)

	reportRows2 := &s4.Rows{}
	err = proto.Unmarshal(report2, reportRows2)
	assert.NoError(t, err)

	// Verify that the same report was produced
	assert.Equal(t, reportRows, reportRows2)
}
