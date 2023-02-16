package soltxm

import (
	"testing"

	"github.com/gagliardetto/solana-go"
	"github.com/gagliardetto/solana-go/rpc"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestSortSignaturesAndResults(t *testing.T) {
	sig := []solana.Signature{
		{0}, {1}, {2}, {3},
	}

	statuses := []*rpc.SignatureStatusesResult{
		{ConfirmationStatus: rpc.ConfirmationStatusProcessed},
		{ConfirmationStatus: rpc.ConfirmationStatusConfirmed},
		nil,
		{ConfirmationStatus: rpc.ConfirmationStatusConfirmed, Err: "ERROR"},
	}

	_, _, err := SortSignaturesAndResults([]solana.Signature{}, statuses)
	require.Error(t, err)

	sig, statuses, err = SortSignaturesAndResults(sig, statuses)
	require.NoError(t, err)

	// new expected order [1, 3, 0, 2]
	assert.Equal(t, rpc.SignatureStatusesResult{ConfirmationStatus: rpc.ConfirmationStatusConfirmed}, *statuses[0])
	assert.Equal(t, rpc.SignatureStatusesResult{ConfirmationStatus: rpc.ConfirmationStatusConfirmed, Err: "ERROR"}, *statuses[1])
	assert.Equal(t, rpc.SignatureStatusesResult{ConfirmationStatus: rpc.ConfirmationStatusProcessed}, *statuses[2])
	assert.True(t, nil == statuses[3])

	assert.Equal(t, solana.Signature{1}, sig[0])
	assert.Equal(t, solana.Signature{3}, sig[1])
	assert.Equal(t, solana.Signature{0}, sig[2])
	assert.Equal(t, solana.Signature{2}, sig[3])
}
