package smoke

import (
	"testing"

	cciptests "github.com/smartcontractkit/chainlink/integration-tests/deployment/ccip/tests"
	"github.com/smartcontractkit/chainlink/v2/core/logger"
)

func Test0002_InitialDeployOnLocal(t *testing.T) {
	lggr := logger.TestLogger(t)
	tenv := cciptests.NewLocalDevEnvironment(t, lggr)
	cciptests.InitialDeployTest(t, tenv)
}
