package smoke_test

//revive:disable:dot-imports
import (
	"testing"

	"github.com/smartcontractkit/chainlink-testing-framework/actions"

	. "github.com/onsi/ginkgo/v2"
)

func Test_Suite(t *testing.T) {
	actions.GinkgoSuite("../")
	RunSpecs(t, "Integration")
}
