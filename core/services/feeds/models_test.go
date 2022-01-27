package feeds

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func Test_JobProposal_CanEditDefinition(t *testing.T) {
	testCases := []struct {
		name   string
		status SpecStatus
		want   bool
	}{
		{
			name:   "pending",
			status: SpecStatusPending,
			want:   true,
		},
		{
			name:   "cancelled",
			status: SpecStatusCancelled,
			want:   true,
		},
		{
			name:   "approved",
			status: SpecStatusApproved,
			want:   false,
		},
		{
			name:   "rejected",
			status: SpecStatusRejected,
			want:   false,
		},
	}

	for _, tc := range testCases {
		tc := tc
		t.Run(tc.name, func(t *testing.T) {
			t.Parallel()

			jp := &JobProposalSpec{Status: tc.status}
			assert.Equal(t, tc.want, jp.CanEditDefinition())
		})
	}
}
