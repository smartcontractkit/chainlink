package changeset_test

import (
	"testing"

	"github.com/smartcontractkit/chainlink-common/pkg/utils/tests"
	"github.com/smartcontractkit/chainlink/deployment/common/changeset"
	"github.com/smartcontractkit/chainlink/deployment/environment/memory"
)

func TestDeployLinkToken(t *testing.T) {
	t.Parallel()
	changeset.DeployLinkTokenTest(t, memory.MemoryEnvironmentConfig{
		Chains: 1,
	})
}

func TestDeployLinkTokenZk(t *testing.T) {
	// Timeouts in CI
	tests.SkipFlakey(t, "https://smartcontract-it.atlassian.net/browse/CCIP-6427")

	t.Parallel()
	changeset.DeployLinkTokenTest(t, memory.MemoryEnvironmentConfig{
		ZkChains: 1,
	})
}
