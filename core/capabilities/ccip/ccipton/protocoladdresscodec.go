package ccipton

import (
	"encoding/binary"
	"encoding/hex"

	chainsel "github.com/smartcontractkit/chain-selectors"

	ccipcommon "github.com/smartcontractkit/chainlink/v2/core/capabilities/ccip/common"
)

// tonRawAddressLength is 4 bytes workchain plus 32 bytes address data.
const tonRawAddressLength = 36

// ProtocolAddressCodec implements ProtocolAddressCodec for TON using stdlib only.
type ProtocolAddressCodec struct{}

var _ ccipcommon.ProtocolAddressCodec = ProtocolAddressCodec{}

func (ProtocolAddressCodec) OracleIDAsAddressBytes(oracleID uint8) ([]byte, error) {
	raw := make([]byte, tonRawAddressLength)
	binary.BigEndian.PutUint32(raw[4:], uint32(oracleID))
	return raw, nil
}

func (ProtocolAddressCodec) TransmitterBytesToString(addr []byte) (string, error) {
	return hex.EncodeToString(addr), nil
}

func init() {
	ccipcommon.RegisterProtocolAddressCodec(chainsel.FamilyTon, ProtocolAddressCodec{})
}
