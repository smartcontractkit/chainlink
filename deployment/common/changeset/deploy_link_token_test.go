package changeset_test

import (
	"testing"

	"github.com/smartcontractkit/chainlink/deployment/ccip/changeset/testhelpers"
)

func TestDeployLinkToken(t *testing.T) {
	testhelpers.DoDeployLinkToken(t, 0)
}
