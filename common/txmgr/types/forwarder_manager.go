package types

import (
	"github.com/smartcontractkit/chainlink/v2/core/services"
)

type ForwarderManager[ADDR any] interface {
	services.ServiceCtx
	ForwarderFor(addr ADDR) (forwarder ADDR, err error)
	// Converts payload to be forwarder-friendly
	ConvertPayload(dest ADDR, origPayload []byte) ([]byte, error)
}
