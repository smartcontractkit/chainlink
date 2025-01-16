package types

import (
	"context"

	"github.com/smartcontractkit/chainlink-common/pkg/services"

	"github.com/smartcontractkit/chainlink-framework/chains"
)

type ForwarderManager[ADDR chains.Hashable] interface {
	services.Service
	ForwarderFor(ctx context.Context, addr ADDR) (forwarder ADDR, err error)
	ForwarderForOCR2Feeds(ctx context.Context, eoa, ocr2Aggregator ADDR) (forwarder ADDR, err error)
	// Converts payload to be forwarder-friendly
	ConvertPayload(dest ADDR, origPayload []byte) ([]byte, error)
}
