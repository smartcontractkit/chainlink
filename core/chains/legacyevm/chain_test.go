package legacyevm_test

import (
	"math/big"
	"testing"

	"github.com/stretchr/testify/assert"

	"github.com/smartcontractkit/chainlink-common/pkg/sqlutil/sqltest"
	"github.com/smartcontractkit/chainlink-common/pkg/utils/mailbox"

	"github.com/smartcontractkit/chainlink/v2/core/chains/legacyevm"
	"github.com/smartcontractkit/chainlink/v2/core/chains/legacyevm/mocks"
	"github.com/smartcontractkit/chainlink/v2/core/internal/testutils/configtest"
)

func TestLegacyChains(t *testing.T) {
	legacyevmCfg := configtest.NewGeneralConfig(t, nil)

	c := mocks.NewChain(t)
	c.On("ID").Return(big.NewInt(7))
	m := map[string]legacyevm.Chain{c.ID().String(): c}

	l := legacyevm.NewLegacyChains(m, legacyevmCfg.EVMConfigs())
	assert.NotNil(t, l.ChainNodeConfigs())
	got, err := l.Get(c.ID().String())
	assert.NoError(t, err)
	assert.Equal(t, c, got)
}

func TestChainOpts_Validate(t *testing.T) {
	cfg := configtest.NewTestGeneralConfig(t)
	tests := []struct {
		name    string
		opts    legacyevm.ChainOpts
		wantErr bool
	}{
		{
			name: "valid",
			opts: legacyevm.ChainOpts{
				ChainConfigs:   cfg.EVMConfigs(),
				DatabaseConfig: cfg.Database(),
				ListenerConfig: cfg.Database().Listener(),
				FeatureConfig:  cfg.Feature(),
				MailMon:        &mailbox.Monitor{},
				DS:             sqltest.NewNoOpDataSource(),
			},
		},
		{
			name:    "invalid",
			opts:    legacyevm.ChainOpts{},
			wantErr: true,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if err := tt.opts.Validate(); (err != nil) != tt.wantErr {
				t.Errorf("ChainOpts.Validate() error = %v, wantErr %v", err, tt.wantErr)
			}
		})
	}
}
