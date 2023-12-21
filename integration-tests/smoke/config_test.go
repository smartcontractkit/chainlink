package smoke

import (
	"fmt"
	"testing"

	"github.com/smartcontractkit/chainlink-testing-framework/logging"

	tc "github.com/smartcontractkit/chainlink/integration-tests/testconfig"
)

func TestConfig(t *testing.T) {
	t.Parallel()
	l := logging.GetTestLogger(t)

	config, err := tc.GetConfig(t.Name(), tc.Smoke, tc.OCR)
	if err != nil {
		t.Fatal(err)
	}

	l.Info().Msg(fmt.Sprintf("%+v", config))
}
