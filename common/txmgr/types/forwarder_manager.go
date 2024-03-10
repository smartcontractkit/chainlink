package types

import (
	"github.com/smartcontractkit/chainlink-common/pkg/services"
	"github.com/smartcontractkit/chainlink/v2/common/types"
)

//go:generate mockery --quiet --name ForwarderManager --output ./mocks/ --case=underscore
type ForwarderManager[ADDR types.Hashable] interface {
	services.Service
	ForwarderFor(addr ADDR) (forwarder ADDR, err error)
	// Converts payload to be forwarder-friendly
	ConvertPayload(dest ADDR, origPayload []byte) ([]byte, error)
}
