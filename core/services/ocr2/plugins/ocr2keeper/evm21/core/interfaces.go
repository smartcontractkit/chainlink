package core

import (
	"context"

	ocr2keepers "github.com/smartcontractkit/ocr2keepers/pkg/v3/types"
)

// UpkeepStateReader is the interface for reading the current state of upkeeps.
//
//go:generate mockery --quiet --name UpkeepStateReader --output ./mocks/ --case=underscore
type UpkeepStateReader interface {
	SelectByWorkIDs(ctx context.Context, workIDs ...string) ([]ocr2keepers.UpkeepState, error)
}
