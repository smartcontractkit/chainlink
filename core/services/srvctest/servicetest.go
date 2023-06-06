package srvctest

import (
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"

	"github.com/smartcontractkit/chainlink/v2/core/internal/testutils"
	"github.com/smartcontractkit/chainlink/v2/core/services"
)

// Start is test helper to automatically Start/Close a ServiceCtx along with a test.
func Start[S services.ServiceCtx](tb testing.TB, s S) S {
	require.NoError(tb, s.Start(testutils.Context(tb)))
	tb.Cleanup(func() { assert.NoError(tb, s.Close()) })
	return s
}
