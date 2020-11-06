package presenters

import (
	"encoding/json"
	"fmt"
	"testing"

	"github.com/ethereum/go-ethereum/common"
	"github.com/smartcontractkit/chainlink/core/store/models"

	"github.com/stretchr/testify/assert"
)

func Test_EthNewTx(t *testing.T) {
	t.Parallel()

	hash := common.BytesToHash([]byte{1, 2, 3})
	nonce := int64(100)
	broadcastBefore := int64(300)
	tx := models.EthTx{
		GasLimit: uint64(5000),
		Nonce:    &nonce,
	}
	txa := models.EthTxAttempt{
		EthTx:                   tx,
		Hash:                    hash,
		State:                   models.EthTxAttemptBroadcast,
		BroadcastBeforeBlockNum: &broadcastBefore,
	}
	ptx := NewEthTxFromAttempt(&txa)

	assert.Equal(t, hash, ptx.Hash)
	assert.Equal(t, "5000", ptx.GasLimit)
	assert.Equal(t, "100", ptx.Nonce)
	assert.Equal(t, "300", ptx.SentAt)
}

func TestTx_MarshalJSON(t *testing.T) {
	t.Parallel()
	hex := "0xfa85d5aa5c48e23b40f5a75d62adfc8036330f9bf86c601229e2bc63e1331d3c"
	want := fmt.Sprintf("{\"hash\":\"%s\"}", hex)

	tx := EthTx{Hash: common.HexToHash(hex)}
	b, err := json.Marshal(&tx)
	assert.NoError(t, err)
	assert.Equal(t, want, string(b))
}
