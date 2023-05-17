package s4_test

import (
	"testing"

	"github.com/smartcontractkit/chainlink/v2/core/internal/testutils"
	"github.com/smartcontractkit/chainlink/v2/core/services/s4"
	"github.com/stretchr/testify/assert"
)

func TestClone_Record(t *testing.T) {
	t.Parallel()

	payload := testutils.Random32Byte()
	r := s4.Record{
		Payload:    payload[:],
		Version:    3,
		Expiration: 123,
	}

	c := r.Clone()
	assert.EqualValues(t, r, c)
}

func TestClone_Metadata(t *testing.T) {
	t.Parallel()

	signature := testutils.Random32Byte()
	m := s4.Metadata{
		State:             s4.ExpiredRecordState,
		HighestExpiration: 123,
		Signature:         signature[:],
	}

	c := m.Clone()
	assert.EqualValues(t, m, c)
}
