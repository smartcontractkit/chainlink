package queue

import (
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"

	"github.com/smartcontractkit/libocr/offchainreporting2plus/ocr3types"
	ocrtypes "github.com/smartcontractkit/libocr/offchainreporting2plus/types"
)

func TestNoopTransmitter(t *testing.T) {
	t.Parallel()

	tr := NewTransmitter("account-1")
	acct, err := tr.FromAccount(t.Context())
	require.NoError(t, err)
	assert.Equal(t, ocrtypes.Account("account-1"), acct)
	require.NoError(t, tr.Transmit(t.Context(), ocrtypes.ConfigDigest{}, 1, ocr3types.ReportWithInfo[[]byte]{}, nil))
}
