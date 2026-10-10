package queue

import (
	"github.com/smartcontractkit/libocr/offchainreporting2plus/ocr3_1types"
)

type PluginConfig struct {
	Limits ocr3_1types.ReportingPluginLimits `json:"limits,omitzero"`
}

func (c *PluginConfig) Validate() error {
	//TODO validate?
	return nil
}
