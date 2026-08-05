package stellar

import (
	"testing"

	"github.com/stretchr/testify/require"

	cldregistry "github.com/smartcontractkit/chainlink-deployments-framework/engine/cld/changeset"
)

// Registers every changeset the way a chainlink-deployments
// domain will and lists the keys, mirroring cld durable-pipeline list.
func TestChangesetRegistryKeys(t *testing.T) {
	r := cldregistry.NewChangesetsRegistry()

	r.Add("0001_deploy_cache", cldregistry.Configure(DeployCache{}).WithEnvInput())
	r.Add("0002_deploy_proxy", cldregistry.Configure(DeployProxy{}).WithEnvInput())
	r.Add("0003_set_feed_configs", cldregistry.Configure(SetFeedConfigs{}).WithEnvInput())
	r.Add("0004_remove_feed_configs", cldregistry.Configure(RemoveFeedConfigs{}).WithEnvInput())
	r.Add("0005_set_feed_frozen", cldregistry.Configure(SetFeedFrozen{}).WithEnvInput())
	r.Add("0006_add_feed_admin", cldregistry.Configure(AddFeedAdmin{}).WithEnvInput())
	r.Add("0007_remove_feed_admin", cldregistry.Configure(RemoveFeedAdmin{}).WithEnvInput())
	r.Add("0008_transfer_ownership", cldregistry.Configure(TransferOwnership{}).WithEnvInput())
	r.Add("0009_accept_ownership", cldregistry.Configure(AcceptOwnership{}).WithEnvInput())
	r.Add("0010_upgrade", cldregistry.Configure(Upgrade{}).WithEnvInput())
	r.Add("0011_recover_tokens", cldregistry.Configure(RecoverTokens{}).WithEnvInput())
	r.Add("0012_set_proxy_cache", cldregistry.Configure(SetProxyCache{}).WithEnvInput())

	keys := r.ListKeys()
	require.Len(t, keys, 12)
	for _, k := range keys {
		t.Logf("registered: %s", k)
	}
}
