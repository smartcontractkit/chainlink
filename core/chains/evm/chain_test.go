package evm_test

import (
	"math/big"
	"testing"

	"github.com/stretchr/testify/assert"

	"github.com/smartcontractkit/chainlink/v2/core/chains/evm"
	"github.com/smartcontractkit/chainlink/v2/core/chains/evm/mocks"
	configtest "github.com/smartcontractkit/chainlink/v2/core/internal/testutils/configtest/v2"
)

func TestLegacyChains(t *testing.T) {
	evmCfg := configtest.NewGeneralConfig(t, nil)

	c := mocks.NewChain(t)
	c.On("ID").Return(big.NewInt(7))
	m := map[string]evm.Chain{c.ID().String(): c}

	l := evm.NewLegacyChains(m, evmCfg.EVMConfigs())
	assert.NotNil(t, l.ChainNodeConfigs())
	got, err := l.Get(c.ID().String())
	assert.NoError(t, err)
	assert.Equal(t, c, got)

}
