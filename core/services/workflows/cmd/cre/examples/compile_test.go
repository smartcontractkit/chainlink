package examples

import (
	"path/filepath"
	"testing"

	"github.com/stretchr/testify/require"

	"github.com/smartcontractkit/chainlink/v2/core/internal/testutils/wasmtest"
)

const pathPrefix = "core/services/workflows/cmd/cre/examples"

func Test_AllExampleWorkflowsCompileToWASM(t *testing.T) {
	paths := []string{
		"legacy/data_feeds",
		"v2/http_read",
		"v2/simple_cron",
		"v2/simple_cron_with_config",
		"v2/simple_cron_with_secrets",
		"v2/empty",
	}
	for _, path := range paths {
		t.Run(path, func(t *testing.T) {
			t.Parallel()
			binary := wasmtest.CreateTestBinary(filepath.Join(pathPrefix, path), false, t)
			require.NotEmpty(t, binary)
		})
	}
}
