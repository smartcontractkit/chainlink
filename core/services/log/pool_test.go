package log

import (
	"math/big"
	"testing"

	"github.com/ethereum/go-ethereum/common"
	"github.com/ethereum/go-ethereum/core/types"
	"github.com/stretchr/testify/require"
)

func Test_logPool_addLog(t *testing.T) {
	p := newLogPool()
	l := types.Log{BlockHash: common.BigToHash(big.NewInt(123456)), Index: 42}
	p.addLog(l)
	require.Len(t, p.logsByBlockHash[l.BlockHash], 1)
	p.addLog(l)
	require.Len(t, p.logsByBlockHash[l.BlockHash], 1)
}
