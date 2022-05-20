package performance_test

//revive:disable:dot-imports
import (
	"testing"

	"github.com/smartcontractkit/chainlink-testing-framework/actions"
	"github.com/smartcontractkit/chainlink-testing-framework/utils"

	. "github.com/onsi/ginkgo/v2"
)

func Test_Suite(t *testing.T) {
	actions.GinkgoSuite(utils.ProjectRoot)
	RunSpecs(t, "Profiling")
}
