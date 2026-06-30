package examples

//go:generate go run ../../../../../internal/testutils/wasmtest/generator/main.go -pkg core/services/workflows/cmd/cre/examples/legacy/data_feeds
//go:generate go run ../../../../../internal/testutils/wasmtest/generator/main.go -pkg core/services/workflows/cmd/cre/examples/v2/http_read
//go:generate go run ../../../../../internal/testutils/wasmtest/generator/main.go -pkg core/services/workflows/cmd/cre/examples/v2/simple_cron
//go:generate go run ../../../../../internal/testutils/wasmtest/generator/main.go -pkg core/services/workflows/cmd/cre/examples/v2/simple_cron -compress
//go:generate go run ../../../../../internal/testutils/wasmtest/generator/main.go -pkg core/services/workflows/cmd/cre/examples/v2/simple_cron_with_config
//go:generate go run ../../../../../internal/testutils/wasmtest/generator/main.go -pkg core/services/workflows/cmd/cre/examples/v2/simple_cron_with_secrets
//go:generate go run ../../../../../internal/testutils/wasmtest/generator/main.go -pkg core/services/workflows/cmd/cre/examples/v2/empty
import (
	"path/filepath"
	"testing"

	"github.com/stretchr/testify/require"

	"github.com/smartcontractkit/chainlink/v2/core/internal/testutils/wasmtest"
)

const pathPrefix = "core/services/workflows/cmd/cre/examples"

func Test_AllExampleWorkflowsCompileToWASM(t *testing.T) {
	t.Parallel()
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
			binary := wasmtest.GetTestBinary(t, filepath.Join(pathPrefix, path), false)
			require.NotEmpty(t, binary)
		})
	}
}
