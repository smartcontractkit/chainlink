package s4_test

import (
	"testing"
	"time"

	"github.com/smartcontractkit/chainlink/v2/core/internal/testutils"
	"github.com/smartcontractkit/chainlink/v2/core/services/ocr2/plugins/s4"
	s4_svc "github.com/smartcontractkit/chainlink/v2/core/services/s4"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func Test_MarshalUnmarshalRows(t *testing.T) {
	t.Parallel()

	const n = 1000
	rows := generateTestRows(t, n, time.Minute)

	data, err := s4.MarshalRows(rows)
	require.NoError(t, err)

	rr, err := s4.UnmarshalRows(data)
	require.NoError(t, err)
	require.Len(t, rr, n)

	data2, err := s4.MarshalRows(rr)
	require.NoError(t, err)
	require.Equal(t, data, data2)
}

func Test_MarshalUnmarshalQuery(t *testing.T) {
	t.Parallel()

	const n = 100
	rows := generateTestOrmRows(t, n, time.Minute)
	ormVersions := rowsToShapshotRows(rows)

	snapshot := make([]*s4.SnapshotRow, len(ormVersions))
	for i, v := range ormVersions {
		snapshot[i] = &s4.SnapshotRow{
			Address: v.Address.Bytes(),
			Slotid:  uint32(v.SlotId),
			Version: v.Version,
		}
	}
	addressRange := s4_svc.NewFullAddressRange()
	data, err := s4.MarshalQuery(snapshot, addressRange)
	require.NoError(t, err)

	qq, ar, err := s4.UnmarshalQuery(data)
	require.NoError(t, err)
	require.Len(t, qq, n)
	require.Equal(t, addressRange, ar)
}

func Test_VerifySignature(t *testing.T) {
	t.Parallel()

	rows := generateTestRows(t, 2, time.Minute)
	err := rows[0].VerifySignature()
	require.NoError(t, err)

	rows[1].Signature[0] = ^rows[1].Signature[0]
	err = rows[1].VerifySignature()
	require.Error(t, err)

	t.Run("address with leading zeros", func(t *testing.T) {
		pk, addr := testutils.NewPrivateKeyAndAddress(t)
		for addr[0] != 0 {
			pk, addr = testutils.NewPrivateKeyAndAddress(t)
		}
		rows := generateTestRows(t, 1, time.Minute)
		rows[0].Address = addr.Big().Bytes()
		env := &s4_svc.Envelope{
			Address:    addr.Bytes(),
			SlotID:     uint(rows[0].Slotid),
			Version:    rows[0].Version,
			Expiration: rows[0].Expiration,
			Payload:    rows[0].Payload,
		}
		sig, err := env.Sign(pk)
		assert.NoError(t, err)
		rows[0].Signature = sig

		err = rows[0].VerifySignature()
		require.NoError(t, err)
	})
}
