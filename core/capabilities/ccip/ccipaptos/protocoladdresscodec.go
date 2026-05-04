package ccipaptos

import (
	"encoding/binary"
	"encoding/hex"

	chainsel "github.com/smartcontractkit/chain-selectors"

	ccipcommon "github.com/smartcontractkit/chainlink/v2/core/capabilities/ccip/common"
)

// ProtocolAddressCodec implements ProtocolAddressCodec for Aptos using stdlib only.
type ProtocolAddressCodec struct{}

var _ ccipcommon.ProtocolAddressCodec = ProtocolAddressCodec{}

func (ProtocolAddressCodec) OracleIDAsAddressBytes(oracleID uint8) ([]byte, error) {
	addr := make([]byte, 32)
	binary.BigEndian.PutUint32(addr[28:], uint32(oracleID))
	return addr, nil
}

func (ProtocolAddressCodec) TransmitterBytesToString(addr []byte) (string, error) {
	return hex.EncodeToString(addr), nil
}

func init() {
	ccipcommon.RegisterProtocolAddressCodec(chainsel.FamilyAptos, ProtocolAddressCodec{})
	// Sui and Aptos both use 32-byte Move addresses for these protocol-only methods.
	ccipcommon.RegisterProtocolAddressCodec(chainsel.FamilySui, ProtocolAddressCodec{})
}
