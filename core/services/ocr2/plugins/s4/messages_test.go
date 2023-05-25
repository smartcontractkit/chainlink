package s4_test

import (
	"testing"
	"time"

	"github.com/smartcontractkit/chainlink/v2/core/services/ocr2/plugins/s4"
	s4_svc "github.com/smartcontractkit/chainlink/v2/core/services/s4"

	"github.com/stretchr/testify/assert"
)

func Test_MarshalUnmarshalRows(t *testing.T) {
	t.Parallel()

	const n = 100
	rows := generateTestRows(t, n, time.Minute)

	t.Run("no address range", func(t *testing.T) {
		data, err := s4.MarshalRows(rows, nil)
		assert.NoError(t, err)

		rr, ar, err := s4.UnmarshalRows(data)
		assert.NoError(t, err)
		assert.Nil(t, ar)
		assert.Len(t, rr, n)
	})

	t.Run("with address range", func(t *testing.T) {
		ar, err := s4_svc.NewInitialAddressRangeForIntervals(16)
		assert.NoError(t, err)

		data, err := s4.MarshalRows(rows, ar)
		assert.NoError(t, err)

		rr, arr, err := s4.UnmarshalRows(data)
		assert.NoError(t, err)
		assert.Len(t, rr, n)
		assert.NotNil(t, arr)
		assert.Equal(t, ar.MaxAddress, arr.MaxAddress)
		assert.Equal(t, ar.MinAddress, arr.MinAddress)
	})
}
