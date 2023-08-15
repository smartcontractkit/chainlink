package core

import (
	"context"

	ocr2keepers "github.com/smartcontractkit/ocr2keepers/pkg/v3/types"
)

// UpkeepStateReader is the interface for reading the current state of upkeeps.
type UpkeepStateReader interface {
	SelectByWorkIDsInRange(ctx context.Context, start, end int64, workIDs ...string) ([]ocr2keepers.UpkeepState, error)
}
