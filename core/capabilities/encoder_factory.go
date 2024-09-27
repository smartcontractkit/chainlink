package capabilities

import (
	"fmt"

	"github.com/smartcontractkit/chainlink-common/pkg/capabilities/consensus/ocr3/types"
	"github.com/smartcontractkit/chainlink-common/pkg/logger"
	"github.com/smartcontractkit/chainlink-common/pkg/values"
	"github.com/smartcontractkit/chainlink/v2/core/services/relay/evm"
)

func NewEncoder(name string, config *values.Map, lggr logger.Logger) (types.Encoder, error) {
	switch name {
	case "EVM":
		return evm.NewEVMEncoder(config)
	// TODO: add a "no-op" encoder for users who only want to use dynamic ones?
	default:
		return nil, fmt.Errorf("encoder %s not supported", name)
	}
}
