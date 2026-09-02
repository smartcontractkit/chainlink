package matrix_test

import (
	"context"
	"os"
	"path/filepath"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"

	"github.com/smartcontractkit/chainlink/v2/tools/ci/internal/matrix"
)

func TestBuildInMemoryMatrix(t *testing.T) {
	t.Parallel()

	tmpDir := t.TempDir()
	configFile := filepath.Join(tmpDir, "in-memory-tests.json")
	content := `[
  {"name":"ccip_fees_test.go","test":"^(Test_CCIPFees)$","timeout":"20m","parallel":4,"plugins":false,"runs_on":"cpu=8/ram=32","free_disk":false,"aptos":"","sui":""},
  {"name":"Test_CCIPMessaging_Solana2EVM","test":"Test_CCIPMessaging_Solana2EVM","timeout":"18m","parallel":4,"plugins":true,"runs_on":"","free_disk":false,"aptos":"","sui":""}
]`
	require.NoError(t, os.WriteFile(configFile, []byte(content), 0o600))

	res, err := matrix.BuildInMemoryMatrix(context.Background(), matrix.InMemoryOptions{
		ConfigFile: configFile,
		RunID:      "999888",
		RunAttempt: "3",
		SpotFlag:   "spot=pco",
	})
	require.NoError(t, err)
	require.Len(t, res, 2)

	// Entry 0: custom runs_on
	first := res[0]
	assert.Equal(t, "ccip_fees_test.go", first.Name)
	assert.Equal(t, "ccip_fees_test.go", first.TestID) // sanitized name
	assert.Equal(t, 25, first.JobTimeout)              // 20m + 5 = 25
	assert.Equal(t, "runs-on=999888-0-3/cpu=8/ram=32/family=m7i+m8i/spot=pco/image=ubuntu24-full-x64/extras=s3-cache+tmpfs", first.RunsOn)
	assert.Equal(t, 4, first.Parallel)
	assert.False(t, first.Plugins)

	// Entry 1: empty runs_on -> ubuntu-latest
	second := res[1]
	assert.Equal(t, "Test_CCIPMessaging_Solana2EVM", second.Name)
	assert.Equal(t, "Test_CCIPMessaging_Solana2EVM", second.TestID)
	assert.Equal(t, 23, second.JobTimeout) // 18m + 5 = 23
	assert.Equal(t, "ubuntu-latest", second.RunsOn)
	assert.True(t, second.Plugins)
}

func TestSanitizeTestID(t *testing.T) {
	t.Parallel()

	assert.Equal(t, "ccip_fees_test.go", matrix.SanitizeTestID("ccip_fees_test.go"))
	assert.Equal(t, "Test_Foo_Bar_Baz.123-v1", matrix.SanitizeTestID("Test/Foo Bar:Baz.123-v1"))
}
