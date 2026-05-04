package common_test

import (
	"fmt"
	"testing"

	sel "github.com/smartcontractkit/chain-selectors"
	"github.com/stretchr/testify/require"

	ccipocr3common "github.com/smartcontractkit/chainlink-common/pkg/types/ccipocr3"

	ccipcommon "github.com/smartcontractkit/chainlink/v2/core/capabilities/ccip/common"
)

type testProtocolAddressCodec struct {
	prefix string
}

func (c testProtocolAddressCodec) OracleIDAsAddressBytes(oracleID uint8) ([]byte, error) {
	return []byte{oracleID}, nil
}

func (c testProtocolAddressCodec) TransmitterBytesToString(addr []byte) (string, error) {
	return fmt.Sprintf("%s-%x", c.prefix, addr), nil
}

func TestProtocolAddressCodecRegistryUsesSnapshot(t *testing.T) {
	ccipcommon.RegisterProtocolAddressCodec(sel.FamilyEVM, testProtocolAddressCodec{prefix: "first"})
	registry := ccipcommon.ProtocolAddressCodecs()
	ccipcommon.RegisterProtocolAddressCodec(sel.FamilyEVM, testProtocolAddressCodec{prefix: "second"})

	selector := ccipocr3common.ChainSelector(sel.ETHEREUM_MAINNET_OPTIMISM_1.Selector)

	oracleAddress, err := registry.OracleIDAsAddressBytes(7, selector)
	require.NoError(t, err)
	require.Equal(t, []byte{7}, oracleAddress)

	transmitter, err := registry.TransmitterBytesToString([]byte{0x42}, selector)
	require.NoError(t, err)
	require.Equal(t, "first-42", transmitter)
}
