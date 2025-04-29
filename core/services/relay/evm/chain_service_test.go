package evm

import (
	"math/big"
	"testing"
	"time"

	"github.com/ethereum/go-ethereum/common"
	"github.com/stretchr/testify/require"

	gethtypes "github.com/ethereum/go-ethereum/core/types"
	evmtypes "github.com/smartcontractkit/chainlink-evm/pkg/types"
)

func TestConverters(t *testing.T) {
	t.Parallel()

	t.Run("convert head", func(t *testing.T) {
		head := evmtypes.Head{
			Timestamp: time.Unix(100000, 100),
			Number:    100,
			Hash:      common.HexToHash("0x123"),
		}
		result := convertHead(&head)
		require.Equal(t, result.Hash, head.Hash.Bytes())
	})

	t.Run("convert transaction", func(t *testing.T) {
		tx := gethtypes.NewTransaction(
			1,
			common.HexToAddress("0xabc123"),
			big.NewInt(1000),
			21000,
			big.NewInt(1e9),
			[]byte{1, 2, 3},
		)

		result := convertTransaction(tx)
		require.NotNil(t, result)
		require.Equal(t, tx.Hash().Hex(), result.Hash)
		require.Equal(t, tx.Nonce(), result.Nonce)
		require.Equal(t, tx.Gas(), result.Gas)
		require.Equal(t, tx.GasPrice(), result.GasPrice)
		require.Equal(t, tx.Value(), result.Value)
		require.Equal(t, tx.To().Hex(), result.To)
		require.Equal(t, tx.Data(), result.Data)
	})
}
