package v2

import (
	"path/filepath"
	"testing"

	"github.com/stretchr/testify/require"

	"github.com/smartcontractkit/chainlink/v2/core/internal/testutils/wasmtest"
)

const pathPrefix = "core/services/workflows/cmd/cre/examples/v2"

func Test_AllExampleWorkflowsCompileToWASM(t *testing.T) {
	paths := []string{
		"simple_cron",
		"simple_cron_with_config",
	}
	for _, path := range paths {
		t.Run(path, func(t *testing.T) {
			t.Parallel()
			binary := wasmtest.CreateTestBinary(filepath.Join(pathPrefix, path), false, t)
			require.NotEmpty(t, binary)
		})
	}
}
