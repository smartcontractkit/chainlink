package adapters_test

import (
	"encoding/hex"
	"fmt"
	"math/big"
	"strings"
	"testing"

	"github.com/ethereum/go-ethereum/accounts/abi"

	"github.com/tidwall/gjson"

	"github.com/ethereum/go-ethereum/common"
	"github.com/ethereum/go-ethereum/common/hexutil"

	"github.com/smartcontractkit/chainlink/core/adapters"
	"github.com/smartcontractkit/chainlink/core/internal/cltest"
	"github.com/smartcontractkit/chainlink/core/store/models"
	"github.com/smartcontractkit/chainlink/core/utils"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

// var adapterUnderTest = adapters.EthTxEncode{
// 	Address: common.HexToAddress(
// 		"0xdeadbeefdeadbeefdeadbeefdeadbeefdeadbeef"),
// 	MethodName: "verifyVRFProof",
// 	Types: map[string]string{
// 		"gammaX": "uint256", "gammaY": "uint256", "c": "uint256", "s": "uint256"},
// 	Order:    []string{"gammaX", "gammaY", "c", "s"},
// 	GasPrice: models.NewBig(big.NewInt(1 << 44)), // ~20k Gwei
// 	GasLimit: 500000,
// }

func TestEthTxEncodeAdapter_Perform_ConfirmedWithJSON(t *testing.T) {
	uint256Type, err := abi.NewType("uint256", []abi.ArgumentMarshaling{})
	var adapterUnderTest = adapters.EthTxEncode{
		Address: common.HexToAddress(
			"0xdeadbeefdeadbeefdeadbeefdeadbeefdeadbeef"),
		FunctionABI: abi.Method{
			Name:    "verifyVRFProof",
			RawName: "verifyVRFProof",
			Const:   false,
			Inputs: []abi.Argument{
				abi.Argument{
					Name: "gammaX",
					Type: uint256Type,
				},
				abi.Argument{
					Name: "gammaY",
					Type: uint256Type,
				},
				abi.Argument{
					Name: "c",
					Type: uint256Type,
				},
				abi.Argument{
					Name: "s",
					Type: uint256Type,
				},
			},
			Outputs: []abi.Argument{},
		},
		GasPrice: models.NewBig(big.NewInt(1 << 44)), // ~20k Gwei
		GasLimit: 500000,
	}

	rawInput := `
    {
       "result": {
        "gammaX":
          "0xa2e03a05b089db7b79cd0f6655d6af3e2d06bd0129f87f9f2155612b4e2a41d8",
        "gammaY":
          "0x0a1dadcabf900bdfb6484e9a4390bffa6ccd666a565a991f061faf868cc9fce8",
        "c":
          "0xf82b4f9161ab41ae7c11e7deb628024ef9f5e9a0bca029f0ccb5cb534c70be31",
        "s":
          "0x2b1049accb1596a24517f96761b22600a690ee5c6b6cadae3fa522e7d95ba338"
       }
    }
`
	require.True(t, gjson.Valid(rawInput), "invalid result json: %s", rawInput)
	inputValue := gjson.Parse(rawInput).Get("result").Value().(map[string]interface{})

	types := []string{}
	for _, input := range adapterUnderTest.FunctionABI.Inputs {
		types = append(types, input.Type.String())
	}
	fullSignatureHash := utils.MustHash(fmt.Sprintf(
		"%s(%s)", adapterUnderTest.FunctionABI.RawName, strings.Join(types, ",")))
	selector := []string{"0x" + hex.EncodeToString(fullSignatureHash[:4])}

	values := []string{}
	for _, input := range adapterUnderTest.FunctionABI.Inputs {
		values = append(values, inputValue[input.Name].(string)[2:])
	}

	expectedAsHex := strings.Join(append(selector, values...), "")
	t.Parallel()
	app, cleanup := cltest.NewApplicationWithKey(t)
	defer cleanup()
	store := app.Store

	ethMock, err := app.MockStartAndConnect()
	require.NoError(t, err)

	hash := cltest.NewHash()
	sentAt := uint64(23456)
	confirmed := sentAt + 1
	ethMock.Register("eth_sendRawTransaction", hash,
		func(_ interface{}, data ...interface{}) error {
			rlp := data[0].([]interface{})[0].(string)
			tx, err := utils.DecodeEthereumTx(rlp)
			assert.NoError(t, err)
			assert.Equal(t, adapterUnderTest.Address.String(), tx.To().String())
			assert.Equal(t, expectedAsHex, hexutil.Encode(tx.Data()))
			return nil
		})
	receipt := models.TxReceipt{Hash: hash, BlockNumber: cltest.Int(confirmed)}
	ethMock.Register("eth_getTransactionReceipt", receipt)
	input := cltest.RunResultWithData(rawInput)
	responseData := adapterUnderTest.Perform(input, store)
	assert.False(t, responseData.HasError())
	from := cltest.GetAccountAddress(t, store)
	assert.NoError(t, err)
	ethMock.EventuallyAllCalled(t)
	txs, err := store.TxFrom(from)
	require.Len(t, txs, 1)
	assert.Len(t, txs[0].Attempts, 1)
}
