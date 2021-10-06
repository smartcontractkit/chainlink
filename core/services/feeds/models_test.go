package feeds

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func Test_JobProposal_CanEdit(t *testing.T) {
	testCases := []struct {
		name   string
		status JobProposalStatus
		want   bool
	}{
		{
			name:   "pending",
			status: JobProposalStatusPending,
			want:   true,
		},
		{
			name:   "cancelled",
			status: JobProposalStatusCancelled,
			want:   true,
		},
		{
			name:   "approved",
			status: JobProposalStatusApproved,
			want:   false,
		},
		{
			name:   "rejected",
			status: JobProposalStatusRejected,
			want:   false,
		},
	}

	for _, tc := range testCases {
		tc := tc
		t.Run(tc.name, func(t *testing.T) {
			t.Parallel()

			jp := &JobProposal{Status: tc.status}
			assert.Equal(t, tc.want, jp.CanEditSpec())
		})
	}
}
