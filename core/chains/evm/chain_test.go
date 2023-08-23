package evm_test

import (
	"math/big"
	"testing"

	"github.com/stretchr/testify/assert"

	"github.com/smartcontractkit/chainlink/v2/core/chains/evm"
	"github.com/smartcontractkit/chainlink/v2/core/chains/evm/mocks"
	configtest "github.com/smartcontractkit/chainlink/v2/core/internal/testutils/configtest/v2"
	"github.com/smartcontractkit/chainlink/v2/core/services/chainlink"
	"github.com/smartcontractkit/chainlink/v2/core/utils"
)

func TestLegacyChains(t *testing.T) {
	evmCfg := configtest.NewGeneralConfig(t, nil)

	c := mocks.NewChain(t)
	c.On("ID").Return(big.NewInt(7))
	m := map[string]evm.Chain{c.ID().String(): c}

	l, err := evm.NewLegacyChains(evmCfg, m)
	assert.NoError(t, err)
	assert.NotNil(t, l.ChainNodeConfigs())
	got, err := l.Get(c.ID().String())
	assert.NoError(t, err)
	assert.Equal(t, c, got)

	l, err = evm.NewLegacyChains(nil, m)
	assert.Error(t, err)
	assert.Nil(t, l)
}

func TestRelayConfigInit(t *testing.T) {
	appCfg := configtest.NewGeneralConfig(t, nil)
	rCfg := evm.RelayerConfig{
		AppConfig: appCfg,
	}

	evmCfg := rCfg.EVMConfigs()
	assert.NotNil(t, evmCfg)

	// test lazy init is done only once
	// note this kind of swapping should never happen in prod
	appCfg2 := configtest.NewGeneralConfig(t, func(c *chainlink.Config, s *chainlink.Secrets) {
		c.EVM[0].ChainID = utils.NewBig(big.NewInt(27))
	})
	rCfg.AppConfig = appCfg2

	newEvmCfg := rCfg.EVMConfigs()
	assert.NotNil(t, newEvmCfg)
	assert.Equal(t, evmCfg, newEvmCfg)
}
