package smoke_test

import (
	"testing"

	"github.com/smartcontractkit/integrations-framework/utils"

	. "github.com/onsi/ginkgo/v2"
)

func Test_Suite(t *testing.T) {
	utils.GinkgoSuite()
	RunSpecs(t, "Integration")
}
