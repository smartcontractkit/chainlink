package evm_test

import (
	"math/big"
	"testing"

	"github.com/stretchr/testify/assert"

	"github.com/smartcontractkit/chainlink/v2/core/chains/evm"
	"github.com/smartcontractkit/chainlink/v2/core/chains/evm/mocks"
	configtest "github.com/smartcontractkit/chainlink/v2/core/internal/testutils/configtest/v2"
	"github.com/smartcontractkit/chainlink/v2/core/services/chainlink"
)

func TestLegacyChains(t *testing.T) {
	evmCfg := configtest.NewGeneralConfig(t, func(c *chainlink.Config, s *chainlink.Secrets) {
	})

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
