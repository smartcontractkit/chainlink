package stellar

import (
	"testing"

	chainselectors "github.com/smartcontractkit/chain-selectors"
	"github.com/stellar/go-stellar-sdk/strkey"
	"github.com/stellar/go-stellar-sdk/xdr"
	"github.com/stretchr/testify/require"

	mcmsstellar "github.com/smartcontractkit/mcms/sdk/stellar"
	mcmstypes "github.com/smartcontractkit/mcms/types"

	"github.com/smartcontractkit/chainlink-stellar/bindings/scval"
)

func testContractAddress(t *testing.T, seed byte) string {
	t.Helper()

	raw := make([]byte, 32)
	for i := range raw {
		raw[i] = seed + byte(i)
	}

	addr, err := strkey.Encode(strkey.VersionByteContract, raw)
	require.NoError(t, err)

	return addr
}

func testAccountAddress(t *testing.T, seed byte) string {
	t.Helper()

	raw := make([]byte, 32)
	raw[0] = seed

	addr, err := strkey.Encode(strkey.VersionByteAccountID, raw)
	require.NoError(t, err)

	return addr
}

func TestForwarderSetConfigBatchOp(t *testing.T) {
	t.Parallel()

	selector := chainselectors.STELLAR_TESTNET.Selector
	forwarderAddr := testContractAddress(t, 1)
	signers := [][32]byte{{1}, {2}, {3}}

	batchOp, err := forwarderSetConfigBatchOp(selector, forwarderAddr, 5, 2, 1, signers)
	require.NoError(t, err)

	require.Equal(t, mcmstypes.ChainSelector(selector), batchOp.ChainSelector)
	require.Len(t, batchOp.Transactions, 1)

	tx := batchOp.Transactions[0]
	require.Equal(t, forwarderAddr, tx.To)
	require.Equal(t, string(ForwarderContract), tx.ContractType)
	require.Contains(t, string(tx.AdditionalFields), chainselectors.FamilyStellar)

	// Decode the payload and assert the exact function name and argument
	// order/values the on-chain forwarder expects; a swapped argument would
	// otherwise only fail at execution time.
	payload, err := mcmsstellar.DecodeSorobanInvokePayload(tx.Data)
	require.NoError(t, err)
	require.Equal(t, "set_config", payload.Function)
	require.Equal(t, []xdr.ScVal{
		scval.Uint32ToScVal(5),
		scval.Uint32ToScVal(2),
		scval.Uint32ToScVal(1),
		scval.Bytes32SliceToScVal(signers),
	}, payload.Args)
}

func TestForwarderSetConfigBatchOp_InvalidTarget(t *testing.T) {
	t.Parallel()

	_, err := forwarderSetConfigBatchOp(chainselectors.STELLAR_TESTNET.Selector, "not-an-address", 5, 2, 1, [][32]byte{{1}})
	require.Error(t, err)
}

func TestAddForwardersBatchOp(t *testing.T) {
	t.Parallel()

	selector := chainselectors.STELLAR_TESTNET.Selector
	forwarderAddr := testContractAddress(t, 1)
	transmitters := []string{
		testAccountAddress(t, 10),
		testAccountAddress(t, 20),
	}

	batchOp, err := addForwardersBatchOp(selector, forwarderAddr, transmitters)
	require.NoError(t, err)

	require.Equal(t, mcmstypes.ChainSelector(selector), batchOp.ChainSelector)
	require.Len(t, batchOp.Transactions, len(transmitters))

	for i, tx := range batchOp.Transactions {
		require.Equal(t, forwarderAddr, tx.To)
		require.Equal(t, string(ForwarderContract), tx.ContractType)

		payload, err := mcmsstellar.DecodeSorobanInvokePayload(tx.Data)
		require.NoError(t, err)
		require.Equal(t, "add_forwarder", payload.Function)
		require.Equal(t, []xdr.ScVal{scval.AddressToScVal(transmitters[i])}, payload.Args)
	}
}

func TestForwarderClearConfigBatchOp(t *testing.T) {
	t.Parallel()

	selector := chainselectors.STELLAR_TESTNET.Selector
	forwarderAddr := testContractAddress(t, 1)

	batchOp, err := forwarderClearConfigBatchOp(selector, forwarderAddr, 5, 2)
	require.NoError(t, err)

	require.Equal(t, mcmstypes.ChainSelector(selector), batchOp.ChainSelector)
	require.Len(t, batchOp.Transactions, 1)

	tx := batchOp.Transactions[0]
	require.Equal(t, forwarderAddr, tx.To)
	require.Equal(t, string(ForwarderContract), tx.ContractType)

	payload, err := mcmsstellar.DecodeSorobanInvokePayload(tx.Data)
	require.NoError(t, err)
	require.Equal(t, "clear_config", payload.Function)
	require.Equal(t, []xdr.ScVal{
		scval.Uint32ToScVal(5),
		scval.Uint32ToScVal(2),
	}, payload.Args)
}

func TestAddForwardersBatchOp_InvalidTarget(t *testing.T) {
	t.Parallel()

	_, err := addForwardersBatchOp(chainselectors.STELLAR_TESTNET.Selector, "not-an-address", []string{testAccountAddress(t, 10)})
	require.Error(t, err)
}
