package evm

import (
	"strings"
	"testing"

	"github.com/ethereum/go-ethereum/common"
	"github.com/stretchr/testify/require"
)

func TestAddressCollection_ToStrings(t *testing.T) {
	t.Parallel()

	hex1 := "0xaAaAaAaaAaAaAaaAaAAAAAAAAaaaAaAaAaaAaaAa"
	hex2 := "0xbBbBBBBbbBBBbbbBbbBbbbbBBbBbbbbBbBbbBBbB"

	ac := AddressCollection{
		common.HexToAddress(hex1),
		common.HexToAddress(hex2),
	}

	acStrings := ac.ToStrings()
	require.Len(t, acStrings, 2)
	require.Equal(t, hex1, acStrings[0])
	require.Equal(t, hex2, acStrings[1])
}

func TestAddressCollection_Scan_Value(t *testing.T) {
	t.Parallel()

	ac := AddressCollection{
		common.HexToAddress(strings.Repeat("AA", 20)),
		common.HexToAddress(strings.Repeat("BB", 20)),
	}

	val, err := ac.Value()
	require.NoError(t, err)

	var acNew AddressCollection
	err = acNew.Scan(val)
	require.NoError(t, err)

	require.Equal(t, ac, acNew)
}
