package changeset_test

import (
	"testing"

	"github.com/smartcontractkit/chainlink/deployment/common/changeset"
)

func TestDeployLinkToken(t *testing.T) {
	t.Parallel()
	changeset.DeployLinkTokenTest(t, 0)
}
