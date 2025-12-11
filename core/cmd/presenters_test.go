package cmd_test

import (
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"

	"github.com/smartcontractkit/chainlink-common/pkg/utils/tests"
	"github.com/smartcontractkit/chainlink/v2/core/cmd"
)

func TestJAID(t *testing.T) {
	t.Parallel()
	tests.BelongsToCISuite(t, "with-db")

	jaid := cmd.JAID{ID: "1"}

	t.Run("GetID", func(t *testing.T) { assert.Equal(t, "1", jaid.GetID()) })
	t.Run("SetID", func(t *testing.T) {
		err := jaid.SetID("2")
		require.NoError(t, err)
		assert.Equal(t, "2", jaid.GetID())
	})
}
