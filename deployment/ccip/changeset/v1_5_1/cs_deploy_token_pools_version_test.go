package v1_5_1

import (
	"testing"

	"github.com/Masterminds/semver/v3"
	"github.com/stretchr/testify/assert"

	cldf "github.com/smartcontractkit/chainlink-deployments-framework/deployment"
	"github.com/smartcontractkit/chainlink/deployment"
	"github.com/smartcontractkit/chainlink/deployment/ccip/shared"
)

// Pins deployedTypeAndVersion against the versions the deploy switch previously assigned inline.
// The datastore key is built from this, so a pool keyed under a version it was not deployed at
// would be unfindable.
func TestDeployedTypeAndVersionMatchesDeploySwitch(t *testing.T) {
	v161 := deployment.Version1_6_1
	for _, tc := range []struct {
		poolType    string
		cfgVersion  semver.Version
		wantVersion semver.Version
	}{
		{"BurnMintTokenPool", v161, deployment.Version1_6_1},
		{"BurnMintTokenPool", semver.Version{}, shared.CurrentTokenPoolVersion},
		{"LockReleaseTokenPool", v161, deployment.Version1_6_1},
		{"LockReleaseTokenPool", semver.Version{}, shared.CurrentTokenPoolVersion},
		{"BurnMintFastTransferTokenPool", semver.Version{}, deployment.Version1_6_3Dev},
		{"BurnMintWithExternalMinterFastTransferTokenPool", semver.Version{}, deployment.Version1_6_0},
		{"HybridWithExternalMinterFastTransferTokenPool", semver.Version{}, deployment.Version1_6_0},
		{"BurnMintWithExternalMinterTokenPool", semver.Version{}, deployment.Version1_6_0},
		{"HybridWithExternalMinterTokenPool", semver.Version{}, deployment.Version1_6_0},
		{"BurnWithFromMintTokenPool", semver.Version{}, shared.CurrentTokenPoolVersion},
		{"BurnFromMintTokenPool", semver.Version{}, shared.CurrentTokenPoolVersion},
	} {
		got := deployedTypeAndVersion(DeployTokenPoolInput{
			Type:    cldf.ContractType(tc.poolType),
			Version: tc.cfgVersion,
		})
		assert.Equal(t, tc.wantVersion.String(), got.Version.String(), "%s cfg=%s", tc.poolType, tc.cfgVersion)
		assert.Equal(t, tc.poolType, string(got.Type))
	}
}
