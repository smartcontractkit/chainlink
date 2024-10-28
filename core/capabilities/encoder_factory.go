package capabilities

import (
	"fmt"

	"github.com/smartcontractkit/chainlink-common/pkg/capabilities/consensus/ocr3"
	"github.com/smartcontractkit/chainlink-common/pkg/capabilities/consensus/ocr3/ocr3cap"
	"github.com/smartcontractkit/chainlink-common/pkg/capabilities/consensus/ocr3/types"
	"github.com/smartcontractkit/chainlink-common/pkg/logger"
	"github.com/smartcontractkit/chainlink-common/pkg/values"

	"github.com/smartcontractkit/chainlink/v2/core/services/relay/evm"
)

func NewEncoder(name string, config *values.Map, lggr logger.Logger) (types.Encoder, error) {
	switch ocr3cap.Encoder(name) {
	case ocr3cap.EncoderEVM:
		return evm.NewEVMEncoder(config)
	case ocr3cap.EncoderValueMap:
		return ocr3.ValueMapEncoder{}, nil
	// TODO: add a "no-op" encoder for users who only want to use dynamic ones?
	// https://smartcontract-it.atlassian.net/browse/CAPPL-88
	default:
		return nil, fmt.Errorf("encoder %s not supported", name)
	}
}
