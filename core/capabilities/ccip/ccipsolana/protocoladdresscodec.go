package ccipsolana

import (
	"encoding/binary"

	"github.com/mr-tron/base58"

	chainsel "github.com/smartcontractkit/chain-selectors"

	ccipcommon "github.com/smartcontractkit/chainlink/v2/core/capabilities/ccip/common"
)

const svmPublicKeyLength = 32

// ProtocolAddressCodec implements ProtocolAddressCodec for Solana using base58 only.
type ProtocolAddressCodec struct{}

var _ ccipcommon.ProtocolAddressCodec = ProtocolAddressCodec{}

func (ProtocolAddressCodec) OracleIDAsAddressBytes(oracleID uint8) ([]byte, error) {
	addr := make([]byte, svmPublicKeyLength)
	binary.LittleEndian.PutUint32(addr, uint32(oracleID))
	return addr, nil
}

func (ProtocolAddressCodec) TransmitterBytesToString(addr []byte) (string, error) {
	return base58.Encode(addr), nil
}

func init() {
	ccipcommon.RegisterProtocolAddressCodec(chainsel.FamilySolana, ProtocolAddressCodec{})
}
