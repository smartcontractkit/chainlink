package txmgr_test

import (
	"bytes"
	"math"
	"math/big"
	"testing"

	gethtypes "github.com/ethereum/go-ethereum/core/types"
	"github.com/ethereum/go-ethereum/crypto"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"

	"github.com/smartcontractkit/chainlink-evm/pkg/utils"
	"github.com/smartcontractkit/chainlink/v2/core/chains/evm/txmgr"
)

func TestTx_GetID(t *testing.T) {
	tx := txmgr.Tx{ID: math.MinInt64}
	assert.Equal(t, "-9223372036854775808", tx.GetID())
}

func TestGetGethSignedTx(t *testing.T) {
	chainID := big.NewInt(3)
	signer := gethtypes.NewCancunSigner(chainID)
	to := utils.NewAddress()
	key, err := crypto.GenerateKey()
	require.NoError(t, err)
	signedTx := gethtypes.MustSignNewTx(key, signer, &gethtypes.LegacyTx{
		Nonce:    42,
		To:       &to,
		Value:    big.NewInt(142),
		Gas:      242,
		GasPrice: big.NewInt(342),
		Data:     []byte{1, 2, 3},
	})
	rlp := new(bytes.Buffer)
	require.NoError(t, signedTx.EncodeRLP(rlp))

	signedRawTx := rlp.Bytes()

	gotSignedTx, err := txmgr.GetGethSignedTx(signedRawTx)
	require.NoError(t, err)
	decodedEncoded := new(bytes.Buffer)
	require.NoError(t, gotSignedTx.EncodeRLP(decodedEncoded))

	require.Equal(t, signedTx.Hash(), gotSignedTx.Hash())
	require.Equal(t, signedRawTx, decodedEncoded.Bytes())
}
