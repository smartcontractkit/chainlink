package benchmark_test

//revive:disable:dot-imports
import (
	"github.com/onsi/ginkgo/v2"
	"testing"

	"github.com/smartcontractkit/chainlink/integration-tests/actions"
)

func Test_Suite(t *testing.T) {
	actions.GinkgoSuite()
	ginkgo.RunSpecs(t, "Benchmark")
}
