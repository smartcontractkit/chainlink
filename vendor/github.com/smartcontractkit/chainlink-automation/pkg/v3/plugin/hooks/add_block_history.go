package hooks

import (
	"fmt"
	"log"

	ocr2keepersv3 "github.com/smartcontractkit/chainlink-automation/pkg/v3"
	"github.com/smartcontractkit/chainlink-automation/pkg/v3/telemetry"
	"github.com/smartcontractkit/chainlink-automation/pkg/v3/types"
)

type AddBlockHistoryHook struct {
	metadata types.MetadataStore
	logger   *log.Logger
}

func NewAddBlockHistoryHook(ms types.MetadataStore, logger *log.Logger) AddBlockHistoryHook {
	return AddBlockHistoryHook{
		metadata: ms,
		logger:   log.New(logger.Writer(), fmt.Sprintf("[%s | build hook:add-block-history]", telemetry.ServiceName), telemetry.LogPkgStdFlags)}
}

func (h *AddBlockHistoryHook) RunHook(obs *ocr2keepersv3.AutomationObservation, limit int) {
	blockHistory := h.metadata.GetBlockHistory()
	if len(blockHistory) > limit {
		blockHistory = blockHistory[:limit]
	}
	obs.BlockHistory = blockHistory
	h.logger.Printf("adding %d blocks to observation", len(blockHistory))
}
