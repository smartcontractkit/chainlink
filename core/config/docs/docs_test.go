package docs_test

import (
	"testing"

	"github.com/smartcontractkit/chainlink-common/pkg/config/configtest"

	"github.com/smartcontractkit/chainlink/v2/core/config/docs"
	"github.com/smartcontractkit/chainlink/v2/core/services/chainlink"
)

func TestDocsTOMLComplete(t *testing.T) {
	configtest.AssertDocsTOMLComplete[chainlink.Config](t, docs.DocsTOML)
}
