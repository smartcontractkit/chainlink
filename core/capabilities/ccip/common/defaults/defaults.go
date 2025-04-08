package defaults

import (
	"github.com/smartcontractkit/chainlink/v2/core/capabilities/ccip/ccipevm"
	"github.com/smartcontractkit/chainlink/v2/core/capabilities/ccip/ccipsolana"
	"github.com/smartcontractkit/chainlink/v2/core/capabilities/ccip/common"
)

// DefaultExtraDataCodec is the default ExtraDataCodec for CCIP initialized with all supported chain families.
var DefaultExtraDataCodec = common.NewExtraDataCodec(
	common.ExtraDataCodecParams{
		EVMExtraDataDecoder:    ccipevm.ExtraDataDecoder{},
		SolanaExtraDataDecoder: ccipsolana.ExtraDataDecoder{},
	})

// DefaultAddressCodec is the default AddressCodec for CCIP initialized with all supported chain families.
var DefaultAddressCodec = common.NewAddressCodec(
	common.AddressCodecParams{
		EVMAddressCodec:    ccipevm.AddressCodec{},
		SolanaAddressCodec: ccipsolana.AddressCodec{},
	})
