package src

import (
	"testing"

	"github.com/gkampitakis/go-snaps/snaps"
)

func TestCreateKeystoneWorkflowJob(t *testing.T) {
	workflowConfig := WorkflowJobSpecConfig{
		JobSpecName:          "keystone_workflow",
		WorkflowOwnerAddress: "0x1234567890abcdef1234567890abcdef12345678",
		FeedIDs:              []string{"feed1", "feed2", "feed3"},
		TargetID:             "target_id",
		TargetAddress:        "0xabcdefabcdefabcdefabcdefabcdefabcdef",
	}

	output := createKeystoneWorkflowJob(workflowConfig)

	snaps.MatchSnapshot(t, output)
}
