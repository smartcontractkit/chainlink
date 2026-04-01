package cre

import (
	"testing"

	"github.com/stretchr/testify/require"

	chainselectors "github.com/smartcontractkit/chain-selectors"
)

func TestOCRExtraSignerFamiliesForFamily(t *testing.T) {
	require.Equal(t, []string{chainselectors.FamilyAptos}, OCRExtraSignerFamiliesForFamily(chainselectors.FamilyAptos))
	require.Equal(t, []string{chainselectors.FamilySolana}, OCRExtraSignerFamiliesForFamily(chainselectors.FamilySolana))
	require.Nil(t, OCRExtraSignerFamiliesForFamily(chainselectors.FamilyEVM))
}

func TestCapabilityToExtraSignerFamiliesCopiesInput(t *testing.T) {
	families := []string{chainselectors.FamilyAptos}
	got := CapabilityToExtraSignerFamilies(families, "cap-a", "cap-b")
	require.Equal(t, map[string][]string{
		"cap-a": {chainselectors.FamilyAptos},
		"cap-b": {chainselectors.FamilyAptos},
	}, got)

	families[0] = chainselectors.FamilySolana
	require.Equal(t, []string{chainselectors.FamilyAptos}, got["cap-a"])
	require.Equal(t, []string{chainselectors.FamilyAptos}, got["cap-b"])
}
