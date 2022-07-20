package cmd_test

import (
	"testing"

	"github.com/stretchr/testify/assert"

	"github.com/smartcontractkit/chainlink/core/cmd"
)

func TestJAID(t *testing.T) {
	t.Parallel()

	jaid := cmd.JAID{ID: "1"}

	t.Run("GetID", func(t *testing.T) { assert.Equal(t, "1", jaid.GetID()) })
	t.Run("SetID", func(t *testing.T) {
		jaid.SetID("2")
		assert.Equal(t, "2", jaid.GetID())
	})
}
