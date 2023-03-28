package types

import (
	"github.com/smartcontractkit/chainlink/common/types"
	"github.com/smartcontractkit/chainlink/core/services"
)

type ForwarderManager[ADDR types.Hashable] interface {
	services.ServiceCtx
	ForwarderFor(addr ADDR) (forwarder ADDR, err error)
	// Converts payload to be forwarder-friendly
	ConvertPayload(dest ADDR, origPayload []byte) ([]byte, error)
}
