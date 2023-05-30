package s4_test

import (
	"testing"
	"time"

	"github.com/smartcontractkit/chainlink/v2/core/services/ocr2/plugins/s4"

	"github.com/stretchr/testify/assert"
)

func Test_MarshalUnmarshalRows(t *testing.T) {
	t.Parallel()

	const n = 100
	rows := generateTestRows(t, n, time.Minute)

	data, err := s4.MarshalRows(rows)
	assert.NoError(t, err)

	rr, err := s4.UnmarshalRows(data)
	assert.NoError(t, err)
	assert.Len(t, rr, n)
}

func Test_MarshalUnmarshalQuery(t *testing.T) {
	t.Parallel()

	const n = 100
	rows := generateTestOrmRows(t, n, time.Minute)
	ormVersions := rowsToVersions(rows)

	versions := make([]*s4.VersionRow, len(ormVersions))
	for i, v := range ormVersions {
		versions[i] = &s4.VersionRow{
			Address: v.Address.Hex(),
			Slotid:  uint32(v.SlotId),
			Version: v.Version,
		}
	}
	data, err := s4.MarshalQuery(versions)
	assert.NoError(t, err)

	qq, err := s4.UnmarshalQuery(data)
	assert.NoError(t, err)
	assert.Len(t, qq, n)
}

func Test_VerifySignature(t *testing.T) {
	t.Parallel()

	rows := generateTestRows(t, 2, time.Minute)
	err := rows[0].VerifySignature()
	assert.NoError(t, err)

	rows[1].Signature[0] = ^rows[1].Signature[0]
	err = rows[1].VerifySignature()
	assert.Error(t, err)
}
