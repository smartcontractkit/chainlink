package adapters

import (
	"chainlink/core/store/models"
	"github.com/stretchr/testify/assert"
	"testing"
)

// In pathological cases, the receipt can be nil.
// Need to ensure we don't panic in this case and return errored output instead
func TestEthTxAdapter_addReceiptToResult(t *testing.T) {
	t.Parallel()

	j := models.JSON{}
	input := *models.NewRunInput(models.NewID(), j, models.RunStatusUnstarted)

	output := addReceiptToResult(nil, input, j)
	assert.True(t, output.HasError())
	assert.EqualError(t, output.Error(), "missing receipt for transaction")
}
