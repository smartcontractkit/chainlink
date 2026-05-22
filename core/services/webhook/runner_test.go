package webhook_test

import (
	"testing"

	"github.com/google/uuid"
	"github.com/stretchr/testify/require"

	"github.com/smartcontractkit/chainlink-common/pkg/utils/jsonserializable"
	"github.com/smartcontractkit/chainlink/v2/core/services/job"
	"github.com/smartcontractkit/chainlink/v2/core/services/webhook"
)

func TestDeprecatedJobRunner(t *testing.T) {
	t.Parallel()

	_, err := webhook.DeprecatedJobRunner{}.RunJob(t.Context(), uuid.New(), "{}", jsonserializable.JSONSerializable{})
	require.Error(t, err)
	require.ErrorIs(t, err, job.ErrJobTypeRemoved)
}
