package smoke

import (
	"testing"

	"github.com/smartcontractkit/chainlink-testing-framework/logging"
)

// TODO: REMOVE WHEN THE CI WORKFLOWS TESTING IS DONE
func TestCICheck_Passing(t *testing.T) {
	l := logging.GetTestLogger(t)
	l.Info().Msg("Quick CI check. Always passing")
}
