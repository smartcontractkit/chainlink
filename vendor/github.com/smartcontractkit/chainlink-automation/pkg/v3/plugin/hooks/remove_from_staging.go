package hooks

import (
	"fmt"
	"log"

	ocr2keepersv3 "github.com/smartcontractkit/chainlink-automation/pkg/v3"
	"github.com/smartcontractkit/chainlink-automation/pkg/v3/telemetry"
	"github.com/smartcontractkit/chainlink-automation/pkg/v3/types"
)

func NewRemoveFromStagingHook(store types.ResultStore, logger *log.Logger) RemoveFromStagingHook {
	return RemoveFromStagingHook{
		store:  store,
		logger: log.New(logger.Writer(), fmt.Sprintf("[%s | pre-build hook:remove-from-staging]", telemetry.ServiceName), telemetry.LogPkgStdFlags),
	}
}

type RemoveFromStagingHook struct {
	store  types.ResultStore
	logger *log.Logger
}

func (hook *RemoveFromStagingHook) RunHook(outcome ocr2keepersv3.AutomationOutcome) {
	toRemove := make([]string, 0, len(outcome.AgreedPerformables))
	for _, result := range outcome.AgreedPerformables {
		toRemove = append(toRemove, result.WorkID)
	}

	hook.logger.Printf("%d results found in outcome for removal", len(toRemove))
	hook.store.Remove(toRemove...)
}
