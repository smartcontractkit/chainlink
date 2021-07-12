package headtracker

import (
	"github.com/smartcontractkit/chainlink/core/logger"
	httypes "github.com/smartcontractkit/chainlink/core/services/headtracker/types"
	"github.com/smartcontractkit/chainlink/core/store/models"
)

var _ httypes.Tracker = &NullTracker{}

type NullTracker struct{}

func (n *NullTracker) HighestSeenHeadFromDB() (*models.Head, error) {
	return nil, nil
}
func (*NullTracker) Start() error             { return nil }
func (*NullTracker) Stop() error              { return nil }
func (*NullTracker) SetLogger(*logger.Logger) {}
