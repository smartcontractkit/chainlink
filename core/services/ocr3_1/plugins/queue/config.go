package queue

import (
	"github.com/smartcontractkit/libocr/offchainreporting2plus/ocr3_1types"
)

const (
	maxObservedBlobHandles = 10   // sanity check on blobs in an observation
	maxObservedEvents      = 1000 // sanity check events in an observation

	//TODO reduce these? connect to cre settings?
	maxKeyValueValueBytes = ocr3_1types.MaxMaxKeyValueValueBytes
)

type PluginConfig struct {
	Limits ocr3_1types.ReportingPluginLimits `json:"limits,omitzero"`
}

func (c *PluginConfig) Validate() error {
	// TODO validate?
	return nil
}
