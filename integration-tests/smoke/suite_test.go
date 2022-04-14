package smoke_test

//revive:disable:dot-imports
import (
	"testing"

	"github.com/smartcontractkit/integrations-framework/actions"
	"github.com/smartcontractkit/integrations-framework/utils"

	. "github.com/onsi/ginkgo/v2"
)

func Test_Suite(t *testing.T) {
	actions.GinkgoSuite(utils.ProjectRoot)
	RunSpecs(t, "Integration")
}
