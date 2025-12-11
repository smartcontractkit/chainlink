package docs_test

import (
	"testing"

	"github.com/smartcontractkit/chainlink-common/pkg/config/configtest"

	"github.com/smartcontractkit/chainlink-common/pkg/utils/tests"
	"github.com/smartcontractkit/chainlink/v2/core/config/docs"
	"github.com/smartcontractkit/chainlink/v2/core/services/chainlink"
)

func TestDocsTOMLComplete(t *testing.T) {
	tests.BelongsToCISuite(t, "unit")
	configtest.AssertDocsTOMLComplete[chainlink.Config](t, docs.DocsTOML)
}
