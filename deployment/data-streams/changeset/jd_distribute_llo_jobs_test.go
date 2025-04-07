package changeset

import (
	"testing"

	"github.com/ethereum/go-ethereum/common"
	"github.com/stretchr/testify/require"

	"github.com/smartcontractkit/chainlink/deployment"
	"github.com/smartcontractkit/chainlink/deployment/common/changeset"
	"github.com/smartcontractkit/chainlink/deployment/data-streams/changeset/testutil"
	"github.com/smartcontractkit/chainlink/deployment/data-streams/jd"
	"github.com/smartcontractkit/chainlink/deployment/data-streams/utils/pointer"
)

func TestDistributeLLOJobSpecs(t *testing.T) {
	t.Parallel()

	e := testutil.NewMemoryEnv(t, false, 1)

	// pick the first EVM chain selector
	chainSelector := e.AllChainSelectors()[0]

	// insert a Configurator address for the given DON
	configuratorAddr := "0x4170ed0880ac9a755fd29b2688956bd959f923f4"
	err := e.ExistingAddresses.Save(chainSelector, configuratorAddr,
		deployment.TypeAndVersion{
			Type:    "Configurator",
			Version: deployment.Version1_0_0,
			Labels:  deployment.NewLabelSet("don-1"),
		})
	require.NoError(t, err)

	config := CsDistributeLLOJobSpecsConfig{
		ChainSelectorEVM: chainSelector,
		Filter: &jd.ListFilter{
			DONID:    1,
			DONName:  "don",
			EnvLabel: "env",
			Size:     0,
		},
		FromBlock:                   0,
		ConfigMode:                  "bluegreen",
		ChannelConfigStoreAddr:      common.HexToAddress("DEAD"),
		ChannelConfigStoreFromBlock: 0,
		ConfiguratorAddress:         configuratorAddr,
		Servers: map[string]string{
			"mercury-pipeline-testnet-producer.TEST.cldev.cloud:1340": "0000005187b1498c0ccb2e56d5ee8040a03a4955822ed208749b474058fc3f9c",
		},
	}

	tests := []struct {
		name       string
		env        deployment.Environment
		config     CsDistributeLLOJobSpecsConfig
		prepConfFn func(CsDistributeLLOJobSpecsConfig) CsDistributeLLOJobSpecsConfig
		wantErr    *string
	}{
		{
			name:   "success",
			env:    e,
			config: config,
		},
		{
			name:   "missing channel config store",
			env:    e,
			config: config,
			prepConfFn: func(c CsDistributeLLOJobSpecsConfig) CsDistributeLLOJobSpecsConfig {
				c.ChannelConfigStoreAddr = common.Address{}
				return c
			},
			wantErr: pointer.To("channel config store address is required"),
		},
		{
			name:   "missing servers",
			env:    e,
			config: config,
			prepConfFn: func(c CsDistributeLLOJobSpecsConfig) CsDistributeLLOJobSpecsConfig {
				c.Servers = nil
				return c
			},
			wantErr: pointer.To("servers map is required"),
		},
	}

	cs := CsDistributeLLOJobSpecs{}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			conf := tt.config
			if tt.prepConfFn != nil {
				conf = tt.prepConfFn(tt.config)
			}
			_, err = changeset.ApplyChangesetsV2(t,
				tt.env,
				[]changeset.ConfiguredChangeSet{
					changeset.Configure(cs, conf),
				},
			)

			if tt.wantErr == nil {
				require.NoError(t, err)
			} else {
				require.Error(t, err)
				require.Contains(t, err.Error(), *tt.wantErr)
			}
		})
	}
}
