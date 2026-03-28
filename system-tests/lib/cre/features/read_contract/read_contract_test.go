package readcontract

import (
	"testing"

	"github.com/stretchr/testify/require"

	chainselectors "github.com/smartcontractkit/chain-selectors"

	"github.com/smartcontractkit/chainlink/system-tests/lib/cre"
	aptosfeature "github.com/smartcontractkit/chainlink/system-tests/lib/cre/features/aptos"
)

type blockchainOutputStub struct {
	chainSelector uint64
	chainFamily   string
}

func (s blockchainOutputStub) ChainSelector() uint64 {
	return s.chainSelector
}

func (s blockchainOutputStub) ChainFamily() string {
	return s.chainFamily
}

func TestAptosCapabilityLabel(t *testing.T) {
	bc := blockchainOutputStub{chainSelector: 1, chainFamily: chainselectors.FamilyAptos}

	t.Run("skips aptos when write aptos feature owns the don", func(t *testing.T) {
		don := &cre.DonMetadata{Flags: []string{cre.ReadContractCapability, cre.WriteAptosCapability}}
		label, skip, err := aptosCapabilityLabel(don, bc)
		require.NoError(t, err)
		require.Empty(t, label)
		require.True(t, skip)
	})

	t.Run("uses aptos label for read-only dons", func(t *testing.T) {
		don := &cre.DonMetadata{Flags: []string{cre.ReadContractCapability}}
		label, skip, err := aptosCapabilityLabel(don, bc)
		require.NoError(t, err)
		require.Equal(t, aptosfeature.CapabilityLabel(1), label)
		require.False(t, skip)
	})
}
