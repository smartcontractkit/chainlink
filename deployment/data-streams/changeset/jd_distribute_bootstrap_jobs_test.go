package changeset

import (
	"testing"

	"github.com/stretchr/testify/require"

	"github.com/smartcontractkit/chainlink/deployment"
	"github.com/smartcontractkit/chainlink/deployment/common/changeset"
	"github.com/smartcontractkit/chainlink/deployment/data-streams/changeset/testutil"
	"github.com/smartcontractkit/chainlink/deployment/data-streams/jd"
)

func TestDistributeBootstrapJobSpecs(t *testing.T) {
	t.Skip("Flaky Test: https://smartcontract-it.atlassian.net/browse/DX-196")
	t.Parallel()

	e := testutil.NewMemoryEnv(t, false, 1)

	// pick the first EVM chain selector
	chainSelector := e.AllChainSelectors()[0]

	// insert a Configurator address for the given DON
	err := e.ExistingAddresses.Save(chainSelector, "0x4170ed0880ac9a755fd29b2688956bd959f923f4",
		deployment.TypeAndVersion{
			Type:    "Configurator",
			Version: deployment.Version1_0_0,
			Labels:  deployment.NewLabelSet("don-1"),
		})
	require.NoError(t, err)

	config := CsDistributeBootstrapJobSpecsConfig{
		ChainSelectorEVM: chainSelector,
		Filter: &jd.ListFilter{
			DONID:    1,
			DONName:  "don",
			EnvLabel: "env",
			Size:     0,
		},
	}

	tests := []struct {
		name    string
		env     deployment.Environment
		config  CsDistributeBootstrapJobSpecsConfig
		wantErr *string
	}{
		{
			name:   "success",
			env:    e,
			config: config,
		},
	}

	cs := CsDistributeBootstrapJobSpecs{}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			_, err = changeset.ApplyChangesetsV2(t,
				tt.env,
				[]changeset.ConfiguredChangeSet{
					changeset.Configure(cs, tt.config),
				},
			)

			if tt.wantErr == nil {
				require.NoError(t, err)
			} else {
				require.Error(t, err)
				require.Contains(t, err.Error(), tt.wantErr)
			}
		})
	}
}
