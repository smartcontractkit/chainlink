package srvctest

import (
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"

	"github.com/smartcontractkit/chainlink/core/internal/testutils"
	"github.com/smartcontractkit/chainlink/core/services"
)

// Start is test helper to automatically Start/Close a ServiceCtx along with a test.
func Start[S services.ServiceCtx](t *testing.T, s S) S {
	require.NoError(t, s.Start(testutils.Context(t)))
	t.Cleanup(func() { assert.NoError(t, s.Close()) })
	return s
}
