package s4_test

import (
	"crypto/ecdsa"
	"crypto/rand"
	"errors"
	"math/big"
	"testing"
	"time"

	"github.com/smartcontractkit/chainlink/v2/core/internal/testutils"
	"github.com/smartcontractkit/chainlink/v2/core/services/ocr2/plugins/s4"
	s4_orm "github.com/smartcontractkit/chainlink/v2/core/services/s4"
	s4_mocks "github.com/smartcontractkit/chainlink/v2/core/services/s4/mocks"

	"github.com/ethereum/go-ethereum/common"
	"github.com/ethereum/go-ethereum/crypto"
	"github.com/smartcontractkit/libocr/offchainreporting2/types"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
	"google.golang.org/protobuf/proto"
)

func createPluginConfig(max uint) *s4.PluginConfig {
	return &s4.PluginConfig{
		MaxObservationEntries: max,
		MaxReportEntries:      max,
		NSnapshotShards:       1,
	}
}

func mustRandomBytes(t *testing.T, n int) []byte {
	b := make([]byte, n)
	k, err := rand.Read(b)
	assert.NoError(t, err)
	assert.Equal(t, n, k)
	return b
}

func mustStringToBig(t *testing.T, bn string) *big.Int {
	assert.NotEmpty(t, bn)
	bigAddress := new(big.Int)
	_, ok := bigAddress.SetString(bn, 10)
	assert.True(t, ok)
	return bigAddress
}

func generateCryptoEntity(t *testing.T) (*ecdsa.PrivateKey, *ecdsa.PublicKey, common.Address) {
	privateKey, err := crypto.GenerateKey()
	assert.NoError(t, err)

	publicKey := privateKey.Public()
	publicKeyECDSA, ok := publicKey.(*ecdsa.PublicKey)
	assert.True(t, ok)

	address := crypto.PubkeyToAddress(*publicKeyECDSA)
	return privateKey, publicKeyECDSA, address
}

func generateTestRows(t *testing.T, n int, ttl time.Duration) []*s4.Row {
	ormRows := generateTestOrmRows(t, n, ttl)
	rows := make([]*s4.Row, n)
	for i := 0; i < n; i++ {
		rows[i] = &s4.Row{
			Address:    ormRows[i].Address,
			Slotid:     uint32(ormRows[i].SlotId),
			Version:    ormRows[i].Version,
			Expiration: ormRows[i].Expiration,
			Payload:    ormRows[i].Payload,
			Signature:  ormRows[i].Signature,
		}
	}
	return rows
}

func generateTestOrmRows(t *testing.T, n int, ttl time.Duration) []*s4_orm.Row {
	rows := make([]*s4_orm.Row, n)
	for i := 0; i < n; i++ {
		priv, _, addr := generateCryptoEntity(t)
		rows[i] = &s4_orm.Row{
			Address:    addr.String(),
			SlotId:     0,
			Version:    0,
			Confirmed:  false,
			Expiration: time.Now().Add(ttl).UnixMilli(),
			Payload:    mustRandomBytes(t, 64),
			UpdatedAt:  time.Now().Add(-time.Second).UnixMilli(),
		}
		env := &s4_orm.Envelope{
			Address:    addr.Bytes(),
			SlotID:     rows[i].SlotId,
			Version:    rows[i].Version,
			Expiration: rows[i].Expiration,
			Payload:    rows[i].Payload,
		}
		sig, err := env.Sign(priv)
		assert.NoError(t, err)
		rows[i].Signature = sig
	}
	return rows
}

func compareRows(t *testing.T, protoRows []*s4.Row, ormRows []*s4_orm.Row) {
	assert.Equal(t, len(ormRows), len(protoRows))
	for i, row := range protoRows {
		assert.Equal(t, row.Address, ormRows[i].Address)
		assert.Equal(t, row.Version, ormRows[i].Version)
		assert.Equal(t, row.Expiration, ormRows[i].Expiration)
		assert.Equal(t, row.Payload, ormRows[i].Payload)
		assert.Equal(t, row.Signature, ormRows[i].Signature)
	}
}

func TestPlugin_Close(t *testing.T) {
	t.Parallel()

	config := createPluginConfig(10)
	orm := s4_mocks.NewORM(t)
	plugin, err := s4.NewConsensusPlugin(config, orm)
	assert.NoError(t, err)

	err = plugin.Close()
	assert.NoError(t, err)
}

func TestPlugin_ShouldTransmitAcceptedReport(t *testing.T) {
	t.Parallel()

	config := createPluginConfig(10)
	orm := s4_mocks.NewORM(t)
	plugin, err := s4.NewConsensusPlugin(config, orm)
	assert.NoError(t, err)

	should, err := plugin.ShouldTransmitAcceptedReport(testutils.Context(t), types.ReportTimestamp{}, nil)
	assert.NoError(t, err)
	assert.False(t, should)
}

func TestPlugin_ShouldAcceptFinalizedReport(t *testing.T) {
	t.Parallel()

	config := createPluginConfig(10)
	orm := s4_mocks.NewORM(t)
	plugin, err := s4.NewConsensusPlugin(config, orm)
	assert.NoError(t, err)

	t.Run("happy", func(t *testing.T) {
		ormRows := make([]*s4_orm.Row, 0)
		rows := generateTestRows(t, 10, time.Minute)
		orm.On("Update", mock.Anything, mock.Anything).Run(func(args mock.Arguments) {
			updateRow := args.Get(0).(*s4_orm.Row)
			ormRows = append(ormRows, updateRow)
		}).Return(nil).Times(10)

		report, err := proto.Marshal(&s4.Rows{
			Rows: rows,
		})
		assert.NoError(t, err)

		should, err := plugin.ShouldAcceptFinalizedReport(testutils.Context(t), types.ReportTimestamp{}, report)
		assert.NoError(t, err)
		assert.True(t, should)
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
		assert.ErrorIs(t, err, testErr)
		assert.False(t, should)
	})
}

func TestPlugin_Query(t *testing.T) {
	t.Parallel()

	config := createPluginConfig(10)
	orm := s4_mocks.NewORM(t)
	plugin, err := s4.NewConsensusPlugin(config, orm)
	assert.NoError(t, err)

	t.Run("happy", func(t *testing.T) {
		ormRows := generateTestOrmRows(t, 10, time.Minute)
		orm.On("GetSnapshot", s4_orm.MinAddress, s4_orm.MaxAddress, mock.Anything).Return(ormRows, nil).Once()

		query, err := plugin.Query(testutils.Context(t), types.ReportTimestamp{})
		assert.NoError(t, err)

		rows := &s4.Rows{}
		err = proto.Unmarshal(query, rows)
		assert.NoError(t, err)

		compareRows(t, rows.Rows, ormRows)
	})

	t.Run("empty", func(t *testing.T) {
		orm.On("GetSnapshot", s4_orm.MinAddress, s4_orm.MaxAddress, mock.Anything).Return(make([]*s4_orm.Row, 0), nil).Once()

		query, err := plugin.Query(testutils.Context(t), types.ReportTimestamp{})
		assert.NoError(t, err)
		assert.Nil(t, query)
	})

	t.Run("query with shards", func(t *testing.T) {
		config.NSnapshotShards = 16

		ormRows := generateTestOrmRows(t, 256, time.Minute)
		for i := 0; i < 256; i++ {
			var thisAddress common.Address
			thisAddress[0] = byte(i)
			ormRows[i].Address = thisAddress.String()
		}

		for i := 0; i <= int(config.NSnapshotShards); i++ {
			from := i * 16
			to := from + 16
			if i == int(config.NSnapshotShards) {
				from = 0
				to = 16
			}
			orm.On("GetSnapshot", mock.Anything, mock.Anything, mock.Anything).Return(ormRows[from:to], nil).Once()

			query, err := plugin.Query(testutils.Context(t), types.ReportTimestamp{})
			assert.NoError(t, err)

			rr := &s4.Rows{}
			err = proto.Unmarshal(query, rr)
			assert.NoError(t, err)

			assert.Len(t, rr.Rows, 16)
			for _, r := range rr.Rows {
				minAddress := mustStringToBig(t, rr.MinAddress)
				maxAddress := mustStringToBig(t, rr.MaxAddress)
				thisAddress := common.HexToAddress(r.Address).Big()
				assert.True(t, thisAddress.Cmp(minAddress) >= 0)
				assert.True(t, thisAddress.Cmp(maxAddress) <= 0)
			}
		}
	})
}

func TestPlugin_Observation(t *testing.T) {
	t.Parallel()

	config := createPluginConfig(10)
	orm := s4_mocks.NewORM(t)
	plugin, err := s4.NewConsensusPlugin(config, orm)
	assert.NoError(t, err)

	ormRows := generateTestOrmRows(t, 10, time.Minute)
	for i, or := range ormRows {
		or.Confirmed = i%2 == 0
	}
	orm.On("GetSnapshot", s4_orm.MinAddress, s4_orm.MaxAddress, mock.Anything).Return(ormRows, nil).Once()

	query, err := plugin.Query(testutils.Context(t), types.ReportTimestamp{})
	assert.NoError(t, err)

	orm.On("DeleteExpired", mock.Anything).Return(nil).Once()
	orm.On("GetSnapshot", mock.Anything, mock.Anything, mock.Anything).Return(ormRows, nil).Once()

	observation, err := plugin.Observation(testutils.Context(t), types.ReportTimestamp{}, query)
	assert.NoError(t, err)

	rows := &s4.Rows{}
	err = proto.Unmarshal(observation, rows)
	assert.NoError(t, err)
	assert.Len(t, rows.Rows, 5)
}

func TestPlugin_Report(t *testing.T) {
	t.Parallel()

	config := createPluginConfig(10)
	orm := s4_mocks.NewORM(t)
	plugin, err := s4.NewConsensusPlugin(config, orm)
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
}
