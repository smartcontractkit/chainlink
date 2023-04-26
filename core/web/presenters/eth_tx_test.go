package presenters

import (
	"testing"

	"github.com/ethereum/go-ethereum/common"
	"github.com/ethereum/go-ethereum/common/hexutil"
	"github.com/manyminds/api2go/jsonapi"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"

	"github.com/smartcontractkit/chainlink/v2/core/assets"
	"github.com/smartcontractkit/chainlink/v2/core/chains/evm/txmgr"
)

func TestEthTxResource(t *testing.T) {
	t.Parallel()

	tx := txmgr.EvmTx{
		ID:             1,
		EncodedPayload: []byte(`{"data": "is wilding out"}`),
		FromAddress:    common.HexToAddress("0x1"),
		ToAddress:      common.HexToAddress("0x2"),
		GasLimit:       uint32(5000),
		State:          txmgr.EthTxConfirmed,
		Value:          assets.NewEthValue(1),
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
			"evmChainID": "0"
		  }
		}
	  }
	`

	assert.JSONEq(t, expected, string(b))

	var (
		nonce           = int64(100)
		hash            = common.BytesToHash([]byte{1, 2, 3})
		gasPrice        = assets.NewWeiI(1000)
		broadcastBefore = int64(300)
	)

	tx.Nonce = &nonce
	txa := txmgr.EvmTxAttempt{
		EthTx:                   tx,
		Hash:                    hash,
		GasPrice:                gasPrice,
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
		  "id": "0x0000000000000000000000000000000000000000000000000000000000010203",
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
			"evmChainID": "0"
		  }
		}
	  }
	`

	assert.JSONEq(t, expected, string(b))
}
