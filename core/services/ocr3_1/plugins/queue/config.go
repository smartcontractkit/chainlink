package queue

import (
	"github.com/smartcontractkit/libocr/offchainreporting2plus/ocr3_1types"
)

const (
	maxObservedBlobHandles = 10   // sanity check on blobs in an observation
	maxObservedEvents      = 1000 // sanity check events in an observation

	//TODO reduce these? connect to cre settings?
	maxObservationBytes                    = ocr3_1types.MaxMaxObservationBytes
	maxKeyValueValueBytes                  = ocr3_1types.MaxMaxKeyValueKeyBytes
	maxKeyValueModifiedKeysPlusValuesBytes = ocr3_1types.MaxMaxKeyValueModifiedKeysPlusValuesBytes
	maxBlobPayloadBytes                    = ocr3_1types.MaxMaxBlobPayloadBytes
)

type PluginConfig struct {
	Limits ocr3_1types.ReportingPluginLimits `json:"limits,omitzero"`
}

func (c *PluginConfig) Validate() error {
	//TODO validate?
	return nil
}
