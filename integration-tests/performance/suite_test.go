package performance_test

//revive:disable:dot-imports
import (
	"testing"

	"github.com/smartcontractkit/chainlink-testing-framework/actions"
	networks "github.com/smartcontractkit/chainlink/integration-tests"

	. "github.com/onsi/ginkgo/v2"
)

func Test_Suite(t *testing.T) {
	actions.GinkgoSuite()
	networks.LoadNetworks("../.env")

	RunSpecs(t, "Profiling")
}
