package presenters

import (
	"math/big"
	"testing"

	"github.com/ethereum/go-ethereum/common"
	"github.com/ethereum/go-ethereum/common/hexutil"
	"github.com/manyminds/api2go/jsonapi"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"

	"github.com/smartcontractkit/chainlink-evm/pkg/assets"
	"github.com/smartcontractkit/chainlink-evm/pkg/gas"
	evmtypes "github.com/smartcontractkit/chainlink-evm/pkg/types"
	txmgrcommon "github.com/smartcontractkit/chainlink-framework/chains/txmgr"
	"github.com/smartcontractkit/chainlink/v2/core/chains/evm/txmgr"
)

func TestEthTxResource(t *testing.T) {
	t.Parallel()

	chainID := big.NewInt(54321)
	tx := txmgr.Tx{
		ID:             1,
		EncodedPayload: []byte(`{"data": "is wilding out"}`),
		FromAddress:    common.HexToAddress("0x1"),
		ToAddress:      common.HexToAddress("0x2"),
		FeeLimit:       uint64(5000),
		ChainID:        chainID,
		State:          txmgrcommon.TxConfirmed,
		Value:          big.Int(assets.NewEthValue(1)),
	}

	r := NewEthTxResource(tx)

	b, err := jsonapi.Marshal(r)
	require.NoError(t, err)

	expected := `
	{
		"data": {
		  "type": "evm_transactions",
		  "id": "",
		  "attributes": {
			"state": "confirmed",
			"data": "0x7b2264617461223a202269732077696c64696e67206f7574227d",
			"from": "0x0000000000000000000000000000000000000001",
			"gasLimit": "5000",
			"gasPrice": "",
			"hash": "0x0000000000000000000000000000000000000000000000000000000000000000",
			"rawHex": "",
			"nonce": "",
			"sentAt": "",
			"to": "0x0000000000000000000000000000000000000002",
			"value": "0.000000000000000001",
			"evmChainID": "54321"
		  }
		}
	  }
	`

	assert.JSONEq(t, expected, string(b))

	var (
		nonce           = evmtypes.Nonce(100)
		hash            = common.BytesToHash([]byte{1, 2, 3})
		gasPrice        = assets.NewWeiI(1000)
		broadcastBefore = int64(300)
	)

	tx.Sequence = &nonce
	txa := txmgr.TxAttempt{
		Tx:                      tx,
		Hash:                    hash,
		TxFee:                   gas.EvmFee{GasPrice: gasPrice},
		SignedRawTx:             hexutil.MustDecode("0xcafe"),
		BroadcastBeforeBlockNum: &broadcastBefore,
	}

	r = NewEthTxResourceFromAttempt(txa)

	b, err = jsonapi.Marshal(r)
	require.NoError(t, err)

	expected = `
	{
		"data": {
		  "type": "evm_transactions",
		  "id": "54321/0x0000000000000000000000000000000000000000000000000000000000010203",
		  "attributes": {
			"state": "confirmed",
			"data": "0x7b2264617461223a202269732077696c64696e67206f7574227d",
			"from": "0x0000000000000000000000000000000000000001",
			"gasLimit": "5000",
			"gasPrice": "1000",
			"hash": "0x0000000000000000000000000000000000000000000000000000000000010203",
			"rawHex": "0xcafe",
			"nonce": "100",
			"sentAt": "300",
			"to": "0x0000000000000000000000000000000000000002",
			"value": "0.000000000000000001",
			"evmChainID": "54321"
		  }
		}
	  }
	`

	assert.JSONEq(t, expected, string(b))
}
