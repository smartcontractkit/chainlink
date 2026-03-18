package readcontract

import (
	"testing"

	"github.com/stretchr/testify/require"

	chainselectors "github.com/smartcontractkit/chain-selectors"

	"github.com/smartcontractkit/chainlink/system-tests/lib/cre"
)

type familyMatcherStub struct {
	family string
}

func (s familyMatcherStub) IsFamily(chainFamily string) bool {
	return s.family == chainFamily
}

func TestShouldSkipPostEnvStartup(t *testing.T) {
	t.Run("skips aptos when write aptos feature owns the don", func(t *testing.T) {
		don := &cre.Don{Flags: []cre.CapabilityFlag{cre.ReadContractCapability, cre.WriteAptosCapability}}
		require.True(t, shouldSkipPostEnvStartup(don, familyMatcherStub{family: chainselectors.FamilyAptos}))
	})

	t.Run("does not skip aptos for read-only dons", func(t *testing.T) {
		don := &cre.Don{Flags: []cre.CapabilityFlag{cre.ReadContractCapability}}
		require.False(t, shouldSkipPostEnvStartup(don, familyMatcherStub{family: chainselectors.FamilyAptos}))
	})

	t.Run("does not skip non-aptos chains", func(t *testing.T) {
		don := &cre.Don{Flags: []cre.CapabilityFlag{cre.ReadContractCapability, cre.WriteAptosCapability}}
		require.False(t, shouldSkipPostEnvStartup(don, familyMatcherStub{family: chainselectors.FamilyEVM}))
	})
}
