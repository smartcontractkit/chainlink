package core

import (
	"context"

	ocr2keepers "github.com/smartcontractkit/chainlink-common/pkg/types/automation"
)

// UpkeepStateReader is the interface for reading the current state of upkeeps.
type UpkeepStateReader interface {
	SelectByWorkIDs(ctx context.Context, workIDs ...string) ([]ocr2keepers.UpkeepState, error)
}
