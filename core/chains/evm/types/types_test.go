package types_test

import (
	"math/rand"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"

	"github.com/smartcontractkit/chainlink/core/chains/evm/types"
	"github.com/smartcontractkit/chainlink/core/internal/cltest"
	"github.com/smartcontractkit/chainlink/core/internal/testutils/pgtest"
	"github.com/smartcontractkit/chainlink/core/utils"
)

func Test_PersistsReadsChain(t *testing.T) {
	db := pgtest.NewGormDB(t)

	val := utils.NewBigI(rand.Int63())
	addr := cltest.NewAddress()
	ks := make(map[string]types.ChainCfg)
	ks[addr.Hex()] = types.ChainCfg{EvmMaxGasPriceWei: val}
	chain := types.Chain{
		ID: *utils.NewBigI(rand.Int63()),
		Cfg: types.ChainCfg{
			KeySpecific: ks,
		},
	}

	require.NoError(t, db.Create(&chain).Error)

	var loadedChain types.Chain
	require.NoError(t, db.First(&loadedChain, "id = ?", chain.ID).Error)

	loadedVal := loadedChain.Cfg.KeySpecific[addr.Hex()].EvmMaxGasPriceWei
	assert.Equal(t, loadedVal, val)
}
