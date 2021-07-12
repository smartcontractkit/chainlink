package presenters

import (
	"testing"

	"github.com/ethereum/go-ethereum/common"
	"github.com/ethereum/go-ethereum/common/hexutil"
	"github.com/manyminds/api2go/jsonapi"
	"github.com/smartcontractkit/chainlink/core/assets"
	"github.com/smartcontractkit/chainlink/core/store/models"
	"github.com/smartcontractkit/chainlink/core/utils"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestEthTxResource(t *testing.T) {
	t.Parallel()

	from := common.HexToAddress("0x1")
	to := common.HexToAddress("0x2")
	tx := models.EthTx{
		ID:             1,
		EncodedPayload: []byte(`{"data": "is wilding out"}`),
		FromAddress:    from,
		ToAddress:      to,
		GasLimit:       uint64(5000),
		State:          models.EthTxConfirmed,
		Value:          assets.NewEthValue(1),
	}

	r := NewEthTxResource(tx)

	b, err := jsonapi.Marshal(r)
	require.NoError(t, err)

	expected := `
	{
		"data": {
		  "type": "transactions",
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
			"value": "0.000000000000000001"
		  }
		}
	  }
	`

	assert.JSONEq(t, expected, string(b))

	var (
		nonce           = int64(100)
		hash            = common.BytesToHash([]byte{1, 2, 3})
		gasPrice        = utils.NewBigI(1000)
		broadcastBefore = int64(300)
	)

	tx.Nonce = &nonce
	txa := models.EthTxAttempt{
		EthTx:                   tx,
		Hash:                    hash,
		GasPrice:                *gasPrice,
		SignedRawTx:             hexutil.MustDecode("0xcafe"),
		BroadcastBeforeBlockNum: &broadcastBefore,
	}

	r = NewEthTxResourceFromAttempt(txa)

	b, err = jsonapi.Marshal(r)
	require.NoError(t, err)

	expected = `
	{
		"data": {
		  "type": "transactions",
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
			"value": "0.000000000000000001"
		  }
		}
	  }
	`

	assert.JSONEq(t, expected, string(b))
}
