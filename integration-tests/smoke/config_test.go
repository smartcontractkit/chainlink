package smoke

import (
	"testing"

	"github.com/stretchr/testify/require"

	tc "github.com/smartcontractkit/chainlink/integration-tests/testconfig"
)

func TestConfig(t *testing.T) {
	t.Parallel()

	_, err := tc.GetConfig(t.Name(), tc.Smoke, tc.OCR)
	if err != nil {
		t.Fatal(err)
	}

	require.True(t, false, "err on purpose")
}
