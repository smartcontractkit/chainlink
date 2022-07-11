package soak_test

//revive:disable:dot-imports
import (
	"testing"

	. "github.com/onsi/ginkgo/v2"
	"github.com/smartcontractkit/chainlink-testing-framework/actions"
)

func Test_Suite(t *testing.T) {
	actions.GinkgoSuite()
	RunSpecs(t, "Soak")
}
