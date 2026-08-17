package examples

import (
	"testing"

	"github.com/stretchr/testify/require"

	"github.com/smartcontractkit/chainlink-common/pkg/wasmbuild"
)

func Test_AllExampleWorkflowsCompileToWASM(t *testing.T) {
	t.Parallel()
	paths := []string{
		"v2/http_read",
		"v2/simple_cron",
		"v2/simple_cron_with_config",
		"v2/simple_cron_with_secrets",
		"v2/empty",
	}
	for _, path := range paths {
		t.Run(path, func(t *testing.T) {
			t.Parallel()
			binary, err := wasmbuild.Compile(t.Context(), wasmbuild.Config{
				PkgDir: path,
			})
			require.NoError(t, err)
			require.NotEmpty(t, binary)
		})
	}
}
