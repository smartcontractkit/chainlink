package triggerfns_test

import (
	"testing"

	_ "github.com/smartcontractkit/chainlink/core/services/fluxmonitor"
	"github.com/smartcontractkit/chainlink/core/store/models/triggerfns"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestTriggerFnsJSONUnmarshaling(t *testing.T) {
	var tfns triggerfns.TriggerFns
	triggerJSON := `{"absoluteThreshold": 0.0059, "relativeThreshold": 0.59}`
	require.NoError(t, tfns.UnmarshalJSON([]byte(triggerJSON)))
	assert.Len(t, tfns, 2)
}
