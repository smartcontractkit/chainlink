package solana_test

import (
	"math/big"
	"net/url"
	"testing"
	"time"

	"github.com/ethereum/go-ethereum/common"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"

	cldfproposalutils "github.com/smartcontractkit/chainlink-deployments-framework/engine/cld/mcms/proposalutils"
	"github.com/smartcontractkit/chainlink-deployments-framework/engine/test/environment"

	"github.com/smartcontractkit/chainlink/deployment/cre/mcms/pkg"
	cresolmcms "github.com/smartcontractkit/chainlink/deployment/cre/mcms/solana"
)

func TestDeploySolanaMCMSConfig_Qualifier(t *testing.T) {
	t.Parallel()

	desc := "platform"
	cfg := cresolmcms.DeploySolanaMCMSConfig{
		ConfigID:      "staging",
		ChainSelector: 16423721717087811551,
		Descriptor:    &desc,
	}

	assertQualifierEqual(t,
		"contract://16423721717087811551/mcmsv2?mcms-config=staging&descriptor=platform",
		cfg.Qualifier(),
	)

	cfgNoDesc := cresolmcms.DeploySolanaMCMSConfig{
		ConfigID:      "staging",
		ChainSelector: 16423721717087811551,
	}
	assertQualifierEqual(t,
		"contract://16423721717087811551/mcmsv2?mcms-config=staging",
		cfgNoDesc.Qualifier(),
	)
}

func assertQualifierEqual(t *testing.T, expected, actual string) {
	t.Helper()

	exp, err := url.Parse(expected)
	require.NoError(t, err)
	got, err := url.Parse(actual)
	require.NoError(t, err)

	assert.Equal(t, exp.Scheme, got.Scheme)
	assert.Equal(t, exp.Host, got.Host)
	assert.Equal(t, exp.Path, got.Path)
	assert.Equal(t, exp.Query(), got.Query())
}

func TestDeploySolanaMCMS_VerifyPreconditions(t *testing.T) {
	t.Parallel()

	env, err := environment.New(t.Context())
	require.NoError(t, err)

	cs := cresolmcms.DeploySolanaMCMS{}
	stagingCfg := testStagingMCMSConfig()

	err = cs.VerifyPreconditions(*env, cresolmcms.DeploySolanaMCMSConfig{
		ConfigID:               "staging",
		MCMSWithTimelockConfig: &stagingCfg,
	})
	require.ErrorContains(t, err, "chainSelector is required")

	err = cs.VerifyPreconditions(*env, cresolmcms.DeploySolanaMCMSConfig{
		ChainSelector:          5009297550715157269,
		ConfigID:               "staging",
		MCMSWithTimelockConfig: &stagingCfg,
	})
	require.ErrorContains(t, err, "not a solana chain")

	err = cs.VerifyPreconditions(*env, cresolmcms.DeploySolanaMCMSConfig{
		ChainSelector:          16423721717087811551,
		MCMSWithTimelockConfig: &stagingCfg,
	})
	require.ErrorContains(t, err, "configId is required")

	err = cs.VerifyPreconditions(*env, cresolmcms.DeploySolanaMCMSConfig{
		ChainSelector: 16423721717087811551,
		ConfigID:      "staging",
	})
	require.ErrorContains(t, err, "mcmsWithTimelockConfig is required")
}

func testStagingMCMSConfig() cldfproposalutils.MCMSWithTimelockConfig {
	eoa := common.HexToAddress("0xA01E9eD15b18D3688D0B84D88a98ed750D56999B")
	d := 5 * time.Second
	return cldfproposalutils.MCMSWithTimelockConfig{
		Proposer:         pkg.MustGetMCMSConfig(1, []common.Address{eoa}, nil),
		Bypasser:         pkg.MustGetMCMSConfig(1, []common.Address{eoa}, nil),
		Canceller:        pkg.MustGetMCMSConfig(1, []common.Address{eoa}, nil),
		TimelockMinDelay: big.NewInt(int64(d.Seconds())),
	}
}
