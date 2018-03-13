package store_test

import (
	"math/big"
	"testing"

	"github.com/ethereum/go-ethereum/common"
	strpkg "github.com/smartcontractkit/chainlink/store"
	"github.com/stretchr/testify/assert"
)

func TestStore_ConfigDefaults(t *testing.T) {
	config := strpkg.NewConfig()
	assert.Equal(t, uint64(0), config.ChainID)
	assert.Equal(t, *big.NewInt(20000000000), config.EthGasPriceDefault)
	assert.Equal(t, "0x514910771AF9Ca656af840dff83E8264EcF986CA", common.HexToAddress(config.LinkContractAddress).String())
}
