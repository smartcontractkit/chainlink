package aptostestutils

import (
	"testing"

	"github.com/aptos-labs/aptos-go-sdk"
	"github.com/aptos-labs/aptos-go-sdk/bcs"
	"github.com/stretchr/testify/require"
)

// FundAccount funds an Aptos account with the given amount of APT.
func FundAccount(t *testing.T, signer aptos.TransactionSigner, to aptos.AccountAddress, amount uint64, client aptos.AptosRpcClient) {
	toBytes, err := bcs.Serialize(&to)
	require.NoError(t, err)
	amountBytes, err := bcs.SerializeU64(amount)
	require.NoError(t, err)
	payload := aptos.TransactionPayload{Payload: &aptos.EntryFunction{
		Module: aptos.ModuleId{
			Address: aptos.AccountOne,
			Name:    "aptos_account",
		},
		Function: "transfer",
		Args: [][]byte{
			toBytes,
			amountBytes,
		},
	}}
	tx, err := client.BuildSignAndSubmitTransaction(signer, payload)
	require.NoError(t, err)
	res, err := client.WaitForTransaction(tx.Hash)
	require.NoError(t, err)
	require.True(t, res.Success, res.VmStatus)
	sender := signer.AccountAddress()
	t.Logf("Funded account %s from %s with %f APT", to.StringLong(), sender.StringLong(), float64(amount)/1e8)
}
