package utils

import (
	"fmt"

	"github.com/smartcontractkit/chainlink/deployment/data-streams/utils/pointer"
)

type JobType *string

var (
	JobTypeLLO    JobType = pointer.To("llo")
	JobTypeStream JobType = pointer.To("stream")
)

const (
	ProductLabel = "data-streams"
)

var (
	LabelJobType     = "job-type"
	LabelStreamID    = "stream-id"
	LabelNodeType    = "node-type"
	LabelEnvironment = "environment"
	LabelProduct     = "product"
)

// DonIdentifier generates a unique identifier for a DON based on its ID and name.
func DonIdentifier(donID uint64, donName string) string {
	return fmt.Sprintf("don-%d-%s", donID, donName)
}
