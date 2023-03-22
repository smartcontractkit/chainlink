package types

import (
	"github.com/smartcontractkit/chainlink/core/services"
)

type ForwarderManager[ADDR any] interface {
	services.ServiceCtx
	// Name change b/c no distinction b/w EOA and contracts for some chains
	GetForwarderFor(addr ADDR) (forwarder ADDR, err error)
	GetForwardedPayload(dest ADDR, origPayload []byte) ([]byte, error)
}
