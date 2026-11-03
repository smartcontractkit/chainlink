package pipeline

import (
	"testing"

	"github.com/ethereum/go-ethereum/common"
	"github.com/stretchr/testify/require"
)

func TestETHTxFromGetters_LiteralAddress(t *testing.T) {
	t.Parallel()

	const address = "0x882969652440ccf14a5dbb9bd53eb21cb1e11e5c"
	var fromAddrs AddressSliceParam

	err := ResolveParam(
		&fromAddrs,
		ethTxFromGetters(address, NewVarsFrom(nil)),
	)

	require.NoError(t, err)
	require.Equal(t, AddressSliceParam{common.HexToAddress(address)}, fromAddrs)
}
