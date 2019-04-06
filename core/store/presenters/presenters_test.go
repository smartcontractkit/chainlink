package presenters

import (
	"testing"

	"github.com/smartcontractkit/chainlink/core/store/models"
	"github.com/stretchr/testify/assert"
)

func Test_NewTx(t *testing.T) {
	t.Parallel()

	tx := models.Tx{
		GasLimit: uint64(5000),
		Nonce:    uint64(100),
		SentAt:   uint64(300),
	}
	ptx := NewTx(&tx)

	assert.Equal(t, "5000", ptx.GasLimit)
	assert.Equal(t, "100", ptx.Nonce)
	assert.Equal(t, "300", ptx.SentAt)
}
