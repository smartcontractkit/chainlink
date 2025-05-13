package defaults

import (
	chainsel "github.com/smartcontractkit/chain-selectors"

	"github.com/smartcontractkit/chainlink/v2/core/capabilities/ccip/ccipevm"
	"github.com/smartcontractkit/chainlink/v2/core/capabilities/ccip/ccipsolana"
	"github.com/smartcontractkit/chainlink/v2/core/capabilities/ccip/common"
)

// DefaultExtraDataCodec is the default ExtraDataCodec for CCIP initialized with all supported chain families.
var DefaultExtraDataCodec = common.NewExtraDataCodec(ccipevm.ExtraDataCodec{}, ccipsolana.ExtraDataCodec{})

// DefaultAddressCodec is the default AddressCodec for CCIP initialized with all supported chain families.
var DefaultAddressCodec = common.NewAddressCodec(ccipevm.AddressCodec{}, ccipsolana.AddressCodec{})

var DefaultCRCW = common.NewCRCW(map[string]common.ChainRWProvider{
	chainsel.FamilyEVM:    ccipevm.ChainCWProvider{},
	chainsel.FamilySolana: ccipsolana.ChainRWProvider{},
})
